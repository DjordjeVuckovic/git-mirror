package cmd

import (
	"fmt"
	"github.com/DjordjeVuckovic/git-mirror/internal/config"
	"github.com/DjordjeVuckovic/git-mirror/internal/mirror"
	"log"

	"github.com/spf13/cobra"
)

var mirrorCmd = &cobra.Command{
	Use:   "mirror",
	Short: "Mirror repositories according to configuration",
	Long:  `Mirror repositories from target to source according to the configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		if verbose {
			fmt.Printf("Loaded configuration with %d mirror jobs\n", len(cfg.Mirrors))
		}

		for i, mirrorConfig := range cfg.Mirrors {
			if verbose {
				fmt.Printf("Processing mirror job %d: %s -> %s\n", i+1, mirrorConfig.Target.URL, mirrorConfig.Source.URL)
			}

			m := mirror.New(mirrorConfig, verbose)
			if err := m.Execute(); err != nil {
				log.Printf("Failed to mirror %s to %s: %v", mirrorConfig.Target.URL, mirrorConfig.Source.URL, err)
				continue
			}

			if verbose {
				fmt.Printf("Successfully mirrored %s to %s\n", mirrorConfig.Target.URL, mirrorConfig.Source.URL)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(mirrorCmd)
}
