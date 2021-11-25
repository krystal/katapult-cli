package console

import (
	"os"
	"sync"

	"github.com/buger/goterm"
	"golang.org/x/term"
)

type gotermTerminal struct {
	raw *term.State
	m   sync.Mutex
}

func (*gotermTerminal) Height() int {
	return goterm.Height()
}

func (*gotermTerminal) Width() int {
	return goterm.Width()
}

func (g *gotermTerminal) Print(items ...interface{}) (int, error) {
	g.m.Lock()
	defer g.m.Unlock()
	if g.raw != nil {
		_ = term.Restore(0, g.raw)
	}
	n, err := goterm.Print(items...)
	if err == nil && g.raw != nil {
		g.raw, err = term.MakeRaw(0)
		if err != nil {
			return 0, err
		}
	}
	return n, err
}

func (g *gotermTerminal) Sync(s string) error {
	// TODO: There are better ways of doing this. I took a few hours to attempt to resync the screen line by line.
	// TODO-1: However, by piping the whole chunk out, I find this stops the tearing a lot.
	g.m.Lock()
	defer g.m.Unlock()
	if g.raw != nil {
		_ = term.Restore(0, g.raw)
	}
	_, err := os.Stdout.Write([]byte(s))
	if err == nil && g.raw != nil {
		g.raw, err = term.MakeRaw(0)
		if err != nil {
			return err
		}
	}
	return err
}

func (g *gotermTerminal) Clear() {
	g.m.Lock()
	defer g.m.Unlock()
	if g.raw != nil {
		_ = term.Restore(0, g.raw)
	}
	goterm.Clear()
	if g.raw != nil {
		g.raw, _ = term.MakeRaw(0)
	}
}

func (g *gotermTerminal) Println(items ...interface{}) (int, error) {
	g.m.Lock()
	defer g.m.Unlock()
	if g.raw != nil {
		_ = term.Restore(0, g.raw)
	}
	n, err := goterm.Println(items...)
	if err == nil && g.raw != nil {
		g.raw, err = term.MakeRaw(0)
		if err != nil {
			return 0, err
		}
	}
	return n, err
}

func (g *gotermTerminal) Flush() {
	g.m.Lock()
	defer g.m.Unlock()
	if g.raw != nil {
		_ = term.Restore(0, g.raw)
	}
	goterm.Flush()
	if g.raw != nil {
		g.raw, _ = term.MakeRaw(0)
	}
}

func (g *gotermTerminal) SignalInterrupt() {
	g.m.Lock()
	if g.raw != nil {
		_ = term.Restore(0, g.raw)
		g.raw = nil
	}
	g.m.Unlock()
	os.Exit(1)
}

func (g *gotermTerminal) MakeRaw() error {
	g.m.Lock()
	raw, err := term.MakeRaw(0)
	if raw != nil {
		g.raw = raw
	}
	g.m.Unlock()
	return err
}

func (g *gotermTerminal) Unraw() error {
	var err error
	g.m.Lock()
	if g.raw != nil {
		err = term.Restore(0, g.raw)
		g.raw = nil
	}
	g.m.Unlock()
	return err
}

func (*gotermTerminal) BufferInputs() bool {
	return true
}
