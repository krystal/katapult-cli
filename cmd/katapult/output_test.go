package main

import (
	"bytes"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenericOutput_JSON(t *testing.T) {
	g := genericOutput{item: "hello world!"}
	b, err := g.JSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"hello world!"`), []byte(b))
}

func TestGenericOutput_YAML(t *testing.T) {
	g := genericOutput{item: []string{"hello", "world"}}
	s, err := g.YAML()
	assert.NoError(t, err)
	assert.Equal(t, getTestData(t, "output.yaml"), s)
}

func Test_renderTemplate(t *testing.T) {
	s, err := renderTemplate("{{.}}", "Hello World!")
	assert.NoError(t, err)
	assert.Equal(t, "Hello World!", s)
}

func TestGenericOutput_Text(t *testing.T) {
	tests := []struct {
		name string

		item            interface{}
		defaultTemplate string
		templateArg     string
		wants           string
	}{
		{
			name:            "default template",
			item:            map[string]int{"a": 1, "b": 2},
			defaultTemplate: getTestData(t, "for_map_tpl.txt"),
			wants:           "a1b2\n",
		},
		{
			name:        "template override",
			item:        map[string]int{"a": 1, "b": 2},
			templateArg: getTestData(t, "for_map_tpl.txt"),
			wants:       "a1b2\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := genericOutput{
				item: tt.item,
				tpl:  tt.defaultTemplate,
			}
			text, err := g.Text(tt.templateArg)
			assert.NoError(t, err)
			assert.Equal(t, text, tt.wants)
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
		wants           string
		wantErr         string
	}{
		{
			name:            "default template",
			item:            map[string]int{"a": 1, "b": 2},
			defaultTemplate: getTestData(t, "for_map_tpl.txt"),
			wants:           "a1b2\n",
		},
		{
			name:         "template override",
			item:         map[string]int{"a": 1, "b": 2},
			templateFlag: getTestData(t, "for_map_tpl.txt"),
			wants:        "a1b2\n",
		},
		{
			name:       "json flag",
			item:       map[string]int{"a": 1, "b": 2},
			outputType: "json",
			wants:     	getTestData(t, "flag.json"),
		},
		{
			name:       "yaml flag",
			item:       map[string]int{"a": 1, "b": 2},
			outputType: "yaml",
			wants: 		getTestData(t, "flag.yaml"),
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
			wrapper := renderOption(func(*cobra.Command, []string) (Output, error) {
				if tt.throws != "" {
					return nil, errors.New(tt.throws)
				}
				return genericOutput{
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
			assert.Equal(t, tt.wants, buf.String())
		})
	}
}
