package main

import "beryju.io/gravity/cmd/cli"

func main() {
	cli.CLICmd.Use = "gravity-cli"
	cli.CLICmd.Execute()
}
