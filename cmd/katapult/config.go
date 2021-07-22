package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/krystal/katapult-cli/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func configCommand(conf *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Print configuration",
		Long:  "Print parsed configuration in YAML/JSON format.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var bs []byte
			var err error
			if strings.ToLower(cmd.Flag("output").Value.String()) == jsonOutput {
				bs, err = json.Marshal(conf.AllSettings())
				if err != nil {
					return fmt.Errorf("unable to marshal config to JSON: %w", err)
				}
			} else {
				bs, err = yaml.Marshal(conf.AllSettings())
				if err != nil {
					return fmt.Errorf("unable to marshal config to YAML: %w", err)
				}
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "---")
			}

			_, _ = cmd.OutOrStdout().Write(append(bs, '\n'))
			return nil
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringP("output", "o", "yaml", "Defines the output type of the config. Can be yaml or json.")

	return cmd
}
