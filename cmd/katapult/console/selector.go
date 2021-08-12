package console

import (
	"container/list"
	"io"
	"math"
	"os"
	"strings"

	"github.com/buger/goterm"
	"golang.org/x/term"
)

const (
	// Clarification string for a single item.
	clarificationStringSingle = " (Press ENTER to make your selection): "

	// Clarification string for multiple items.
	clarificationStringMultiple = " (Press ENTER to select items and ESC when you are done with your selections): "
)

func intMin(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func exactCompare(a, b interface{}) bool {
	switch x := a.(type) {
	case []string:
		// The logic behind this: A string slice is just a fancy pointer with some extra data.
		// If we compare the pointers, we can tell if they are equal.
		return &x[0] == &b.([]string)[0]
	default:
		return a == b
	}
}

// Get query matches. Returns slice, lower case query and length.
func getQueryMatches(query string, hasColumns bool, items interface{}) (interface{}, string, int) {
	query = strings.ToLower(query)
	var matched interface{}
	matchedLen := 0
	if hasColumns {
		// Handle string array matches.
		ungeneric := items.([][]string)
		for _, v := range ungeneric {
			found := false
			for _, x := range v {
				if strings.Contains(strings.ToLower(x), query) {
					found = true
					break
				}
			}
			if found {
				// The query is contained in one of the strings.
				x, _ := matched.([][]string)
				//nolint:gocritic
				matched = append(x, v)
				matchedLen = len(x) + 1
			}
		}
	} else {
		// Handle string matches.
		ungeneric := items.([]string)
		for _, v := range ungeneric {
			if strings.Contains(strings.ToLower(v), query) {
				// The query is contained in the string.
				x, _ := matched.([]string)
				//nolint:gocritic
				matched = append(x, v)
				matchedLen = len(x) + 1
			}
		}
	}
	if matched == nil {
		if hasColumns {
			matched = [][]string{}
		} else {
			matched = []string{}
		}
	}
	return matched, query, matchedLen
}

// Formats the user prompt. Returns the rough line count.
func formatUserPrompt(length, highlightIndex int, hasColumns bool, matched interface{}, query, queryLower string) int {
	var suggestionLen int
	if length == 0 {
		// There's no matches, we should just just print the users input.
		_, _ = goterm.Println(query)
		suggestionLen = len(query)
	} else {
		// We should print it inside the highlighted result.
		var highlighted string
		if hasColumns {
			highlighted = strings.Join(matched.([][]string)[highlightIndex], " / ")
		} else {
			highlighted = matched.([]string)[highlightIndex]
		}
		suggestionLen = len(highlighted)
		index := strings.Index(strings.ToLower(highlighted), queryLower)
		start := highlighted[:index]
		end := highlighted[index+len(query):]
		_, _ = goterm.Println(goterm.Color(start, goterm.BLUE) + query + goterm.Color(end, goterm.BLUE))
	}
	return suggestionLen
}

// Handles column rendering.
func renderColumns(row []string, offset, highlight, width int) {
	// Get the remaining width by subtracting the offset from the console width.
	remainingWidth := width
	if offset > 0 {
		// Subtract the offset from the width.
		remainingWidth = width - offset
	}

	// Get the length per column.
	lengthPerColumn := remainingWidth / len(row)
	if 0 > offset {
		// Subtract the offset from column length.
		lengthPerColumn += offset
	}

	// Go through each column.
	content := ""
	if offset > 0 {
		// Pad out for the offset.
		content = strings.Repeat(" ", offset)
	}
	for _, column := range row {
		// Check if we need to truncate or pad.
		if len(column) >= lengthPerColumn {
			// Truncate
			content += column[:lengthPerColumn]
		} else {
			// Pad
			content += column + strings.Repeat(" ", lengthPerColumn-len(column))
		}
	}
	switch highlight {
	case 1:
		// Content highlight
		content = goterm.Color(content, goterm.YELLOW)
	case 2:
		// Title highlight
		content = goterm.Color(content, goterm.CYAN)
	}
	_, _ = goterm.Println(content)
}

// Handle a standard input.
func handleStandardInput(buf []byte, matchedLen int, multiple, hasColumns bool,
	highlightIndex *int, selectedItems *list.List, matched interface{}, query *string) interface{} {
	switch buf[0] {
	case 3:
		// CTRL+C
		os.Exit(1)
	case 13:
		// Enter
		if matchedLen != 0 {
			if !multiple {
				// If we aren't to match multiple, enter means we should return it.
				if hasColumns {
					return [][]string{matched.([][]string)[*highlightIndex]}
				}
				return []string{matched.([]string)[*highlightIndex]}
			}

			// Handle the selection in a multiple context.
			var item interface{}
			switch x := matched.(type) {
			case []string:
				item = x[*highlightIndex]
			case [][]string:
				item = x[*highlightIndex]
			}
			found := false
			for e := selectedItems.Front(); e != nil; e = e.Next() {
				if exactCompare(e.Value, item) {
					selectedItems.Remove(e)
					found = true
					break
				}
			}
			if !found {
				selectedItems.PushBack(item)
			}
			return nil
		}
	case 27:
		// Escape
		if multiple {
			if !hasColumns {
				// These aren't rows. Treat them as strings.
				a := make([]string, selectedItems.Len())
				i := 0
				for e := selectedItems.Front(); e != nil; e = e.Next() {
					a[i] = e.Value.(string)
					i++
				}
				return a
			}

			// Treat these as rows.
			a := make([][]string, selectedItems.Len())
			i := 0
			for e := selectedItems.Front(); e != nil; e = e.Next() {
				a[i] = e.Value.([]string)
				i++
			}
			return a
		}
	case 127:
		// Backspace
		if len(*query) == 0 {
			return nil
		}
		*query = (*query)[:len(*query)-1]
	default:
		// Character
		*query += string(buf[0])
	}
	*highlightIndex = 0
	return nil
}

// Handle the inputs.
func handleInput(buf []byte, multiple bool, matchedLen, n int, highlightIndex *int,
	hasColumns bool, selectedItems *list.List, query *string,
	matched interface{}) interface{} {
	if n == 1 {
		// Standard input.
		return handleStandardInput(buf, matchedLen, multiple, hasColumns, highlightIndex,
			selectedItems, matched, query)
	}

	// AT&T style key input.
	switch string(buf) {
	case string([]byte{27, 91, 65}):
		// Arrow up
		*highlightIndex--
		if *highlightIndex == -1 {
			// Don't let people try and access index -1.
			*highlightIndex = 0
		}
	case string([]byte{27, 91, 66}):
		// Arrow down
		*highlightIndex++
	default:
		// Something else. Ignore this.
	}
	return nil
}

// Renders a row item.
func renderRowItem(hasColumn bool, highlightIndex, i, width int, v interface{}) {
	if hasColumn {
		// Handle column rendering.
		highlightValue := 0
		if i == highlightIndex {
			highlightValue = 1
		}
		renderColumns(v.([]string), -2, highlightValue, width)
		return
	}

	// Handle string rendering.
	if i == highlightIndex {
		// Highlight this item.
		_, _ = goterm.Println(goterm.Color(v.(string), goterm.YELLOW))
	} else {
		// Print the item.
		_, _ = goterm.Println(v)
	}
}

// items is either []string or [][]string (if columns isn't nil).
//nolint:unparam,funlen
func selectorComponent(
	question string, columns []string, items interface{},
	stdin io.Reader, multiple bool, onRender func(),
) interface{} {
	// Pre-initialize things we need below.
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
		// Get the usable item rows.
		usableItemRows := goterm.Height() - 1
		if 0 >= usableItemRows {
			// Weird. Return status code 1.
			os.Exit(1)
		}

		// Get the width.
		width := goterm.Width()

		// Get the matched items.
		matched, queryLower, matchedLen := getQueryMatches(query, columns != nil, items)
		if highlightIndex >= matchedLen {
			// This means that the user was highlighting over something that
			// is too far down for current query. Show it at the top.
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
		suggestionLen := formatUserPrompt(matchedLen, highlightIndex, columns != nil, matched, query, queryLower)
		roughLines := int(math.Ceil(float64(len(questionFormatted)+suggestionLen) / float64(width)))
		usableItemRows -= roughLines

		// Render the title of each column.
		if columns != nil {
			offset := -2
			if multiple {
				// If there is multiple items, we want to offset the items by the checkbox size.
				offset = 4
			}
			renderColumns(columns, offset, 2, width)
			usableItemRows--
		}

		// Display the rest of the items.
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
					if exactCompare(e.Value, v) {
						_, _ = goterm.Print(goterm.Color("[*] ", goterm.GREEN))
						goto renderItem
					}
				}
				_, _ = goterm.Print(goterm.Color("[ ] ", goterm.RED))
			}

			// Handle rendering the item.
		renderItem:
			renderRowItem(columns != nil, highlightIndex, i, width, v)
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
		raw, err := term.MakeRaw(0)
		if err != nil {
			panic(err)
		}
		n, _ := stdin.Read(buf)
		_ = term.Restore(0, raw)
		if x := handleInput(buf, multiple, matchedLen,
			n, &highlightIndex, columns != nil, selectedItems, &query, matched); x != nil {
			return x
		}
	}
}

// FuzzySelector is used to create a selector.
func FuzzySelector(question string, items []string, stdin io.Reader) string {
	return selectorComponent(question, nil, items, stdin, false, nil).([]string)[0]
}

// FuzzyMultiSelector is used to create a selector with multiple items.
func FuzzyMultiSelector(question string, items []string, stdin io.Reader) []string {
	return selectorComponent(question, nil, items, stdin, true, nil).([]string)
}

// FuzzyTableMultiSelector is used to create a selector with a table.
func FuzzyTableSelector(question string, columns []string, items [][]string, stdin io.Reader) []string {
	return selectorComponent(question, columns, items, stdin, false, nil).([][]string)[0]
}

// FuzzyTableMultiSelector is used to create a selector with a table and multiple items.
func FuzzyTableMultiSelector(question string, columns []string, items [][]string, stdin io.Reader) [][]string {
	return selectorComponent(question, columns, items, stdin, true, nil).([][]string)
}
