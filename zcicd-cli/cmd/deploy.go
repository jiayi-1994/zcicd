package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

func NewDeployCmd(apiURL *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy operations",
	}

	sync := &cobra.Command{
		Use:   "sync [config-id]",
		Short: "Trigger deployment sync",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("%s/api/v1/deploys/%s/sync", *apiURL, args[0])
			resp, err := http.Post(url, "application/json", nil)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			return printJSON(resp.Body)
		},
	}

	status := &cobra.Command{
		Use:   "status [config-id]",
		Short: "Get deployment status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("%s/api/v1/deploys/%s/status", *apiURL, args[0])
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			return printJSON(resp.Body)
		},
	}

	cmd.AddCommand(sync, status)
	return cmd
}
