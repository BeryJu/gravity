package externaldns

import (
	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/externaldns"
	"github.com/spf13/cobra"
)

var ExternalDNSCommand = &cobra.Command{
	Use:     "external-dns",
	Short:   "Allow for external-dns to talk to gravity",
	Version: extconfig.FullVersion(),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := externaldns.New()
		if err != nil {
			return err
		}
		return s.Run()
	},
}
