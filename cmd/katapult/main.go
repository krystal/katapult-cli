package main

import (
	"fmt"
	"github.com/krystal/go-katapult"
	"log"
	"net/url"
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
	rootCmd := &cobra.Command{
		Use:   "katapult",
		Short: "katapult CLI tool",
		Long:  `katapult is a CLI tool for the katapult.io hosting platform.`,
	}

	rootFlags := rootCmd.PersistentFlags()

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

	cobra.OnInitialize(func() {
		// TODO: We probably want to fix the config to remove this without regressions!
		cl.(*katapult.Client).APIKey = configAPIKey
		if configURLFlag != "" {
			u, err := url.Parse(configURLFlag)
			if err != nil {
				panic(err)
			}
			cl.(*katapult.Client).BaseURL = u
		}
	})

	rootCmd.AddCommand(
		versionCommand(),
		configCommand(conf),
		dataCentersCmd(core.NewDataCentersClient(cl)),
		networksCmd(core.NewNetworksClient(cl)),
		organizationsCmd(core.NewOrganizationsClient(cl)),
		createCmd(core.NewOrganizationsClient(cl),
			core.NewDataCentersClient(cl),
			core.NewVirtualMachinePackagesClient(cl),
			core.NewDiskTemplatesClient(cl),
			core.NewIPAddressesClient(cl),
			core.NewSSHKeysClient(cl),
			core.NewTagsClient(cl),
			core.NewVirtualMachineBuildsClient(cl)))

	return rootCmd.Execute()
}

func main() {
	err := run()
	if err != nil {
		log.Printf("A fatal error occurred: %s", err)
		os.Exit(1)
	}
}
