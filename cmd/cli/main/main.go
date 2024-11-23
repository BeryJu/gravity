package main

import (
	"os"

	"beryju.io/gravity/cmd/cli"
)

func main() {
	cli.CLICmd.Use = "gravity-cli"
	err := cli.CLICmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
