package main

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func executeTestCommand(t *testing.T, cmd *cobra.Command, errResult, wants, stderrResult string) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	err := cmd.Execute()
	switch {
	case err == nil:
		// Ignore this.
	case errResult != "":
		assert.Equal(t, errResult, err.Error())
		return
	default:
		t.Fatal(err)
	}
	assert.Equal(t, wants, stdout.String())
	assert.Equal(t, stderrResult, stderr.String())
}
