package console

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/krystal/katapult-cli/internal/golden"
	"github.com/stretchr/testify/assert"
	"golang.org/x/term"
)

type mockTerminal struct {
	buf *bytes.Buffer

	exitSignaled bool
}

func (m mockTerminal) Height() int {
	return 10
}

func (m mockTerminal) Width() int {
	return 200
}

func (m mockTerminal) Print(items ...interface{}) (int, error) {
	return fmt.Fprint(m.buf, items...)
}

func (m mockTerminal) Println(items ...interface{}) (int, error) {
	return fmt.Fprintln(m.buf, items...)
}

func (m mockTerminal) Clear() {
	_, _ = fmt.Fprint(m.buf, "\033[2J")
}

func (m mockTerminal) Flush() {
	// This is platform dependant. Ignore this.
}

func (m *mockTerminal) SignalInterrupt() {
	m.exitSignaled = true
}

func (m mockTerminal) MakeRaw() (*term.State, error) {
	return nil, nil
}

type stdinDripFeeder struct {
	inputs [][]byte
	index  int
}

func (s *stdinDripFeeder) Read(b []byte) (int, error) {
	copy(b, s.inputs[s.index])
	l := len(s.inputs[s.index])
	s.index++
	return l, nil
}

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
				{3},
			},
			shouldExit: true,
			items: []string{
				"hello", "world",
			},
		},
		{
			name: "scroll down on non-row selection menu",
			inputs: [][]byte{
				{27, 91, 66},
				{3},
			},
			shouldExit: true,
			items: []string{
				"hello", "world",
			},
		},
		{
			name: "wrap around on non-row selection menu",
			inputs: [][]byte{
				{27, 91, 66},
				{27, 91, 66},
				{3},
			},
			shouldExit: true,
			items: []string{
				"hello", "world",
			},
		},
		{
			name: "scroll up on non-row selection menu",
			inputs: [][]byte{
				{27, 91, 66},
				{27, 91, 65},
				{3},
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
				{3},
			},
			shouldExit: true,
			items: []string{
				"hello", "world",
			},
		},
		{
			name: "make selection on non-row selection menu",
			inputs: [][]byte{
				{27, 91, 66},
				{13},
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
				{3},
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
				{27, 91, 66},
				{3},
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
				{27, 91, 66},
				{27, 91, 66},
				{3},
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
				{27, 91, 66},
				{27, 91, 65},
				{3},
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
				{3},
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
				{27, 91, 66},
				{13},
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
				{3},
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
				{27, 91, 66},
				{3},
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
				{27, 91, 66},
				{27, 91, 66},
				{3},
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
				{27, 91, 66},
				{27, 91, 65},
				{3},
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
				{3},
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
				{13},
				{27, 91, 66},
				{13},
				{27, 91, 66},
				{27},
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
				{3},
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
				{27, 91, 66},
				{3},
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
				{27, 91, 66},
				{27, 91, 66},
				{3},
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
				{27, 91, 66},
				{27, 91, 65},
				{3},
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
				{3},
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
				{13},
				{27, 91, 66},
				{13},
				{27, 91, 66},
				{27},
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
			stdin := &stdinDripFeeder{inputs: tt.inputs}
			stdout := &mockTerminal{buf: &bytes.Buffer{}}
			res := selectorComponent("test", tt.columns, tt.items, stdin, tt.multiple, stdout)
			if tt.shouldExit {
				assert.Equal(t, tt.shouldExit, stdout.exitSignaled)
			} else {
				assert.Equal(t, tt.result, res)
			}
			buf := stdout.buf
			if golden.Update() {
				golden.Set(t, buf.Bytes())
				return
			}
			assert.Equal(t, string(golden.Get(t)), buf.String())
		})
	}
}
