package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"beryju.io/gravity/pkg/instance"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run Gravity server",
	Run: func(cmd *cobra.Command, args []string) {
		inst := instance.New()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			inst.Stop()
		}()
		inst.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
