package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var outputFlag, templateFlag string

// Output is used to define the interface of outputs.
type Output interface {
	// JSON is used to write out the JSON output.
	JSON(w io.Writer) error

	// YAML is used to write out the YAML output.
	YAML(w io.Writer) error

	// Text is used to render a template. If string is blank, uses the default.
	Text(w io.Writer, template string) error
}

// Used to render a table.
func table(columns []string, rows [][]interface{}) string {
	buf := &bytes.Buffer{}
	t := tablewriter.NewWriter(buf)

	// This is based off the usage for tablewriter
	// source: https://github.com/olekukonko/tablewriter#example-10---set-nowhitespace-and-tablepadding-option
	t.SetAutoWrapText(false)
	t.SetAutoFormatHeaders(true)
	t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetCenterSeparator("")
	t.SetColumnSeparator("")
	t.SetRowSeparator("")
	t.SetHeaderLine(false)
	t.SetBorder(false)
	t.SetTablePadding("\t")
	t.SetNoWhiteSpace(true)
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
	orderedKeys := make([]string, len(m))
	for k := range m {
		orderedKeys[i] = k
		i++
	}
	sort.Strings(orderedKeys)
	i = 0
	for _, k := range orderedKeys {
		v := m[k]
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
	// Use reflect to get the items.
	// We are using this with the slice since we might need to handle many different slice types.
	itemsReflect := reflect.ValueOf(items)

	// Create a slice of all of the rows.
	rows := make([][]interface{}, itemsReflect.Len())

	// Go through each length items.
	for i := 0; i < itemsReflect.Len(); i++ {
		// Create the row.
		row := make([]interface{}, len(keys))

		// Get the value using reflect so we can access fields.
		outerValue := reflect.Indirect(itemsReflect.Index(i))

		// Go through each key which we want from the field.
		for i, k := range keys {
			// Get the locally scoped value.
			value := outerValue

			// Split by dots so we can get properties.
			dotsplit := strings.Split(k, ".")

			// Traverse through each field in the key. len-1 is safe here since split
			// will always return at least 1 item.
			for x := 0; x < len(dotsplit)-1; x++ {
				value = reflect.Indirect(value.FieldByName(dotsplit[x]))
			}

			// Get the item from the struct.
			row[i] = value.FieldByName(dotsplit[len(dotsplit)-1]).Interface()
		}

		// Add to the array.
		rows[i] = row
	}

	// Return the rows array.
	return rows
}

// Used to render the template.
func renderTemplate(w io.Writer, tpl string, data interface{}) error {
	parsed, err := template.New("tpl").Funcs(template.FuncMap{
		"Table":        table,
		"KVMap":        kvMap,
		"StringSlice":  stringSlice,
		"SingleRow":    singleRow,
		"MultipleRows": multipleRows,
	}).Parse(tpl)
	if err != nil {
		return err
	}
	return parsed.Execute(w, data)
}

// Used to implement Output for a variety of use cases.
type genericOutput struct {
	item                interface{}
	defaultTextTemplate string
}

// JSON is used to write out the JSON output.
func (g *genericOutput) JSON(w io.Writer) error {
	b, err := json.MarshalIndent(g.item, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// YAML is used to write out the YAML output.
func (g *genericOutput) YAML(w io.Writer) error {
	b, err := yaml.Marshal(g.item)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// Text is used to render a template. If string is blank, uses the default.
func (g *genericOutput) Text(w io.Writer, template string) error {
	if template == "" {
		// Return the default template.
		template = g.defaultTextTemplate
	}

	// Render the template.
	return renderTemplate(w, template, g.item)
}

// Defines a function that returns a output.
type outputFunc func(cmd *cobra.Command, args []string) (Output, error)

// Used to render a console output of a type. Passes through errors.
func outputWrapper(f outputFunc) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Get stdout.
		out := cmd.OutOrStdout()

		// Call the function.
		output, err := f(cmd, args)
		if err != nil {
			return err
		}

		// Handle calling the correct render function.
		switch strings.ToLower(outputFlag) {
		case "json":
			return output.JSON(out)
		case "yml", "yaml":
			return output.YAML(out)
		default:
			return output.Text(out, templateFlag)
		}
	}
}
