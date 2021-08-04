package main

import (
	"github.com/krystal/katapult-cli/config"
	"github.com/spf13/cobra"
)

const configFormat = "{{ range $key, $value := . }}" +
	"{{ $key }}: {{ $value }}\n{{ end }}"

func configCommand(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Print configuration",
		Long:  "Print parsed configuration in YAML/JSON format.",
		RunE: renderOption(func(cmd *cobra.Command, args []string) (Output, error) {
			return genericOutput{
				item: conf.AllSettings(),
				tpl:  configFormat,
			}, nil
		}),
	}
}
