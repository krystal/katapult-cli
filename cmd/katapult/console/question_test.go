package console

import (
	"bytes"
	"errors"
	"strconv"
	"testing"

	"github.com/krystal/katapult-cli/internal/golden"
	"github.com/stretchr/testify/assert"
)

type questionStdin struct {
	count int
}

func (q *questionStdin) ReadString(delim byte) (string, error) {
	if delim != '\n' {
		return "", errors.New("delim is expected to be new line")
	}
	old := q.count
	q.count++
	if old == 0 {
		return "\n", nil
	}
	return strconv.Itoa(old)+"\n", nil
}

func TestQuestion(t *testing.T) {
	tests := []struct{
		name string

		blank bool
		response string
	} {
		{
			name:  "handle blank input",
			blank: true,
		},
		{
			name:  "handle ignoring blank input",
			response: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			assert.Equal(t,
				Question("test", tt.blank, &questionStdin{}, stdout),
				tt.response)
			if golden.Update() {
				golden.Set(t, stdout.Bytes())
				return
			}
			assert.Equal(t, string(golden.Get(t)), stdout.String())
		})
	}
}
