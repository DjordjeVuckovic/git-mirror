package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	configFile string
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "git-mirror",
	Short: "A tool to mirror Git repositories",
	Long: `Git Mirror is a CLI tool that mirrors a target repository to a source repository.
It supports various authentication methods and can be configured via YAML files.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yaml", "config file path")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
