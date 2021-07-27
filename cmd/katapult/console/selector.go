package console

import (
	"container/list"
	"github.com/buger/goterm"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"strings"
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

func selectorComponent(question string, items []string, stdin io.Reader, multiple bool) []string {
	// Pre-initialise things we need below.
	query := ""
	buf := make([]byte, 3)
	highlightIndex := 0
	var selectedItems *list.List
	if multiple {
		// Allocate a list for selections.
		selectedItems = list.New()
	}
	usableItemRows := goterm.Height() - 2
	if 0 >= usableItemRows {
		// Weird. Return status code 1.
		os.Exit(1)
	}

	// Loop until we match.
	for {
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
		if multiple {
			// Add the multiple clarification string onto the question.
			_, _ = goterm.Print(goterm.Color(question+clarificationStringMultiple, goterm.GREEN))
		} else {
			// Add the single clarification string on o the question.
			_, _ = goterm.Print(goterm.Color(question+clarificationStringSingle, goterm.GREEN))
		}

		// Format the query.
		if len(matched) == 0 {
			// There's no matches, we should just just print the users input.
			_, _ = goterm.Println(query)
		} else {
			// We should print it inside the highlighted result.
			highlighted := matched[highlightIndex]
			index := strings.Index(strings.ToLower(highlighted), queryLower)
			start := highlighted[:index]
			end := highlighted[index+len(query):]
			_, _ = goterm.Println(goterm.Color(start, goterm.BLUE) + query + goterm.Color(end, goterm.BLUE))
		}

		// Display the rest of the items.
		matchedLen := len(matched)
		matchedStart := 0
		if highlightIndex >= usableItemRows {
			matchedStart = highlightIndex - usableItemRows + 1
		}
		for i := matchedStart; i < intMin(matchedLen, matchedStart + usableItemRows); i++ {
			// Get the match.
			v := matched[i]

			// Handle rendering the item.
			if i == highlightIndex {
				// Highlight this item.
				_, _ = goterm.
					Println(goterm.Color(v, goterm.YELLOW))
			} else {
				// Check if it is selected.
				found := false
				if multiple {
					for e := selectedItems.Front(); e != nil; e = e.Next() {
						if e.Value.(string) == v {
							found = true
							break
						}
					}
				}

				// Print this item in a different color depending if it is highlighted or not.
				if found {
					_, _ = goterm.Println(goterm.Color(v, goterm.CYAN))
				} else {
					_, _ = goterm.Println(v)
				}
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
			case 127:
				// Backspace
				if len(query) == 0 {
					continue
				}
				query = query[:len(query)-1]
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
				}
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
	// TODO: Handle escape for multiple!
}

// FuzzySelector is used to create a selector which also fuzzy searches the items and allows for the selection of one item.
func FuzzySelector(question string, items []string, stdin io.Reader) string {
	return selectorComponent(question, items, stdin, false)[0]
}

// FuzzyMultiSelector is used to create a selector which also fuzzy searches the items and allows for the selection of multiple items.
func FuzzyMultiSelector(question string, items []string, stdin io.Reader) []string {
	return selectorComponent(question, items, stdin, true)
}
