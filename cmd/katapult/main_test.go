package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertCobraCommand(t *testing.T, cmd *cobra.Command, errResult, want, stderrResult string) {
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

	assert.Equal(t, want, stdout.String())
	assert.Equal(t, stderrResult, stderr.String())
}

func assertCobraCommandReturnStdout(t *testing.T, cmd *cobra.Command, errResult, stderrResult string) string {
	t.Helper()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	err := cmd.Execute()

	if errResult != "" {
		require.EqualError(t, err, errResult)
		return stdout.String()
	}
	require.NoError(t, err)

	assert.Equal(t, stderrResult, stderr.String())
	return stdout.String()
}

type mockOrganizationsListClient struct {
	orgs   []*core.Organization
	throws string
}

func (m mockOrganizationsListClient) List(context.Context) ([]*core.Organization, *katapult.Response, error) {
	if m.throws != "" {
		return nil, nil, errors.New(m.throws)
	}
	return m.orgs, nil, nil
}

type mockDataCentersClient struct {
	dcs    []*core.DataCenter
	throws string
}

func (m mockDataCentersClient) List(context.Context) ([]*core.DataCenter, *katapult.Response, error) {
	if m.throws != "" {
		return nil, nil, errors.New(m.throws)
	}
	return m.dcs, nil, nil
}

func (m mockDataCentersClient) Get(
	_ context.Context, ref core.DataCenterRef) (*core.DataCenter, *katapult.Response, error) {
	if m.throws != "" {
		return nil, nil, errors.New(m.throws)
	}
	for _, v := range m.dcs {
		if v.Permalink == ref.Permalink {
			return v, nil, nil
		}
	}
	return nil, nil, fmt.Errorf("unknown datacentre")
}

var fixtureOrganizations = []*core.Organization{
	{
		ID:        "loge",
		Name:      "Loge Enthusiasts",
		SubDomain: "loge",
	},
	{
		ID:        "testing",
		Name:      "testing, testing, 123",
		SubDomain: "test",
	},
}

var fixtureDataCenters = []*core.DataCenter{
	{
		ID:        "dc_9UVoPiUQoI1cqtRd",
		Name:      "hello",
		Permalink: "POG1",
		Country: &core.Country{
			ID:   "POG",
			Name: "Pogland",
		},
	},
	{
		ID:        "dc_9UVoPiUQoI1cqtR0",
		Name:      "hello",
		Permalink: "GB1",
		Country: &core.Country{
			ID:   "UK",
			Name: "United Kingdom",
		},
	},
}
