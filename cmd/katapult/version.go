package main

import (
	"strconv"
	"time"

	_ "embed"

	"github.com/spf13/cobra"
)

var (
	Version string
	Commit  string
	Date    string
)

const unknownPlaceholder = "undefined"

type versionInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	populated bool
}

func (v *versionInfo) Populate() {
	if v.populated {
		return
	}

	if Version != "" {
		v.Version = Version
	} else {
		v.Version = unknownPlaceholder
	}

	if Commit != "" {
		v.Commit = Commit
	} else {
		v.Commit = unknownPlaceholder
	}

	if Date != "" {
		ts, err := strconv.Atoi(Date)
		if err == nil {
			v.Date = time.Unix(int64(ts), 0).UTC().String()
		}
	}
	if v.Date == "" {
		v.Date = unknownPlaceholder
	}

	v.populated = true
}

//go:embed formatdata/version.txt
var versionFormat string

func versionCommand() *cobra.Command {
	prettyVersion := &versionInfo{}
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Long:  `Print the version number of katapult CLI tool.`,
		RunE: outputWrapper(func(cmd *cobra.Command, args []string) (Output, error) {
			prettyVersion.Populate()

			return &genericOutput{
				item:                prettyVersion,
				defaultTextTemplate: versionFormat,
			}, nil
		}),
	}

	return versionCmd
}
