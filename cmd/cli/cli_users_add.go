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
	Args: func(cmd *cobra.Command, args []string) error {
		// Optionally run one of the validators provided by cobra
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		fmt.Printf("Enter the password for %s: ", username)
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}
		fmt.Println("")
		hr, err := apiClient.RolesApiApi.ApiPutUsers(cmd.Context()).Username(username).AuthAPIUsersPutInput(api.AuthAPIUsersPutInput{
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
