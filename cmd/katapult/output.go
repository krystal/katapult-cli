package main

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/cobra"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

// Output is used to define the interface of outputs.
type Output interface {
	// JSON is used to return a string of the JSON output.
	JSON() (json.RawMessage, error)

	// YAML is used to return a string of the YAML output.
	YAML() (string, error)

	// Text is used to render a template. If string is blank, uses the default.
	Text(template string) (string, error)
}

// Used when an item can just output YAML for the user readable output.
type genericOutput struct {
	item interface{}
	tpl  string
}

// JSON is used to return a string of the JSON output.
func (g genericOutput) JSON() (json.RawMessage, error) {
	b, err := json.Marshal(g.item)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// YAML is used to return a string of the YAML output.
func (g genericOutput) YAML() (string, error) {
	b, err := yaml.Marshal(g.item)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Used to render the template.
func renderTemplate(tpl string, data interface{}) (string, error) {
	parsed, err := template.New("tpl").Parse(tpl)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err = parsed.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Text is used to render a template. If string is blank, uses the default.
func (g genericOutput) Text(template string) (string, error) {
	if template == "" {
		// Return the default template.
		template = g.tpl
	}

	// If not, render the template.
	return renderTemplate(template, g.item)
}

// Used to render a console output of a single option. Passes through errors.
func renderOption(f func(cmd *cobra.Command, args []string) (Output, error)) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		output, err := f(cmd, args)
		if err != nil {
			return err
		}
		flags := cmd.PersistentFlags()
		out := cmd.OutOrStdout()
		renderType, err := flags.GetString("output")
		if err != nil {
			return err
		}
		switch strings.ToLower(renderType) {
		case "json":
			b, err := output.JSON()
			if err != nil {
				return err
			}
			return json.NewEncoder(out).Encode(b)
		case "yml", "yaml":
			b, err := output.YAML()
			if err != nil {
				return err
			}
			return yaml.NewEncoder(out).Encode(b)
		default:
			tpl, err := flags.GetString("template")
			if err != nil {
				return err
			}
			s, err := output.Text(tpl)
			if err != nil {
				return err
			}
			_, _ = out.Write([]byte(s))
			return nil
		}
	}
}
