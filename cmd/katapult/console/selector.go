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

func selectorComponent(question string, items []string, stdin io.Reader, multiple bool, onRender func()) []string {
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
		var matched []string
		queryLower := strings.ToLower(query)
		for _, v := range items {
			if strings.Contains(strings.ToLower(v), query) {
				// The query is contained in the string.
				matched = append(matched, v)
			}
		}
		if matched == nil {
			matched = []string{}
		}
		if highlightIndex >= len(matched) {
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
		if len(matched) == 0 {
			// There's no matches, we should just just print the users input.
			_, _ = goterm.Println(query)
			suggestionLen = len(query)
		} else {
			// We should print it inside the highlighted result.
			highlighted := matched[highlightIndex]
			suggestionLen = len(highlighted)
			index := strings.Index(strings.ToLower(highlighted), queryLower)
			start := highlighted[:index]
			end := highlighted[index+len(query):]
			_, _ = goterm.Println(goterm.Color(start, goterm.BLUE) + query + goterm.Color(end, goterm.BLUE))
		}
		roughLines := int(math.Ceil(float64(len(questionFormatted) + suggestionLen) / float64(width)))
		usableItemRows -= roughLines

		// Display the rest of the items.
		matchedLen := len(matched)
		matchedStart := 0
		if highlightIndex >= usableItemRows {
			matchedStart = highlightIndex - usableItemRows + 1
		}
		for i := matchedStart; i < intMin(matchedLen, matchedStart + usableItemRows); i++ {
			// Get the match.
			v := matched[i]

			// Handle rendering selections in a multiple context.
			if multiple {
				for e := selectedItems.Front(); e != nil; e = e.Next() {
					if e.Value.(string) == v {
						_, _ = goterm.Print(goterm.Color("[*] ", goterm.GREEN))
						goto renderItem
					}
				}
				_, _ = goterm.Print(goterm.Color("[ ] ", goterm.RED))
			}

			// Handle rendering the item.
			renderItem:
			if i == highlightIndex {
				// Highlight this item.
				_, _ = goterm.Println(goterm.Color(v, goterm.YELLOW))
			} else {
				// Print the item.
				_, _ = goterm.Println(v)
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
				if len(matched) != 0 {
					if !multiple {
						// If we aren't to match multiple, enter means we should return it.
						return []string{matched[highlightIndex]}
					}

					// Handle the selection in a multiple context.
					item := matched[highlightIndex]
					found := false
					for e := selectedItems.Front(); e != nil; e = e.Next() {
						if e.Value.(string) == item {
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
	return selectorComponent(question, items, stdin, false, nil)[0]
}

// FuzzyMultiSelector is used to create a selector which also fuzzy searches the items and allows for the selection of multiple items.
func FuzzyMultiSelector(question string, items []string, stdin io.Reader) []string {
	return selectorComponent(question, items, stdin, true, nil)
}
