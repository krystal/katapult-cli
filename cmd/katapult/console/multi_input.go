package console

import (
	"io"
	"math"
	"strings"

	"github.com/buger/goterm"
	"github.com/krystal/katapult-cli/internal/keystrokes"
	"golang.org/x/term"
)

// InputField is used to define a field which users can input text into.
type InputField struct {
	// Optional defines if the field can be blank.
	Optional bool `json:"optional"`

	// Name defines the input boxes name.
	Name string `json:"name"`

	// Description defines the input boxes description.
	Description string `json:"description"`
}

func prepStringForTableView(s string, sLen int, l int, pad bool) []string {
	if l >= sLen {
		// The length is more or equal to the string length.
		// Return a 1 length slice that contains spacing(total
		// length minus string length) plus the string.
		padLen := l - sLen
		return []string{
			addSideBorder(s+strings.Repeat(" ", padLen), pad),
		}
	}

	// Defines the chunk count.
	chunkCount := int(math.Ceil(float64(sLen) / float64(l)))

	// Defines the string slice.
	a := make([]string, chunkCount)
	for i := 0; i < chunkCount; i++ {
		// Defines the previous count of items.
		prevItems := i * l

		// If prevItems + sLen is greater than the string
		// length, set it to the length.
		endIndex := prevItems + l
		if endIndex > sLen {
			endIndex = sLen
		}

		// Add the string.
		a[i] = s[prevItems:endIndex]
	}

	// If the final item isn't padded, we should pad it accordingly.
	final := a[chunkCount-1]
	if len(final) != l {
		a[chunkCount-1] = final + strings.Repeat(" ", l-len(final))
	}

	// Add side borders.
	for i, v := range a {
		a[i] = addSideBorder(v, pad)
	}

	// Return the slice.
	return a
}

func renderCursor(content string, width, highlighted int, active bool) string {
	contentLen := len(content)
	if highlighted > contentLen {
		// Set the highlight to the end of the string.
		highlighted = contentLen
	}
	cursor := "▓"
	if !active {
		// We should switch to the inactive cursor.
		cursor = "░"
	}
	if width > contentLen {
		// Greater than means there's enough room for the cursor too.
		padAmount := width - contentLen - 1
		if 0 > padAmount {
			// Do not panic here.
			padAmount = 0
		}
		pad := strings.Repeat(" ", padAmount)
		if contentLen == highlighted {
			// Return the content with the cursor at the end.
			return content + cursor + pad
		}
		return content[:highlighted] + cursor + content[highlighted:] + pad
	}

	// Create a width - 1 chunk and recall it.
	wm1 := width - 1
	if 0 > wm1 {
		// only room for 1 icon.
		return cursor
	}

	if highlighted == contentLen {
		// If the content length is equal to the highlighted index, cut out the width - 1 and cursor.
		return content[contentLen-wm1:] + cursor
	}

	if wm1 > highlighted {
		// In this case, we should recall it with the start.
		return renderCursor(content[:wm1], width, highlighted, active)
	}

	// Slice the text for part before the highlighted bit.
	beforeHighlighted := highlighted - wm1
	textBefore := content[beforeHighlighted:highlighted]
	return renderCursor(textBefore, width, len(textBefore), active)
}

func addSideBorder(s string, pad bool) string {
	if pad {
		return "│ " + s + " │"
	}
	return "│" + s + "│"
}

func createInputChunk(content string, width, highlighted int, active bool) [3]string {
	// Handle if there isn't enough length.
	if 3 > width {
		return [3]string{}
	}

	// Defines the top and bottom of the chunk.
	topBottomShared := strings.Repeat("─", width-2)
	top := "┌" + topBottomShared + "┐"
	bottom := "└" + topBottomShared + "┘"

	// Return the array.
	return [3]string{
		// Defines the top border.
		top,

		// Defines the input box.
		addSideBorder(renderCursor(content, width-2, highlighted, active), false),

		// Defines the bottom border.
		bottom,
	}
}

func renderInputField(field InputField, content string, highlighted int, width int, active bool) []string {
	// If the width is less than 10, return blank.
	// It's technically more space than we need, but the UX will be so unusable at this point.
	if 10 > width {
		return []string{}
	}

	// Defines the top and bottom of the field.
	topBottomShared := strings.Repeat("─", width-2)
	top := "┌" + topBottomShared + "┐"
	bottom := "└" + topBottomShared + "┘"

	// Add a red asterisk to the title if required.
	name := field.Name
	nameLen := len(name)
	if !field.Optional {
		name = goterm.Color("* ", goterm.RED) + name
		nameLen += 2
	}

	// Render the title into chunks.
	titleChunks := prepStringForTableView(name, nameLen, width-4, true)

	// Render the description into chunks.
	descriptionChunks := prepStringForTableView(field.Description, len(field.Description), width-4, true)

	// Create the input box.
	inputChunk := createInputChunk(content, width-4, highlighted, active)

	// We should render the following:
	// Top line, title, blank line, description, blank line, input field (with one space each side),
	// bottom line
	totalLen := 5 + len(titleChunks) + len(descriptionChunks)
	toRender := make([]string, totalLen)
	toRender[0] = top
	copy(toRender[1:], titleChunks)
	index := len(titleChunks) + 1
	copy(toRender[index:], descriptionChunks)
	index += len(descriptionChunks)
	for _, v := range inputChunk {
		toRender[index] = addSideBorder(v, true)
		index++
	}
	toRender[index] = bottom

	// Return the freshly created slice.
	return toRender
}

type consoleChunk struct {
	content           string
	fullyContains     []int
	partiallyContains []int
	lines             int
}

func handleKeypress(buf []byte, n, activeIndex int, fields []InputField,
	highlightedIndexes []int, fieldsContent []string, terminal TerminalInterface) (int, bool) {
	// Handle single byte.
	if n == 1 {
		switch buf[0] {
		case 3:
			// CTRL+C
			terminal.SignalInterrupt()
			return activeIndex, true
		case 13:
			// Enter
			for i, v := range fields {
				content := fieldsContent[i]
				if content == "" && !v.Optional {
					// Non-optional field missing.
					return activeIndex, false
				}
			}
			return activeIndex, true
		case 127:
			// Backspace
			stringIndex := highlightedIndexes[activeIndex]
			if stringIndex == 0 {
				// Impossible to backspace at zero index.
				return activeIndex, false
			}
			highlightedIndexes[activeIndex]--
			s := fieldsContent[activeIndex]
			fieldsContent[activeIndex] = s[:stringIndex-1] + s[stringIndex:]
		default:
			// Character
			stringIndex := highlightedIndexes[activeIndex]
			highlightedIndexes[activeIndex]++
			s := fieldsContent[activeIndex]
			if stringIndex == len(s) {
				// Add to the string.
				fieldsContent[activeIndex] = s + string(buf[0])
			} else {
				// Add into the middle of the string.
				fieldsContent[activeIndex] = s[:stringIndex] + string(buf[0]) + s[stringIndex:]
			}
		}
		return activeIndex, false
	}

	// AT&T style key input.
	switch string(buf) {
	case string(keystrokes.UpArrow):
		activeIndex--
		if activeIndex == -1 {
			// Don't let people try and access index -1.
			activeIndex = 0
		}
	case string(keystrokes.DownArrow):
		activeIndex++
		if activeIndex == len(fields) {
			// Loop to the start.
			activeIndex = 0
		}
	case string(keystrokes.LeftArrow):
		stringIndex := highlightedIndexes[activeIndex]
		if stringIndex != 0 {
			highlightedIndexes[activeIndex]--
		}
	case string(keystrokes.RightArrow):
		stringIndex := highlightedIndexes[activeIndex]
		content := fieldsContent[activeIndex]
		stringIndex++
		if stringIndex > len(content) {
			stringIndex = len(content)
		}
		highlightedIndexes[activeIndex] = stringIndex
	default:
		// Something else. Ignore this.
	}
	return activeIndex, false
}

func chunkForConsole(renderedFields [][]string, usableItemRows int) []consoleChunk {
	consoleSizedChunks := make([]consoleChunk, 0, 1)
	startFrom := 0
	for len(renderedFields) != startFrom {
		lines := 0
		chunk := ""
		fullyContains := []int{}
		partial := false
		for ; startFrom < len(renderedFields); startFrom++ {
			v := renderedFields[startFrom]
			for _, line := range v {
				if lines+1 > usableItemRows {
					partial = true
					break
				}
				lines++
				chunk += line + "\n"
			}
			if partial {
				// Partial means we should add it to the array of partial chunks and stop.
				consoleSizedChunks = append(consoleSizedChunks, consoleChunk{
					content:           chunk,
					fullyContains:     fullyContains,
					partiallyContains: []int{startFrom},
					lines:             lines,
				})
				break
			} else {
				// Add to fully contains.
				fullyContains = append(fullyContains, startFrom)
			}
		}
		if !partial {
			// No partials. We should append.
			consoleSizedChunks = append(consoleSizedChunks, consoleChunk{
				content:           chunk,
				fullyContains:     fullyContains,
				partiallyContains: []int{},
				lines:             lines,
			})
		}
	}
	return consoleSizedChunks
}

// MultiInput is used to format multiple input slots to the display.
func MultiInput(fields []InputField, stdin io.Reader, terminal TerminalInterface) []string {
	// Ensure the terminal isn't nil.
	if terminal == nil {
		terminal = gotermTerminal{}
	}

	// Defines the active index.
	activeIndex := 0

	// Defines the highlighted index in all fields.
	highlightedIndexes := make([]int, len(fields))

	// Defines the content for all fields.
	fieldsContent := make([]string, len(fields))

	// Loop until we are done.
	buf := make([]byte, 3)
	for {
		// Get the usable item rows.
		usableItemRows := terminal.Height() - 1
		if 0 >= usableItemRows {
			// Weird. Return status code 1.
			terminal.SignalInterrupt()
			return nil
		}

		// Get the width.
		width := terminal.Width()

		// Clear the terminal.
		terminal.Clear()

		// Render all fields so we can get the height of them all.
		renderedFields := make([][]string, len(fields))
		for i, v := range fields {
			renderedFields[i] = renderInputField(v, fieldsContent[i], highlightedIndexes[i], width, i == activeIndex)
		}

		// Create console sized chunks for all of the fields.
		consoleSizedChunks := chunkForConsole(renderedFields, usableItemRows)

		// Find the best chunk to display.
		containsInt := func(a []int, i int) bool {
			for _, v := range a {
				if v == i {
					return true
				}
			}
			return false
		}
		for _, chunk := range consoleSizedChunks {
			// First pass will work for most use cases where console is reasonably sized.
			if containsInt(chunk.fullyContains, activeIndex) {
				_, _ = terminal.Print(chunk.content)
				_, _ = terminal.Print(strings.Repeat("\n", usableItemRows-chunk.lines))
				goto postChunkPrint
			}
		}
		for _, chunk := range consoleSizedChunks {
			// Second pass will work in poorly sized consoles.
			if containsInt(chunk.partiallyContains, activeIndex) {
				_, _ = terminal.Print(chunk.content)
				_, _ = terminal.Print(strings.Repeat("\n", usableItemRows-chunk.lines))
				break
			}
		}

	postChunkPrint:
		// Flush out the output.
		terminal.Flush()

		// Handle keypresses.
		raw, err := terminal.MakeRaw()
		if err != nil {
			panic(err)
		}
		n, _ := stdin.Read(buf)
		if raw != nil {
			_ = term.Restore(0, raw)
		}
		var ret bool
		activeIndex, ret = handleKeypress(buf, n, activeIndex, fields, highlightedIndexes, fieldsContent, terminal)
		if ret {
			return fieldsContent
		}
	}
}
