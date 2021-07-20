package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func testAssertCommand(t *testing.T, cmd *cobra.Command, errResult, want, stderrResult string) {
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

	assert.Equal(t, want, stdout.String())
	assert.Equal(t, stderrResult, stderr.String())
}
