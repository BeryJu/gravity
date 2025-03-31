package externaldns

import (
	"fmt"
	"net/url"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/externaldns"
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

var ExternalDNSCommand = &cobra.Command{
	Use:   "external-dns",
	Short: "Allow for external-dns to talk to gravity",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		url, err := url.Parse(apiUrl)
		if err != nil {
			logger.Error("failed to parse API URL", zap.Error(err))
			return
		}

		config := api.NewConfiguration()
		config.Host = url.Host
		config.Scheme = url.Scheme
		if apiToken != "" {
			config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", apiToken))
		}
		config.UserAgent = fmt.Sprintf("gravity-cli/%s", extconfig.FullVersion())
		apiClient = api.NewAPIClient(config)
	},
	Run: func(cmd *cobra.Command, args []string) {
		s := externaldns.New(apiClient)
		s.Run()
	},
}

func init() {
	logger = extconfig.Get().Logger().Named("external-dns")
	ExternalDNSCommand.PersistentFlags().StringVarP(&apiUrl, "host", "s", "", "API Host")
	ExternalDNSCommand.PersistentFlags().StringVarP(&apiToken, "token", "t", "", "API Token")
}
