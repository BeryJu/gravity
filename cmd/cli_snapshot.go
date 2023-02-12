package cmd

import (
	"github.com/spf13/cobra"
	"go.etcd.io/etcd/etcdutl/v3/etcdutl"
)

func init() {
	sc := etcdutl.NewSnapshotCommand()
	sc.PersistentFlags().StringVarP(&etcdutl.OutputFormat, "write-out", "w", "table", "set the output format (fields, json, protobuf, simple, table)")
	sc.RegisterFlagCompletionFunc("write-out", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"fields", "json", "protobuf", "simple", "table"}, cobra.ShellCompDirectiveDefault
	})
	cliCmd.AddCommand(sc)
}
