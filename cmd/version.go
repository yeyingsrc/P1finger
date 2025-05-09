package cmd

import (
	"github.com/P001water/P1finger/cmd/vars"
	"github.com/projectdiscovery/gologger"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the P1finger Version",
	Long:  "Print the P1finger Version",
	Run: func(cmd *cobra.Command, args []string) {
		gologger.Info().Msgf("Current Version %s", vars.P1fingerVer)
	},
}
