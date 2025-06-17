package cmd

import (
	"fmt"
	"github.com/DjordjeVuckovic/git-mirror/internal/config"
	"log"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration file",
	Long:  `Validate the configuration file syntax and required fields.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			log.Fatalf("Configuration validation failed: %v", err)
		}

		fmt.Printf("Configuration file '%s' is valid\n", configFile)
		fmt.Printf("Found %d mirror configurations\n", len(cfg.Mirrors))

		for i, mirror := range cfg.Mirrors {
			fmt.Printf("  Mirror %d:\n", i+1)
			fmt.Printf("    Target: %s (%s)\n", mirror.Target.URL, mirror.Target.Auth.Method)
			fmt.Printf("    Source: %s (%s)\n", mirror.Source.URL, mirror.Source.Auth.Method)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
