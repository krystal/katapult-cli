package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

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

func versionCommand() *cobra.Command {
	prettyVersion := &versionInfo{}
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Long:  `Print the version number of katapult CLI tool.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			prettyVersion.Populate()

			if strings.ToLower(cmd.Flag("output").Value.String()) == "json" {
				j, err := json.Marshal(prettyVersion)
				if err != nil {
					return err
				}
				fmt.Println(string(j))
			} else {
				fmt.Printf("katapult %s (katapult-cli)\n", prettyVersion.Version)
				fmt.Println("---")
				fmt.Printf("Version: %s\n", prettyVersion.Version)
				fmt.Printf("GitCommit: %s\n", prettyVersion.Commit)
				fmt.Printf("BuildDate: %s\n", prettyVersion.Date)
			}
			return nil
		},
	}

	flags := versionCmd.PersistentFlags()
	flags.StringP("output", "o", "text", "Defines the output type of the config. Can be text or json.")

	return versionCmd
}
