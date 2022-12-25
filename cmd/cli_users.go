package cmd

import "github.com/spf13/cobra"

var cliUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Commands related to user management",
}

func init() {
	cliCmd.AddCommand(cliUsersCmd)
}
