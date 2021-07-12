package main

import (
	"fmt"
	"github.com/krystal/go-katapult"
	"log"
	"os"

	"github.com/krystal/katapult-cli/config"
	"github.com/spf13/cobra"
)

func run() error {
	var (
		configFileFlag string
		configURLFlag  string
		configAPIKey   string
	)

	conf, err := config.New()
	if err != nil {
		return err
	}
	if configFileFlag != "" {
		conf.SetConfigFile(configFileFlag)
	}

	err = conf.Load()
	if err != nil {
		return err
	}

	var cl *katapult.Client
	rootCmd := &cobra.Command{
		Use:   "katapult",
		Short: "katapult CLI tool",
		Long:  `katapult is a CLI tool for the katapult.io hosting platform.`,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// TODO: Pointer problems here because the commands are built before here.
			cl, err = newClient(conf)
			return err
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configFileFlag, "config", "c", "",
		"config file (default: $HOME/.katapult/katapult.yaml)")

	rootCmd.PersistentFlags().StringVar(&configURLFlag, "api-url", "", fmt.Sprintf(
		"URL for Katapult API (default: %s)", config.Defaults.APIURL,
	))
	err = conf.BindPFlag("api_url", rootCmd.PersistentFlags().Lookup("api-url"))
	if err != nil {
		return err
	}
	rootCmd.PersistentFlags().StringVar(&configAPIKey, "api-key", "", fmt.Sprintf(
		"Katapult API Key (default: %s)", config.Defaults.APIKey,
	))
	err = conf.BindPFlag("api_key", rootCmd.PersistentFlags().Lookup("api-key"))
	if err != nil {
		return err
	}

	rootCmd.AddCommand(
		versionCommand(),
		configCommand(conf),
		dataCentersCmd(cl),
		networksCmd(cl),
		organizationsCmd(cl),
		virtualMachinesCmd(cl),
	)

	return rootCmd.Execute()
}

func main() {
	err := run()
	if err != nil {
		log.Printf("A fatal error occurred: %s", err)
		os.Exit(1)
	}
}
