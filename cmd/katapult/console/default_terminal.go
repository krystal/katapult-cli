package console

import (
	"os"

	"github.com/buger/goterm"
	"golang.org/x/term"
)

type gotermTerminal struct{}

func (gotermTerminal) Height() int {
	return goterm.Height()
}

func (gotermTerminal) Width() int {
	return goterm.Width()
}

func (gotermTerminal) Print(items ...interface{}) (int, error) {
	return goterm.Print(items...)
}

func (gotermTerminal) Clear() {
	goterm.Clear()
}

func (gotermTerminal) Println(items ...interface{}) (int, error) {
	return goterm.Println(items...)
}

func (gotermTerminal) Flush() {
	goterm.Flush()
}

func (gotermTerminal) SignalInterrupt() {
	os.Exit(1)
}

func (gotermTerminal) MakeRaw() (*term.State, error) {
	return term.MakeRaw(0)
}
