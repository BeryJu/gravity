package cmd

import (
	"github.com/spf13/cobra"
)

var cliUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "User Management",
}

func init() {
	cliCmd.AddCommand(cliUsersCmd)
}
