package console

import (
	"github.com/buger/goterm"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"strings"
)

// FuzzySelector is used to create a selector which also fuzzy searches the items
func FuzzySelector(question string, items []string, stdin io.Reader) string {
	// Pre-initialise things we need below.
	query := ""
	buf := make([]byte, 3)
	highlightIndex := 0

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
		_, _ = goterm.Print(goterm.Color(question+": ", goterm.GREEN))
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
		for i, v := range matched {
			if i == highlightIndex {
				// Highlight this item.
				_, _ = goterm.Println(goterm.Color(v, goterm.YELLOW))
			} else {
				// Print this item.
				_, _ = goterm.Println(v)
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
					return matched[highlightIndex]
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
}
