package main

import (
	_ "embed"

	"github.com/krystal/katapult-cli/config"
	"github.com/spf13/cobra"
)

//go:embed formatdata/config.txt
var configFormat string

func configCommand(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Print configuration",
		Long:  "Print parsed configuration in YAML/JSON format.",
		RunE: outputWrapper(func(cmd *cobra.Command, args []string) (Output, error) {
			return &genericOutput{
				item:                conf.AllSettings(),
				defaultTextTemplate: configFormat,
			}, nil
		}),
	}
}
