package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

func NewBuildCmd(apiURL *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build operations",
	}

	trigger := &cobra.Command{
		Use:   "trigger [config-id]",
		Short: "Trigger a build",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("%s/api/v1/builds/%s/run", *apiURL, args[0])
			resp, err := http.Post(url, "application/json", nil)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			return printJSON(resp.Body)
		},
	}

	status := &cobra.Command{
		Use:   "status [run-id]",
		Short: "Get build run status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("%s/api/v1/builds/runs/%s", *apiURL, args[0])
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			return printJSON(resp.Body)
		},
	}

	cmd.AddCommand(trigger, status)
	return cmd
}

func printJSON(r io.Reader) error {
	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	var out map[string]interface{}
	if err := json.Unmarshal(body, &out); err != nil {
		fmt.Println(string(body))
		return nil
	}
	pretty, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(pretty))
	return nil
}
