package cmd

import (
	"fmt"
	"syscall"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var username string

var addUserCmd = &cobra.Command{
	Use:   "addUser",
	Short: "Add a user for the API",
	Run: func(cmd *cobra.Command, args []string) {
		if username == "" {
			panic("Must set -u")
		}
		inst := instance.NewInstance()
		inst.ForRole("api").AddEventListener(instance.EventTopicInstanceBootstrapped, func(ev *roles.Event) {
			fmt.Printf("Enter the password for %s: ", username)
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				panic(err)
			}
			fmt.Println("")
			api := api.New(inst.ForRole("api"))
			err = api.CreateUser(username, string(bytePassword))
			if err != nil {
				panic(err)
			}
			inst.Stop()
		})
		inst.Start()
	},
}

func init() {
	rootCmd.AddCommand(addUserCmd)
	addUserCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "set Username")
}
