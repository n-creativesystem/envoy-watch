package cmd

import (
	"github.com/n-creativesystem/envoy-watch/version"
	"github.com/spf13/cobra"
)

func NewCmdVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			ver := version.Version
			rev := version.Revision
			cmd.Printf("ver. %s, rev. %s\n", ver, rev)
		},
	}
	return cmd
}
