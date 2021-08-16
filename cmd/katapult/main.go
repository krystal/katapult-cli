package main

import (
	"fmt"
	"os"

	"github.com/krystal/go-katapult/core"
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

	rootCmd := &cobra.Command{
		Use:   "katapult",
		Short: "katapult CLI tool",
		Long:  `katapult is a CLI tool for the katapult.io hosting platform.`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		SilenceUsage: true,
	}

	rootFlags := rootCmd.PersistentFlags()

	rootFlags.StringVarP(&outputFlag, "output", "o", "", "output type (yaml, json, text)")
	rootFlags.StringVar(&templateFlag, "format", "", "defines the output template for text")

	rootFlags.StringVarP(&configFileFlag, "config", "c", "",
		"config file (default: $HOME/.katapult/katapult.yaml)")

	rootFlags.StringVar(&configURLFlag, "api-url", "", fmt.Sprintf(
		"URL for Katapult API (default: %s)", config.Defaults.APIURL,
	))
	err = conf.BindPFlag("api_url", rootFlags.Lookup("api-url"))
	if err != nil {
		return err
	}
	rootFlags.StringVar(&configAPIKey, "api-key", "", fmt.Sprintf(
		"Katapult API Key (default: %s)", config.Defaults.APIKey,
	))
	err = conf.BindPFlag("api_key", rootFlags.Lookup("api-key"))
	if err != nil {
		return err
	}

	err = rootCmd.ParseFlags(os.Args)
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

	cl, err := newClient(conf)
	if err != nil {
		return err
	}

	rootCmd.AddCommand(
		versionCommand(),
		configCommand(conf),
		dataCentersCmd(core.NewDataCentersClient(cl)),
		networksCmd(core.NewNetworksClient(cl)),
		organizationsCmd(core.NewOrganizationsClient(cl)),
		virtualMachinesCmd(core.NewVirtualMachinesClient(cl)),
	)

	return rootCmd.Execute()
}

func main() {
	err := run()
	if err != nil {
		// Ensure we exit with status code 1. The actual printing is done by Cobra.
		os.Exit(1)
	}
}
