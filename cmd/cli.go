package cmd

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/extconfig"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	apiUrl   string
	apiToken string
)

var (
	apiClient *api.APIClient
	logger    *zap.Logger
)

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Interact with a running Gravity server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		url, err := url.Parse(apiUrl)
		if err != nil {
			logger.Error("failed to parse API URL", zap.Error(err))
			return
		}

		config := api.NewConfiguration()
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
	logger = extconfig.Get().Logger().Named("cli")
	cliCmd.PersistentFlags().StringVarP(&apiUrl, "host", "s", defUrl, "API Host")
	cliCmd.PersistentFlags().StringVarP(&apiToken, "token", "t", "", "API Token")
	rootCmd.AddCommand(cliCmd)
}

func checkApiError(hr *http.Response, err error) {
	if err == nil {
		return
	}
	b, _ := io.ReadAll(hr.Body)
	logger.Error("failed to send request", zap.String("response", string(b)), zap.Error(err))
}
