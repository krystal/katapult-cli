package main

import (
	"github.com/krystal/go-katapult/core"
	"github.com/krystal/katapult-cli/config"
	"github.com/spf13/cobra"
)

func authCommand(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "auth",
		Args:  cobra.ExactArgs(1),
		Short: "Authenticate user",
		Long:  "Authenticates the user with a token. The argument should be an API token.",
		RunE: outputWrapper(func(cmd *cobra.Command, args []string) (Output, error) {
			// Get the token.
			token := args[0]

			// Set the token in the config. This is okay to mutate since we are going
			// to exit after this in any case.
			conf.APIToken = token

			// Check if the config works.
			c, err := newClient(conf)
			if err != nil {
				return nil, err
			}
			if _, _, err = core.NewDataCentersClient(c).List(cmd.Context()); err != nil {
				return nil, err
			}

			// Write the config.
			conf.SetDefault("api_token", token)
			if err = conf.WriteConfig(); err != nil {
				return nil, err
			}

			// Return the output.
			return &genericOutput{
				item:                map[string]bool{"success": true},
				defaultTextTemplate: "Successfully authenticated.\n",
			}, nil
		}),
	}
}
