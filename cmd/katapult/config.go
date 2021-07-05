package main

import (
	"fmt"
	"github.com/krystal/katapult-cli/pkg/config"
	"log"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func configCommand(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Print configuration",
		Long:  `Print parsed configuration in YAML format.`,
		Run: func(cmd *cobra.Command, args []string) {
			bs, err := yaml.Marshal(conf.AllSettings())
			if err != nil {
				log.Fatalf("unable to marshal config to YAML: %v", err)
			}

			fmt.Println("---")
			fmt.Println(string(bs))
		},
	}
}
