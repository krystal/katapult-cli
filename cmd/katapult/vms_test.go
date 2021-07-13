
package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const stdoutVMsList = `- loge-enthusiasts-1 (example.com) [some_id]: ROCK-6
- loge-enthusiasts-2 (2.example.com) [some_id_2]: ROCK-6`

func TestVMs_List(t *testing.T) {
	cmd := virtualMachinesCmd(mockAPIClient{})
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"list"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, stdoutVMsList, stdout.String())
}

func TestVMs_Poweroff(t *testing.T) {
	cmd := virtualMachinesCmd(mockAPIClient{})
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"poweroff", "--id", "1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, stdoutVMsList, stdout.String())
}

func TestVMs_Stop(t *testing.T) {
	cmd := virtualMachinesCmd(mockAPIClient{})
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"stop", "--id", "1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, stdoutVMsList, stdout.String())
}

func TestVMs_Start(t *testing.T) {
	cmd := virtualMachinesCmd(mockAPIClient{})
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"start", "--id", "1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, stdoutVMsList, stdout.String())
}

func TestVMs_Reset(t *testing.T) {
	cmd := virtualMachinesCmd(mockAPIClient{})
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"reset", "--id", "1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, stdoutVMsList, stdout.String())
}
