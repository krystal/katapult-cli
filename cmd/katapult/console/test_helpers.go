package console

import (
	"bytes"
	"fmt"

	"golang.org/x/term"
)

// TerminalInterface defines a interface for a compatible terminal. Used for unit testing.
type TerminalInterface interface {
	Height() int
	Width() int
	Print(items ...interface{}) (int, error)
	Println(items ...interface{}) (int, error)
	Clear()
	Flush()
	SignalInterrupt()
	MakeRaw() (*term.State, error)
}

// MockTerminal is used to define a terminal mock for unit tests.
type MockTerminal struct {
	Buffer bytes.Buffer

	CustomWidth  int
	ExitSignaled bool
}

// Height implements TerminalInterface.
func (m *MockTerminal) Height() int {
	return 10
}

// Width implements TerminalInterface.
func (m *MockTerminal) Width() int {
	if m.CustomWidth == 0 {
		return 200
	}
	return m.CustomWidth
}

// Print implements TerminalInterface.
func (m *MockTerminal) Print(items ...interface{}) (int, error) {
	return fmt.Fprint(&m.Buffer, items...)
}

// Println implements TerminalInterface.
func (m *MockTerminal) Println(items ...interface{}) (int, error) {
	return fmt.Fprintln(&m.Buffer, items...)
}

// Clear implements TerminalInterface.
func (m *MockTerminal) Clear() {
	_, _ = fmt.Fprint(&m.Buffer, "\033[2J")
}

// Flush implements TerminalInterface.
func (m *MockTerminal) Flush() {
	// This is platform dependant. Ignore this.
}

// SignalInterrupt implements TerminalInterface.
func (m *MockTerminal) SignalInterrupt() {
	m.ExitSignaled = true
}

// MakeRaw implements TerminalInterface.
func (m *MockTerminal) MakeRaw() (*term.State, error) {
	return nil, nil
}

// StdinDripFeeder is used to define a io.Reader designed to drip feed in different inputs.
type StdinDripFeeder struct {
	Inputs [][]byte
	Index  int
}

// Read implements io.Reader.
func (s *StdinDripFeeder) Read(b []byte) (int, error) {
	copy(b, s.Inputs[s.Index])
	l := len(s.Inputs[s.Index])
	s.Index++
	return l, nil
}
