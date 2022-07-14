package cmd

import (
	"fmt"
	"syscall"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/api/auth"
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
		rootInst := instance.NewInstance()
		inst := rootInst.ForRole("tests")
		inst.AddEventListener(instance.EventTopicInstanceBootstrapped, func(ev *roles.Event) {
			api := api.New(inst)
			am := auth.NewAuthProvider(api, inst)
			fmt.Printf("Enter the password for %s: ", username)
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				panic(err)
			}
			fmt.Println("")
			err = am.CreateUser(username, string(bytePassword))
			if err != nil {
				panic(err)
			}
			rootInst.Stop()
		})
		rootInst.Start()
	},
}

func init() {
	rootCmd.AddCommand(addUserCmd)
	addUserCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "set Username")
}
