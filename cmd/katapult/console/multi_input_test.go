package console

import (
	"strings"
	"testing"

	"github.com/krystal/katapult-cli/internal/keystrokes"

	"github.com/krystal/katapult-cli/internal/golden"
	"github.com/stretchr/testify/assert"
)

func Test_prepStringForTableView(t *testing.T) {
	tests := []struct {
		name string

		content string
		chunks  []string
	}{
		{
			name:    "blank",
			content: "",
			chunks:  []string{"│     │"},
		},
		{
			name:    "under length",
			content: "hey",
			chunks:  []string{"│hey  │"},
		},
		{
			name:    "equal to length",
			content: "hello",
			chunks:  []string{"│hello│"},
		},
		{
			name:    "split equal length",
			content: "helloworld",
			chunks:  []string{"│hello│", "│world│"},
		},
		{
			name:    "split pad",
			content: "helloxd",
			chunks:  []string{"│hello│", "│xd   │"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.chunks, prepStringForTableView(tt.content, len(tt.content), 5, false))
		})
	}
}

func Test_renderCursor(t *testing.T) {
	tests := []struct {
		name string

		content     string
		width       int
		highlighted int
		expected    string
		active      bool
	}{
		{
			name:        "under width inactive",
			content:     "hello",
			width:       10,
			highlighted: 5,
			expected:    "hello░    ",
		},
		{
			name:        "under width midpoint inactive",
			content:     "hello",
			width:       10,
			highlighted: 4,
			expected:    "hell░o    ",
		},
		{
			name:        "full minus 1 inactive",
			content:     "aaa",
			width:       4,
			highlighted: 3,
			expected:    "aaa░",
		},
		{
			name:        "full minus 1 midpoint inactive",
			content:     "aaa",
			width:       4,
			highlighted: 2,
			expected:    "aa░a",
		},
		{
			name:        "content overflow first chunk midpoint inactive",
			content:     "aaax",
			width:       4,
			highlighted: 2,
			expected:    "aa░a",
		},
		{
			name:        "content overflow first chunk end inactive",
			content:     "aaax",
			width:       4,
			highlighted: 3,
			expected:    "aaa░",
		},
		{
			name:        "content overflow one over inactive",
			content:     "abcd",
			width:       3,
			highlighted: 3,
			expected:    "bc░",
		},
		{
			name:        "content overflow not first chunk inactive",
			content:     "abcd",
			width:       3,
			highlighted: 4,
			expected:    "cd░",
		},

		{
			name:        "under width active",
			content:     "hello",
			width:       10,
			highlighted: 5,
			expected:    "hello▓    ",
			active:      true,
		},
		{
			name:        "under width midpoint active",
			content:     "hello",
			width:       10,
			highlighted: 4,
			expected:    "hell▓o    ",
			active:      true,
		},
		{
			name:        "full minus 1 active",
			content:     "aaa",
			width:       4,
			highlighted: 3,
			expected:    "aaa▓",
			active:      true,
		},
		{
			name:        "full minus 1 midpoint active",
			content:     "aaa",
			width:       4,
			highlighted: 2,
			expected:    "aa▓a",
			active:      true,
		},
		{
			name:        "content overflow first chunk midpoint active",
			content:     "aaax",
			width:       4,
			highlighted: 2,
			expected:    "aa▓a",
			active:      true,
		},
		{
			name:        "content overflow first chunk end active",
			content:     "aaax",
			width:       4,
			highlighted: 3,
			expected:    "aaa▓",
			active:      true,
		},
		{
			name:        "content overflow one over active",
			content:     "abcd",
			width:       3,
			highlighted: 3,
			expected:    "bc▓",
			active:      true,
		},
		{
			name:        "content overflow not first chunk active",
			content:     "abcd",
			width:       3,
			highlighted: 4,
			expected:    "cd▓",
			active:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, renderCursor(tt.content, tt.width, tt.highlighted, tt.active))
		})
	}
}

func createFormattedInputTestResult(expectInput string) [3]string {
	topBottomShared := strings.Repeat("─", len(expectInput)-2)
	top := "┌" + topBottomShared + "┐"
	bottom := "└" + topBottomShared + "┘"
	return [3]string{
		top, "│" + expectInput + "│", bottom,
	}
}

func Test_createInputChunk(t *testing.T) {
	tests := []struct {
		name string

		content     string
		width       int
		highlighted int
		expected    [3]string
		active      bool
	}{
		{
			name:        "under width",
			content:     "hello",
			width:       12,
			highlighted: 5,
			expected:    createFormattedInputTestResult("hello▓    "),
			active:      true,
		},
		{
			name:        "under width midpoint",
			content:     "hello",
			width:       12,
			highlighted: 4,
			expected:    createFormattedInputTestResult("hell▓o    "),
			active:      true,
		},
		{
			name:        "full minus 1",
			content:     "aaa",
			width:       6,
			highlighted: 3,
			expected:    createFormattedInputTestResult("aaa▓"),
			active:      true,
		},
		{
			name:        "full minus 1 midpoint",
			content:     "aaa",
			width:       6,
			highlighted: 2,
			expected:    createFormattedInputTestResult("aa▓a"),
			active:      true,
		},
		{
			name:        "content overflow first chunk midpoint",
			content:     "aaax",
			width:       6,
			highlighted: 2,
			expected:    createFormattedInputTestResult("aa▓a"),
			active:      true,
		},
		{
			name:        "content overflow first chunk end",
			content:     "aaax",
			width:       6,
			highlighted: 3,
			expected:    createFormattedInputTestResult("aaa▓"),
			active:      true,
		},
		{
			name:        "content overflow one over",
			content:     "abcd",
			width:       5,
			highlighted: 3,
			expected:    createFormattedInputTestResult("bc▓"),
			active:      true,
		},
		{
			name:        "content overflow not first chunk",
			content:     "abcd",
			width:       5,
			highlighted: 4,
			expected:    createFormattedInputTestResult("cd▓"),
			active:      true,
		},
		{
			name:        "inactive cursor",
			content:     "abcd",
			width:       5,
			highlighted: 4,
			expected:    createFormattedInputTestResult("cd░"),
			active:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, createInputChunk(tt.content, tt.width, tt.highlighted, tt.active))
		})
	}
}

func Test_renderInputField(t *testing.T) {
	tests := []struct {
		name string

		field       InputField
		content     string
		highlighted int
		width       int
		active      bool
	}{
		{
			name: "under width inactive",
			field: InputField{
				Optional:    false,
				Name:        "a",
				Description: "b",
			},
			content:     "hello",
			highlighted: 5,
			width:       14,
		},
		{
			name: "under width active",
			field: InputField{
				Optional:    false,
				Name:        "a",
				Description: "b",
			},
			content:     "hello",
			highlighted: 5,
			width:       14,
			active:      true,
		},
		{
			name: "optional",
			field: InputField{
				Optional:    true,
				Name:        "a",
				Description: "b",
			},
			content:     "hello",
			highlighted: 5,
			width:       14,
			active:      true,
		},
		{
			name: "title overflow",
			field: InputField{
				Optional:    false,
				Name:        "abcabcabcabcabc",
				Description: "b",
			},
			content:     "hello",
			highlighted: 5,
			width:       14,
			active:      true,
		},
		{
			name: "title overflow optional",
			field: InputField{
				Optional:    true,
				Name:        "abcabcabcabcabc",
				Description: "b",
			},
			content:     "hello",
			highlighted: 5,
			width:       14,
			active:      true,
		},
		{
			name: "description overflow",
			field: InputField{
				Optional:    true,
				Name:        "a",
				Description: "abcabcabcabcabc",
			},
			content:     "hello",
			highlighted: 5,
			width:       14,
			active:      true,
		},
		{
			name: "input exact length",
			field: InputField{
				Optional:    true,
				Name:        "a",
				Description: "b",
			},
			content:     "hello12",
			highlighted: 7,
			width:       14,
			active:      true,
		},
		{
			name: "input overflow",
			field: InputField{
				Optional:    true,
				Name:        "a",
				Description: "b",
			},
			content:     "hello123",
			highlighted: 7,
			width:       14,
			active:      true,
		},
		{
			name: "input overflow scroll",
			field: InputField{
				Optional:    true,
				Name:        "a",
				Description: "b",
			},
			content:     "hello123",
			highlighted: 8,
			width:       14,
			active:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := renderInputField(tt.field, tt.content, tt.highlighted, tt.width, tt.active)
			var stringRep string
			if s == nil {
				stringRep = "null"
			} else {
				stringRep = strings.Join(s, "\n")
			}
			if golden.Update() {
				golden.Set(t, []byte(stringRep))
			}
			assert.Equal(t, string(golden.Get(t)), stringRep)
		})
	}
}

func TestMultiInput(t *testing.T) {
	tests := []struct {
		name string

		inputs     [][]byte
		fields     []InputField
		shouldExit bool
		result     []string
	}{
		{
			name: "left right arrows",
			inputs: [][]byte{
				{'a'},
				{'c'},
				{'2'},
				{'3'},
				keystrokes.LeftArrow,
				keystrokes.LeftArrow,
				{'1'},
				keystrokes.LeftArrow, keystrokes.LeftArrow,
				keystrokes.LeftArrow, keystrokes.RightArrow,
				{'b'},
				keystrokes.Enter,
			},
			fields: []InputField{
				{
					Name:        "a",
					Description: "b",
				},
			},
			shouldExit: false,
			result:     []string{"abc123"},
		},
		{
			name: "required enforcement",
			inputs: [][]byte{
				keystrokes.Enter, keystrokes.DownArrow,
				{'b'},
				keystrokes.UpArrow,
				{'a'},
				keystrokes.Enter,
			},
			fields: []InputField{
				{
					Optional:    false,
					Name:        "top",
					Description: "the top item",
				},
				{
					Optional:    true,
					Name:        "overflow",
					Description: "the overflow item",
				},
			},
			shouldExit: false,
			result:     []string{"a", "b"},
		},
		{
			name:   "overflow",
			inputs: [][]byte{{'a'}, keystrokes.DownArrow, {'b'}, keystrokes.Enter},
			fields: []InputField{
				{
					Optional:    false,
					Name:        "top",
					Description: "the top item",
				},
				{
					Optional:    true,
					Name:        "overflow",
					Description: "the overflow item",
				},
			},
			shouldExit: false,
			result:     []string{"a", "b"},
		},
		{
			name:       "ctrl c",
			inputs:     [][]byte{keystrokes.CTRLC},
			fields:     []InputField{{}},
			shouldExit: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdin := &StdinDripFeeder{Inputs: tt.inputs}
			stdout := &MockTerminal{CustomWidth: 50}
			res := MultiInput(tt.fields, stdin, stdout)
			if tt.shouldExit {
				assert.Equal(t, tt.shouldExit, stdout.ExitSignaled)
			} else {
				assert.Equal(t, tt.result, res)
			}
			if golden.Update() {
				golden.Set(t, stdout.Buffer.Bytes())
				return
			}
			assert.Equal(t, string(golden.Get(t)), stdout.Buffer.String())
		})
	}
}
