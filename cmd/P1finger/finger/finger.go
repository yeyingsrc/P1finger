package finger

import (
	"github.com/P001water/P1finger/cmd"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(fingersCmd)
}

var fingersCmd = &cobra.Command{
	Use:   "finger",
	Short: "Operations on the P1finger fingerprint database",
	Long:  "Operations on the P1finger fingerprint database",
}
