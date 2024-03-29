package console

import (
	"io"
	"strings"

	"github.com/buger/goterm"
)

// Defines the string reader interface.
type stringReader interface {
	ReadString(delim byte) (string, error)
}

// Question is used to define a basic console question.
func Question(question string, blankAcceptable bool, bufferedStdin stringReader, stdout io.Writer) string {
	for {
		// Print the question.
		_, _ = stdout.Write([]byte(goterm.Color(question, goterm.CYAN) + " "))

		// Read stdin.
		text, _ := bufferedStdin.ReadString('\n')
		text = text[:len(text)-1]
		if text != "" || blankAcceptable {
			return strings.TrimSuffix(text, "\r")
		}
	}
}
