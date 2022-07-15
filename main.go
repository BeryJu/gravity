package main

import (
	"math/rand"
	"time"

	"beryju.io/gravity/cmd"
)

func main() {
	rand.Seed(time.Now().Unix())
	cmd.Execute()
}
