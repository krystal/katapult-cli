package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/krystal/katapult-cli/internal/client"
	"github.com/krystal/katapult-cli/pkg/config"
	"github.com/spf13/cobra"
)

func run() error {
	var (
		configFileFlag string
		configURLFlag  string
		configAPIKey   string
	)

	conf := config.New()
	if configFileFlag != "" {
		conf.SetConfigFile(configFileFlag)
	}

	err := conf.Load()
	if err != nil {
		return err
	}

	ctx := context.Background()
	cl := client.New(ctx, conf)
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

	rootCmd.AddCommand(
		versionCommand(),
		configCommand(conf),
		dataCentersCmd(cl),
	)

	return rootCmd.Execute()
}

func main() {
	err := run()
	if err != nil {
		log.Printf("A fatal error occured: %s", err)
		os.Exit(1)
	}
}