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

func Test_genericOutput(t *testing.T) {
	assert.Implements(t, (*Output)(nil), &genericOutput{})
}

func Test_genericOutput_JSON(t *testing.T) {
	tests := []struct {
		name string

		item interface{}
	}{
		{
			name: "string",
			item: "hello world!",
		},
		{
			name: "number",
			item: 21,
		},
		{
			name: "boolean",
			item: true,
		},
		{
			name: "null",
			item: nil,
		},
		{
			name: "array",
			item: []string{"hello", "world"},
		},
		{
			name: "object",
			item: map[string]string{"hello": "world"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &genericOutput{item: tt.item}
			buf := &bytes.Buffer{}
			assert.NoError(t, g.JSON(buf))
			if golden.Update() {
				golden.Set(t, buf.Bytes())
				return
			}
			assert.Equal(t, golden.Get(t), buf.Bytes())
		})
	}
}

func Test_genericOutput_YAML(t *testing.T) {
	tests := []struct {
		name string

		item interface{}
	}{
		{
			name: "string",
			item: "hello world!",
		},
		{
			name: "number",
			item: 21,
		},
		{
			name: "boolean",
			item: true,
		},
		{
			name: "null",
			item: nil,
		},
		{
			name: "array",
			item: []string{"hello", "world"},
		},
		{
			name: "object",
			item: map[string]string{"hello": "world"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &genericOutput{item: tt.item}
			buf := &bytes.Buffer{}
			assert.NoError(t, g.YAML(buf))
			if golden.Update() {
				golden.Set(t, buf.Bytes())
				return
			}
			assert.Equal(t, golden.Get(t), buf.Bytes())
		})
	}
}

func Test_renderTemplate(t *testing.T) {
	buf := &bytes.Buffer{}
	err := renderTemplate(buf, "{{.}}", "Hello World!")
	assert.NoError(t, err)
	assert.Equal(t, "Hello World!", buf.String())
}

func Test_genericOutput_Text(t *testing.T) {
	tests := []struct {
		name string

		item            interface{}
		defaultTemplate string
		templateArg     string
	}{
		{
			name:            "default template",
			item:            map[string]int{"a": 1, "b": 2},
			defaultTemplate: forMapTpl,
		},
		{
			name:        "template override",
			item:        map[string]int{"a": 1, "b": 2},
			templateArg: forMapTpl,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := genericOutput{
				item:                tt.item,
				defaultTextTemplate: tt.defaultTemplate,
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

func Test_outputWrapper(t *testing.T) {
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
			defaultTemplate: forMapTpl,
		},
		{
			name:         "template override",
			item:         map[string]int{"a": 1, "b": 2},
			templateFlag: forMapTpl,
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
					item:                tt.item,
					defaultTextTemplate: tt.defaultTemplate,
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
