package cmd

import (
	_ "unsafe"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/spf13/cobra"
	_ "go.etcd.io/etcd/etcdctl/v3/ctlv3"
	"go.etcd.io/etcd/etcdctl/v3/ctlv3/command"
)

//go:linkname etcdctlCommand go.etcd.io/etcd/etcdctl/v3/ctlv3.rootCmd
var etcdctlCommand *cobra.Command

//go:linkname globalFlags go.etcd.io/etcd/etcdctl/v3/ctlv3.globalFlags
var globalFlags command.GlobalFlags

func init() {
	globalFlags.Endpoints = []string{extconfig.Get().Etcd.Endpoint}
	cliCmd.AddCommand(etcdctlCommand)
}
