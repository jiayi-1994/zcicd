package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

func NewRolloutCmd(apiURL *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollout",
		Short: "Rollout operations (canary/blue-green)",
	}

	status := &cobra.Command{
		Use:   "status [config-id]",
		Short: "Get rollout status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("%s/api/v1/deploys/%s/rollout", *apiURL, args[0])
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			return printJSON(resp.Body)
		},
	}

	promote := &cobra.Command{
		Use:   "promote [config-id]",
		Short: "Promote rollout to next step",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("%s/api/v1/deploys/%s/rollout/promote", *apiURL, args[0])
			resp, err := http.Post(url, "application/json", nil)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			return printJSON(resp.Body)
		},
	}

	abort := &cobra.Command{
		Use:   "abort [config-id]",
		Short: "Abort current rollout",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("%s/api/v1/deploys/%s/rollout/abort", *apiURL, args[0])
			resp, err := http.Post(url, "application/json", nil)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			return printJSON(resp.Body)
		},
	}

	cmd.AddCommand(status, promote, abort)
	return cmd
}
