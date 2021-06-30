package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

type versionInfo struct {
	Version   string
	Commit    string
	Date      string
	populated bool
}

func (v *versionInfo) Populate() {
	if v.populated {
		return
	}

	if Version != "" {
		v.Version = Version
	} else {
		v.Version = "undefined"
	}

	if Commit != "" {
		v.Commit = Commit
	} else {
		v.Commit = "undefined"
	}

	if Date != "" {
		ts, err := strconv.Atoi(Date)
		if err == nil {
			v.Date = time.Unix(int64(ts), 0).UTC().String()
		}
	}
	if v.Date == "" {
		v.Date = "undefined"
	}

	v.populated = true
}

var (
	Version string
	Commit  string
	Date    string

	prettyVersion = &versionInfo{}
	versionCmd    = &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Long:  `Print the version number of katapult CLI tool.`,
		Run: func(cmd *cobra.Command, args []string) {
			prettyVersion.Populate()

			fmt.Printf("katapult %s (katapult-cli)\n", prettyVersion.Version)
			fmt.Println("---")
			fmt.Printf("Version: %s\n", prettyVersion.Version)
			fmt.Printf("GitCommit: %s\n", prettyVersion.Commit)
			fmt.Printf("BuildDate: %s\n", prettyVersion.Date)
		},
	}
)
