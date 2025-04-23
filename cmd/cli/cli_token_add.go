package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cliTokensAddCmd = &cobra.Command{
	Use:   "add username",
	Short: "Add a token for the API",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		token, hr, err := apiClient.RolesApiApi.ApiPutTokens(cmd.Context()).Username(username).Execute()
		if err != nil {
			checkApiError(hr, err)
			return
		}
		fmt.Println(token)
	},
}

func init() {
	cliTokensCmd.AddCommand(cliTokensAddCmd)
}
