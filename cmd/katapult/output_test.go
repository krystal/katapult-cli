package main

import (
	"bytes"
	"errors"
	"testing"

	"github.com/krystal/katapult-cli/internal/golden"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenericOutput_JSON(t *testing.T) {
	g := genericOutput{item: "hello world!"}
	buf := &bytes.Buffer{}
	err := g.JSON(buf)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"hello world!"`), buf.Bytes())
}

func TestGenericOutput_YAML(t *testing.T) {
	g := genericOutput{item: []string{"hello", "world"}}
	buf := &bytes.Buffer{}
	err := g.YAML(buf)
	assert.NoError(t, err)
	if golden.Update() {
		golden.Set(t, buf.Bytes())
		return
	}
	assert.Equal(t, golden.Get(t), buf.Bytes())
}

func Test_renderTemplate(t *testing.T) {
	buf := &bytes.Buffer{}
	err := renderTemplate(buf, "{{.}}", "Hello World!")
	assert.NoError(t, err)
	assert.Equal(t, "Hello World!", buf.String())
}

func TestGenericOutput_Text(t *testing.T) {
	tests := []struct {
		name string

		item            interface{}
		defaultTemplate string
		templateArg     string
	}{
		{
			name:            "default template",
			item:            map[string]int{"a": 1, "b": 2},
			defaultTemplate: getForMapTpl(t),
		},
		{
			name:        "template override",
			item:        map[string]int{"a": 1, "b": 2},
			templateArg: getForMapTpl(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := genericOutput{
				item: tt.item,
				tpl:  tt.defaultTemplate,
			}
			buf := &bytes.Buffer{}
			err := g.Text(buf, tt.templateArg)
			assert.NoError(t, err)
			if golden.Update() {
				golden.Set(t, buf.Bytes())
				return
			}
			assert.Equal(t, buf.Bytes(), golden.Get(t))
		})
	}
}

func Test_renderOption(t *testing.T) {
	tests := []struct {
		name string

		item            interface{}
		outputType      string
		defaultTemplate string
		templateFlag    string
		throws          string
		wantErr         string
	}{
		{
			name:            "default template",
			item:            map[string]int{"a": 1, "b": 2},
			defaultTemplate: getForMapTpl(t),
		},
		{
			name:         "template override",
			item:         map[string]int{"a": 1, "b": 2},
			templateFlag: getForMapTpl(t),
		},
		{
			name:       "json flag",
			item:       map[string]int{"a": 1, "b": 2},
			outputType: "json",
		},
		{
			name:       "yaml flag",
			item:       map[string]int{"a": 1, "b": 2},
			outputType: "yaml",
		},
		{
			name:    "test throw",
			wantErr: "testing testing 123",
			throws:  "testing testing 123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock cobra command.
			mockCmd := &cobra.Command{}

			// Create a buffer to store the stdout.
			buf := &bytes.Buffer{}
			mockCmd.SetOut(buf)

			// Defines the function wrapper. A tiny basic function is specified to return the test output.
			wrapper := outputWrapper(func(*cobra.Command, []string) (Output, error) {
				if tt.throws != "" {
					return nil, errors.New(tt.throws)
				}
				return &genericOutput{
					item: tt.item,
					tpl:  tt.defaultTemplate,
				}, nil
			})

			// Call the wrapped function and check the error is what we want.
			outputFlag = tt.outputType
			templateFlag = tt.templateFlag
			err := wrapper(mockCmd, nil)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			outputFlag = ""
			templateFlag = ""

			// Check stdout.
			if golden.Update() {
				golden.Set(t, buf.Bytes())
			}
			assert.Equal(t, golden.Get(t), buf.Bytes())
		})
	}
}
