package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var outputFlag, templateFlag string

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

// Used to render a table.
func table(columns []string, rows [][]interface{}) string {
	buf := &bytes.Buffer{}
	t := tablewriter.NewWriter(buf)
	t.SetHeader(columns)
	strrows := make([][]string, len(rows))
	for i, row := range rows {
		strrow := make([]string, len(row))
		for x, v := range row {
			strrow[x] = fmt.Sprint(v)
		}
		strrows[i] = strrow
	}
	t.AppendBulk(strrows)
	t.Render()
	return buf.String()
}

// Handle mapping a KV to an array.
func kvMap(m map[string]interface{}) [][]interface{} {
	a := make([][]interface{}, len(m))
	i := 0
	for k, v := range m {
		a[i] = []interface{}{k, v}
		i++
	}
	return a
}

// Used to make a string slice.
func stringSlice(items ...string) []string {
	return items
}

// Used to return a single row.
func singleRow(items ...interface{}) [][]interface{} {
	return [][]interface{}{items}
}

// Used to return multiple rows.
func multipleRows(items interface{}, keys ...string) [][]interface{} {
	a := make([][]interface{}, len(keys))
	itemsReflect := reflect.ValueOf(items)
	for i := 0; i < itemsReflect.Len(); i++ {
		x := make([]interface{}, len(keys))
		value := reflect.Indirect(itemsReflect.Index(i))
		for i, k := range keys {
			dotsplit := strings.Split(k, ".")
			for x := 0; x < len(dotsplit)-1; x++ {
				value = reflect.Indirect(value.FieldByName(dotsplit[x]))
			}
			x[i] = value.FieldByName(dotsplit[len(dotsplit)-1]).Interface()
		}
		a[i] = x
	}
	return a
}

// Used to render the template.
func renderTemplate(tpl string, data interface{}) (string, error) {
	parsed, err := template.New("tpl").Funcs(template.FuncMap{
		"Table":        table,
		"KVMap":        kvMap,
		"StringSlice":  stringSlice,
		"SingleRow":    singleRow,
		"MultipleRows": multipleRows,
	}).Parse(tpl)
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

// Defines a function that returns a output.
type outputFunc func(cmd *cobra.Command, args []string) (Output, error)

// Used to render a console output of a single option. Passes through errors.
func renderOption(f outputFunc) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Get the objects we require from the cache.
		out := cmd.OutOrStdout()

		// Call the function.
		output, err := f(cmd, args)
		if err != nil {
			return err
		}

		// Handle calling the correct render function.
		switch strings.ToLower(outputFlag) {
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
			s, err := output.Text(templateFlag)
			if err != nil {
				return err
			}
			_, _ = out.Write([]byte(s))
			return nil
		}
	}
}
