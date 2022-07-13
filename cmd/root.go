package cmd

import (
	"math/rand"
	"os"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gravity",
	Short:   "Start gravity instance",
	Version: extconfig.FullVersion(),
	Run: func(cmd *cobra.Command, args []string) {
		rand.Seed(time.Now().Unix())
		inst := instance.NewInstance()
		defer inst.Stop()
		inst.Start()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
