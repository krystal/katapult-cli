package main

import (
	"os"

	"github.com/krystal/go-katapult/core"
	"github.com/krystal/katapult-cli/config"
	"github.com/spf13/cobra"
)

func run() error {
	var (
		configFileFlag string
		configURLFlag  string
		configAPIToken string
	)

	conf, err := config.New()
	if err != nil {
		return err
	}

	var help bool
	rootCmd := &cobra.Command{
		Use:   "katapult",
		Short: "katapult CLI tool",
		Long:  `katapult is a CLI tool for the katapult.io hosting platform.`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if help {
				err = cmd.Usage()
				if err != nil {
					return err
				}
				os.Exit(0)
			}
			return nil
		},
		SilenceUsage: true,
	}

	rootFlags := rootCmd.PersistentFlags()

	rootFlags.BoolVarP(&help, "help", "h", false, "Display the help for the command/root.")

	rootFlags.StringVarP(&outputFlag, "output", "o", "", "output type (yaml, json, text)")
	rootFlags.StringVar(&templateFlag, "format", "", "defines the output template for text")

	rootFlags.StringVar(&configFileFlag, "config-path", "",
		"config file (default: $HOME/.katapult/katapult.yaml)")

	apiURLDefault := config.Defaults.APIURL
	if apiURLDefault != "" {
		apiURLDefault = " (default: " + apiURLDefault + ")"
	}
	rootFlags.StringVar(&configURLFlag, "api-url", "", "URL for Katapult API"+apiURLDefault)
	err = conf.BindPFlag("api_url", rootFlags.Lookup("api-url"))
	if err != nil {
		return err
	}
	tokenDefault := config.Defaults.APIToken
	if tokenDefault != "" {
		tokenDefault = " (default: " + tokenDefault + ")"
	}
	rootFlags.StringVar(&configAPIToken, "api-token", "", "Katapult API Token"+tokenDefault)
	err = conf.BindPFlag("api_token", rootFlags.Lookup("api-token"))
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

	if err = conf.Load(); err != nil {
		return err
	}

	cl, err := newClient(conf)
	if err != nil {
		return err
	}

	rootCmd.AddCommand(
		authCommand(conf),
		versionCommand(),
		configCommand(conf),
		dataCentersCmd(core.NewDataCentersClient(cl)),
		networksCmd(core.NewNetworksClient(cl)),
		organizationsCmd(core.NewOrganizationsClient(cl)),
		virtualMachinesCmd(
			core.NewVirtualMachinesClient(cl),
			core.NewOrganizationsClient(cl),
			core.NewDataCentersClient(cl),
			core.NewVirtualMachinePackagesClient(cl),
			core.NewDiskTemplatesClient(cl),
			core.NewIPAddressesClient(cl),
			core.NewSSHKeysClient(cl),
			core.NewTagsClient(cl),
			core.NewVirtualMachineBuildsClient(cl),
			nil),
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
