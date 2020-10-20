package main

import (
	"context"
	"fmt"
	"os"

	"github.com/krystal/go-katapult"
	"github.com/krystal/katapult-cli/internal/client"
	"github.com/krystal/katapult-cli/pkg/config"
	"github.com/spf13/cobra"
)

type contextKey int

const (
	ctxKey contextKey = iota
)

var (
	// flags
	configFileFlag string
	configURLFlag  string
	configAPIKey   string

	ctx context.Context
	cl  *katapult.Client

	conf    = config.New()
	rootCmd = &cobra.Command{
		Use:   "katapult",
		Short: "katapult CLI tool",
		Long:  `katapult is a CLI tool for the katapult.io hosting platform.`,
	}
)

func init() {
	cobra.OnInitialize(initConfig, initClient)

	rootFlags := rootCmd.PersistentFlags()

	rootFlags.StringVarP(&configFileFlag, "config", "c", "",
		"config file (default: $HOME/.katapult/katapult.yaml)")

	rootFlags.StringVar(&configURLFlag, "api-url", "", fmt.Sprintf(
		"URL for Katapult API (default: %s)", config.Defaults.APIURL,
	))
	conf.BindPFlag("api_url", rootFlags.Lookup("api-url"))

	rootFlags.StringVar(&configAPIKey, "api-key", "", fmt.Sprintf(
		"Katapult API Key (default: %s)", config.Defaults.APIKey,
	))
	conf.BindPFlag("api_key", rootFlags.Lookup("api-key"))
}

func er(msg interface{}) {
	fmt.Fprintln(os.Stderr, "Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if configFileFlag != "" {
		conf.SetConfigFile(configFileFlag)
	}

	_ = conf.Load()
}

func initClient() {
	ctx = context.WithValue(context.Background(), ctxKey, "katapult")
	cl = client.New(ctx, conf)
}

func main() {
	rootCmd.Execute()
}
