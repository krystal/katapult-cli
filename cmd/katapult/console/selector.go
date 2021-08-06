package console

import (
	"container/list"
	"io"
	"math"
	"os"
	"strings"

	"github.com/buger/goterm"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	// Clarification string for a single item
	clarificationStringSingle = " (Press ENTER to make your selection): "

	// Clarification string for multiple items
	clarificationStringMultiple = " (Press ENTER to select items and ESC when you are done with your selections): "
)

func intMin(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// items is either []string or [][]string (if columns isn't nil)
func selectorComponent(question string, columns []string, items interface{}, stdin io.Reader, multiple bool, onRender func()) interface{} {
	// Pre-initialise things we need below.
	query := ""
	buf := make([]byte, 3)
	highlightIndex := 0
	var selectedItems *list.List
	if multiple {
		// Allocate a list for selections.
		selectedItems = list.New()
	}

	// Loop until we match.
	for {
	loopStart:
		// Get the usable item rows.
		usableItemRows := goterm.Height() - 1
		if 0 >= usableItemRows {
			// Weird. Return status code 1.
			os.Exit(1)
		}

		// Get the width.
		width := goterm.Width()

		// Get the matched items.
		var matched interface{}
		queryLower := strings.ToLower(query)
		if columns == nil {
			// Handle string matches.
			for _, v := range items.([]string) {
				if strings.Contains(strings.ToLower(v), query) {
					// The query is contained in the string.
					a, _ := matched.([]string)
					matched = append(a, v)
				}
			}
		} else {
			// Handle string array matches.
			for _, v := range items.([][]string) {
				found := false
				for _, x := range v {
					if strings.Contains(strings.ToLower(x), query) {
						found = true
						break
					}
				}
				if found {
					// The query is contained in one of the strings.
					a, _ := matched.([][]string)
					matched = append(a, v)
				}
			}
		}
		if matched == nil {
			if columns == nil {
				matched = []string{}
			} else {
				matched = [][]string{}
			}
		}
		stackLen := func(iface interface{}) int {
			if columns == nil {
				return len(matched.([]string))
			}
			return len(matched.([][]string))
		}
		if highlightIndex >= stackLen(matched) {
			highlightIndex = 0
		}

		// Clear the terminal.
		goterm.Clear()

		// Asks the question.
		var questionFormatted string
		if multiple {
			// Add the multiple clarification string onto the question.
			questionFormatted = goterm.Color(question+clarificationStringMultiple, goterm.GREEN)
		} else {
			// Add the single clarification string on o the question.
			questionFormatted = goterm.Color(question+clarificationStringSingle, goterm.GREEN)
		}
		_, _ = goterm.Print(questionFormatted)

		// Format the query.
		var suggestionLen int
		if stackLen(matched) == 0 {
			// There's no matches, we should just just print the users input.
			_, _ = goterm.Println(query)
			suggestionLen = len(query)
		} else {
			// We should print it inside the highlighted result.
			var highlighted string
			if columns == nil {
				highlighted = matched.([]string)[highlightIndex]
			} else {
				highlighted = strings.Join(matched.([][]string)[highlightIndex], " / ")
			}
			suggestionLen = len(highlighted)
			index := strings.Index(strings.ToLower(highlighted), queryLower)
			start := highlighted[:index]
			end := highlighted[index+len(query):]
			_, _ = goterm.Println(goterm.Color(start, goterm.BLUE) + query + goterm.Color(end, goterm.BLUE))
		}
		roughLines := int(math.Ceil(float64(len(questionFormatted)+suggestionLen) / float64(width)))
		usableItemRows -= roughLines

		// Handle column rendering.
		renderColumns := func(row []string, highlight bool) {
			// Get the length per column.
			lengthPerColumn := width / len(row)

			// Go through each column.
			content := ""
			for _, column := range row {
				// Check if we need to truncate or pad.
				if len(column) >= lengthPerColumn {
					// Truncate
					content += column[:lengthPerColumn]
				} else {
					// Pad
					padding := make([]byte, lengthPerColumn - len(column))
					for i := range padding {
						padding[i] = ' '
					}
					content += column
					content += string(padding)
				}
			}
			if highlight {
				content = goterm.Color(content, goterm.YELLOW)
			}
			_, _ = goterm.Println(content)
		}

		// Render the title of each column.
		if columns != nil {
			renderColumns(columns, false)
			usableItemRows--
		}

		// Display the rest of the items.
		matchedLen := stackLen(matched)
		matchedStart := 0
		if highlightIndex >= usableItemRows {
			matchedStart = highlightIndex - usableItemRows + 1
		}
		for i := matchedStart; i < intMin(matchedLen, matchedStart+usableItemRows); i++ {
			// Get the match.
			var v interface{}
			switch x := matched.(type) {
			case []string:
				v = x[i]
			case [][]string:
				v = x[i]
			}

			// Handle rendering selections in a multiple context.
			if multiple {
				for e := selectedItems.Front(); e != nil; e = e.Next() {
					if e.Value == v {
						_, _ = goterm.Print(goterm.Color("[*] ", goterm.GREEN))
						goto renderItem
					}
				}
				_, _ = goterm.Print(goterm.Color("[ ] ", goterm.RED))
			}

			// Handle rendering the item.
		renderItem:
			if columns == nil {
				// Handle string rendering.
				if i == highlightIndex {
					// Highlight this item.
					_, _ = goterm.Println(goterm.Color(v.(string), goterm.YELLOW))
				} else {
					// Print the item.
					_, _ = goterm.Println(v)
				}
			} else {
				// Handle column rendering.
				renderColumns(v.([]string), i == highlightIndex)
			}
		}

		// Print a bunch of new lines if there's less items than console rows.
		if usableItemRows > matchedLen {
			blankLines := usableItemRows - matchedLen
			for i := 0; i < blankLines; i++ {
				_, _ = goterm.Println()
			}
		}

		// Flush out the output.
		goterm.Flush()

		// Call the on render event if it exists.
		if onRender != nil {
			onRender()
		}

		// Wait for user input.
		raw, err := terminal.MakeRaw(0)
		if err != nil {
			panic(err)
		}
		n, _ := stdin.Read(buf)
		_ = terminal.Restore(0, raw)

		if n == 1 {
			// Standard input.
			switch buf[0] {
			case 3:
				// CTRL+C
				os.Exit(1)
			case 13:
				// Enter
				if stackLen(matched) != 0 {
					if !multiple {
						// If we aren't to match multiple, enter means we should return it.
						if columns == nil {
							return []string{matched.([]string)[highlightIndex]}
						}
						return [][]string{matched.([][]string)[highlightIndex]}
					}

					// Handle the selection in a multiple context.
					var item interface{}
					switch x := matched.(type) {
					case []string:
						 item = x[highlightIndex]
					case [][]string:
						item = x[highlightIndex]
					}
					found := false
					for e := selectedItems.Front(); e != nil; e = e.Next() {
						if e.Value == item {
							selectedItems.Remove(e)
							found = true
							break
						}
					}
					if !found {
						selectedItems.PushBack(item)
					}
					goto loopStart
				}
			case 27:
				// Escape
				if multiple {
					a := make([]string, selectedItems.Len())
					i := 0
					for e := selectedItems.Front(); e != nil; e = e.Next() {
						a[i] = e.Value.(string)
						i++
					}
					return a
				}
			case 127:
				// Backspace
				if len(query) == 0 {
					continue
				}
				query = query[:len(query)-1]
			default:
				// Character
				query += string(buf[0])
			}
			highlightIndex = 0
		} else {
			// AT&T style key input.
			switch string(buf) {
			case string([]byte{27, 91, 65}):
				// Arrow up
				highlightIndex--
				if highlightIndex == -1 {
					highlightIndex = 0
				}
			case string([]byte{27, 91, 66}):
				// Arrow down
				highlightIndex++
			default:
				// Something else. Ignore this.
			}
		}
	}
}

// FuzzySelector is used to create a selector which also fuzzy searches the items and allows for the selection of one item.
func FuzzySelector(question string, items []string, stdin io.Reader) string {
	return selectorComponent(question, nil, items, stdin, false, nil).([]string)[0]
}

// FuzzyMultiSelector is used to create a selector which also fuzzy searches the items and allows for the selection of multiple items.
func FuzzyMultiSelector(question string, items []string, stdin io.Reader) []string {
	return selectorComponent(question, nil, items, stdin, true, nil).([]string)
}

// FuzzyTableSelector is used to create a selector which also fuzzy searches the table and allows for the selection of one item.
func FuzzyTableSelector(question string, columns []string, items [][]string, stdin io.Reader) []string {
	return selectorComponent(question, columns, items, stdin, false, nil).([][]string)[0]
}

// FuzzyTableMultiSelector is used to create a selector which also fuzzy searches the table and allows for the selection of multiple items.
func FuzzyTableMultiSelector(question string, columns []string, items [][]string, stdin io.Reader) [][]string {
	return selectorComponent(question, columns, items, stdin, true, nil).([][]string)
}
