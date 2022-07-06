package cmd

import (
	"math/rand"
	"os"
	"time"

	"github.com/beryju/dns-dhcp-etcd-thingy/internal/instance"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dns-dhcp-etcd-thingy",
	Short: "A brief description of your application",
	Run: func(cmd *cobra.Command, args []string) {
		rand.Seed(time.Now().Unix())
		instance.NewInstance()
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
	log.SetFormatter(&log.JSONFormatter{
		DisableHTMLEscape: true,
	})
}
