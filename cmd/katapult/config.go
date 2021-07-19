package main

import (
	"fmt"

	"github.com/krystal/katapult-cli/config"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func configCommand(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Print configuration",
		Long:  `Print parsed configuration in YAML format.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			bs, err := yaml.Marshal(conf.AllSettings())
			if err != nil {
				return fmt.Errorf("unable to marshal config to YAML: %w", err)
			}

			out := cmd.OutOrStdout()
			_, _ = fmt.Fprintln(out, "---")
			_, _ = out.Write(append(bs, '\n'))
			return nil
		},
	}
}
