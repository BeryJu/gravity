package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var cliTokensAddCmd = &cobra.Command{
	Use:   "add username",
	Short: "Add a token for the API",
	Args: func(cmd *cobra.Command, args []string) error {
		// Optionally run one of the validators provided by cobra
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		token, _, err := apiClient.RolesApiApi.ApiPutTokens(cmd.Context()).Username(username).Execute()
		if err != nil {
			logger.Error("failed to add token", zap.Error(err))
			return
		}
		fmt.Println(token)
	},
}

func init() {
	cliTokensCmd.AddCommand(cliTokensAddCmd)
}
