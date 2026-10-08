package server

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
	RunE: func(cmd *cobra.Command, args []string) error {
		inst := instance.New()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			inst.Stop()
		}()
		return inst.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
