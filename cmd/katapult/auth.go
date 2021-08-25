package main

import (
	"github.com/krystal/go-katapult/core"
	"github.com/krystal/katapult-cli/config"
	"github.com/spf13/cobra"
)

const authFormat = `Successfully authenticated. Here is your current list of organizations:
{{ Table (StringSlice "Name" "Subdomain") (MultipleRows . "Name" "SubDomain") }}`

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
			orgs, _, err := core.NewOrganizationsClient(c).List(cmd.Context())
			if err != nil {
				return nil, err
			}

			// Write the config.
			conf.SetDefault("api_token", token)
			if err = conf.WriteConfig(); err != nil {
				return nil, err
			}

			// Return the output.
			return &genericOutput{
				item:                orgs,
				defaultTextTemplate: authFormat,
			}, nil
		}),
	}
}
