package cli

import (
	"fmt"
	"net/http"
	"syscall"

	"beryju.io/gravity/api"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var cliUsersAddCmd = &cobra.Command{
	Use:   "add username",
	Short: "Add a user for the API",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		fmt.Printf("Enter the password for %s: ", username)
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}
		fmt.Println("")
		hr, err := apiClient.RolesApiAPI.ApiPutUsers(cmd.Context()).Username(username).AuthAPIUsersPutInput(api.AuthAPIUsersPutInput{
			Password: string(bytePassword),
			Permissions: []api.AuthPermission{
				{
					Path:    api.PtrString("/*"),
					Methods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodHead, http.MethodDelete},
				},
			},
		}).Execute()
		checkApiError(hr, err)
	},
}

func init() {
	cliUsersCmd.AddCommand(cliUsersAddCmd)
}
