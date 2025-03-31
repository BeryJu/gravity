package main

import (
	"os"

	externaldns "beryju.io/gravity/cmd/external-dns"
)

func main() {
	err := externaldns.ExternalDNSCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}
