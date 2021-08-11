package main

import (
	"bytes"
	"embed"
	"github.com/krystal/katapult-cli/internal/golden"
	"path"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertCobraCommand(t *testing.T, cmd *cobra.Command, errResult, stderrResult string) {
	t.Helper()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	err := cmd.Execute()

	if errResult != "" {
		require.EqualError(t, err, errResult)
		return
	}
	require.NoError(t, err)

	if golden.Update() {
		golden.Set(t, stdout.Bytes())
		return
	}
	assert.Equal(t, string(golden.Get(t)), stdout.String())
	assert.Equal(t, stderrResult, stderr.String())
}

//go:embed testdata
var testdata embed.FS

func getTestData(t *testing.T, filename string) string {
	t.Helper()
	fp := path.Join("testdata", t.Name(), filename)
	b, err := testdata.ReadFile(fp)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
