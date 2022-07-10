package cmd

import (
	"math/rand"
	"os"
	"time"

	"beryju.io/gravity/pkg/instance"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gravity",
	Short: "A brief description of your application",
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

func init() {
	log.SetLevel(log.TraceLevel)
	// log.SetFormatter(&log.JSONFormatter{
	// 	DisableHTMLEscape: true,
	// })
}
