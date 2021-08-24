package console

import (
	"testing"

	"github.com/krystal/katapult-cli/internal/golden"
	"github.com/krystal/katapult-cli/internal/keystrokes"
	"github.com/stretchr/testify/assert"
)

func TestSelector(t *testing.T) {
	tests := []struct {
		name string

		inputs     [][]byte
		columns    []string
		items      interface{}
		multiple   bool
		shouldExit bool
		result     interface{}
	}{
		// Non-row selection

		{
			name: "display non-row selection menu",
			inputs: [][]byte{
				keystrokes.CTRLC,
			},
			shouldExit: true,
			items: []string{
				"hello", "world",
			},
		},
		{
			name: "scroll down on non-row selection menu",
			inputs: [][]byte{
				keystrokes.DownArrow,
				keystrokes.CTRLC,
			},
			shouldExit: true,
			items: []string{
				"hello", "world",
			},
		},
		{
			name: "wrap around on non-row selection menu",
			inputs: [][]byte{
				keystrokes.DownArrow,
				keystrokes.DownArrow,
				keystrokes.CTRLC,
			},
			shouldExit: true,
			items: []string{
				"hello", "world",
			},
		},
		{
			name: "scroll up on non-row selection menu",
			inputs: [][]byte{
				keystrokes.DownArrow,
				keystrokes.UpArrow,
				keystrokes.CTRLC,
			},
			shouldExit: true,
			items: []string{
				"hello", "world",
			},
		},
		{
			name: "search non-row selection menu",
			inputs: [][]byte{
				{'h'},
				{'e'},
				keystrokes.CTRLC,
			},
			shouldExit: true,
			items: []string{
				"hello", "world",
			},
		},
		{
			name: "make selection on non-row selection menu",
			inputs: [][]byte{
				keystrokes.DownArrow,
				keystrokes.Enter,
			},
			items: []string{
				"hello", "world",
			},
			result: []string{"world"},
		},

		// Row selection

		{
			name: "display row selection menu",
			inputs: [][]byte{
				keystrokes.CTRLC,
			},
			shouldExit: true,
			columns:    []string{"test"},
			items: [][]string{
				{"hello"}, {"world"},
			},
		},
		{
			name: "scroll down on row selection menu",
			inputs: [][]byte{
				keystrokes.DownArrow,
				keystrokes.CTRLC,
			},
			shouldExit: true,
			columns:    []string{"test"},
			items: [][]string{
				{"hello"}, {"world"},
			},
		},
		{
			name: "wrap around on row selection menu",
			inputs: [][]byte{
				keystrokes.DownArrow,
				keystrokes.DownArrow,
				keystrokes.CTRLC,
			},
			shouldExit: true,
			columns:    []string{"test"},
			items: [][]string{
				{"hello"}, {"world"},
			},
		},
		{
			name: "scroll up on row selection menu",
			inputs: [][]byte{
				keystrokes.DownArrow,
				keystrokes.UpArrow,
				keystrokes.CTRLC,
			},
			shouldExit: true,
			columns:    []string{"test"},
			items: [][]string{
				{"hello"}, {"world"},
			},
		},
		{
			name: "search row selection menu",
			inputs: [][]byte{
				{'h'},
				{'e'},
				keystrokes.CTRLC,
			},
			shouldExit: true,
			columns:    []string{"test"},
			items: [][]string{
				{"hello"}, {"world"},
			},
		},
		{
			name: "make selection on row selection menu",
			inputs: [][]byte{
				keystrokes.DownArrow,
				keystrokes.Enter,
			},
			columns: []string{"test"},
			items: [][]string{
				{"hello"}, {"world"},
			},
			result: [][]string{{"world"}},
		},

		// Non-row multi selection

		{
			name: "display non-row multi selection menu",
			inputs: [][]byte{
				keystrokes.CTRLC,
			},
			multiple:   true,
			shouldExit: true,
			items: []string{
				"hello", "world",
			},
		},
		{
			name: "scroll down on non-row multi selection menu",
			inputs: [][]byte{
				keystrokes.DownArrow,
				keystrokes.CTRLC,
			},
			multiple:   true,
			shouldExit: true,
			items: []string{
				"hello", "world",
			},
		},
		{
			name: "wrap around on non-row multi selection menu",
			inputs: [][]byte{
				keystrokes.DownArrow,
				keystrokes.DownArrow,
				keystrokes.CTRLC,
			},
			multiple:   true,
			shouldExit: true,
			items: []string{
				"hello", "world",
			},
		},
		{
			name: "scroll up on non-row multi selection menu",
			inputs: [][]byte{
				keystrokes.DownArrow,
				keystrokes.UpArrow,
				keystrokes.CTRLC,
			},
			multiple:   true,
			shouldExit: true,
			items: []string{
				"hello", "world",
			},
		},
		{
			name: "search non-row multi selection menu",
			inputs: [][]byte{
				{'h'},
				{'e'},
				keystrokes.CTRLC,
			},
			multiple:   true,
			shouldExit: true,
			items: []string{
				"hello", "world",
			},
		},
		{
			name: "make selection on non-row multi selection menu",
			inputs: [][]byte{
				keystrokes.Enter,
				keystrokes.DownArrow,
				keystrokes.Enter,
				keystrokes.DownArrow,
				keystrokes.Escape,
			},
			multiple: true,
			items: []string{
				"hello", "world", "xd",
			},
			result: []string{"hello", "world"},
		},

		// Row multi selection

		{
			name: "display row multi selection menu",
			inputs: [][]byte{
				keystrokes.CTRLC,
			},
			multiple:   true,
			shouldExit: true,
			columns:    []string{"test", "123"},
			items: [][]string{
				{"hello", "a"}, {"world", "b"},
			},
		},
		{
			name: "scroll down on row multi selection menu",
			inputs: [][]byte{
				keystrokes.DownArrow,
				keystrokes.CTRLC,
			},
			multiple:   true,
			shouldExit: true,
			columns:    []string{"test", "123"},
			items: [][]string{
				{"hello", "a"}, {"world", "b"},
			},
		},
		{
			name: "wrap around on row multi selection menu",
			inputs: [][]byte{
				keystrokes.DownArrow,
				keystrokes.DownArrow,
				keystrokes.CTRLC,
			},
			multiple:   true,
			shouldExit: true,
			columns:    []string{"test", "123"},
			items: [][]string{
				{"hello", "a"}, {"world", "b"},
			},
		},
		{
			name: "scroll up on row multi selection menu",
			inputs: [][]byte{
				keystrokes.DownArrow,
				keystrokes.UpArrow,
				keystrokes.CTRLC,
			},
			multiple:   true,
			shouldExit: true,
			columns:    []string{"test", "123"},
			items: [][]string{
				{"hello", "a"}, {"world", "b"},
			},
		},
		{
			name: "search row multi selection menu",
			inputs: [][]byte{
				{'h'},
				{'e'},
				keystrokes.CTRLC,
			},
			multiple:   true,
			shouldExit: true,
			columns:    []string{"test", "123"},
			items: [][]string{
				{"hello", "a"}, {"world", "b"},
			},
		},
		{
			name: "make selection on row multi selection menu",
			inputs: [][]byte{
				keystrokes.Enter,
				keystrokes.DownArrow,
				keystrokes.Enter,
				keystrokes.DownArrow,
				keystrokes.Escape,
			},
			multiple: true,
			columns:  []string{"test", "123"},
			items: [][]string{
				{"hello", "a"}, {"world", "b"}, {"xd", "c"},
			},
			result: [][]string{{"hello", "a"}, {"world", "b"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdin := &StdinDripFeeder{Inputs: tt.inputs}
			stdout := &MockTerminal{}
			res := selectorComponent("test", tt.columns, tt.items, stdin, tt.multiple, stdout)
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
