package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/extconfig"
	"github.com/spf13/cobra"
)

var apiUrl string
var apiToken string

var apiClient *api.APIClient

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Interact with a running Gravity server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		url, err := url.Parse(apiUrl)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		config := api.NewConfiguration()
		config.Debug = true
		config.Host = url.Host
		config.Scheme = url.Scheme
		if url.Scheme == "unix" {
			config.Scheme = "http"
			config.Host = "socket"
			config.HTTPClient = &http.Client{
				Transport: &http.Transport{
					DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
						return net.Dial("unix", url.Path)
					},
				},
			}
		}
		if apiToken != "" {
			config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", apiToken))
		}
		apiClient = api.NewAPIClient(config)
	},
}

func init() {
	defUrl := "unix:///var/run/gravity.sock"
	if extconfig.Get().Debug {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		defUrl = fmt.Sprintf("unix://%s/gravity.sock", cwd)
	}
	cliCmd.PersistentFlags().StringVarP(&apiUrl, "host", "s", defUrl, "API Host")
	cliCmd.PersistentFlags().StringVarP(&apiToken, "token", "t", "", "API Token")
	rootCmd.AddCommand(cliCmd)
}
