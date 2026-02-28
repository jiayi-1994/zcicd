package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zcicd/zcicd-cli/cmd"
)

var apiURL string

var rootCmd = &cobra.Command{
	Use:   "zcicd",
	Short: "ZCI/CD Platform CLI",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "http://localhost:8080", "API server URL")
	rootCmd.AddCommand(cmd.NewBuildCmd(&apiURL))
	rootCmd.AddCommand(cmd.NewDeployCmd(&apiURL))
	rootCmd.AddCommand(cmd.NewRolloutCmd(&apiURL))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
