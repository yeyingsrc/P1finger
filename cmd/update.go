package cmd

import (
	"bufio"
	"fmt"
	"github.com/P001water/P1-github-selfupdate/selfupdate"
	"github.com/projectdiscovery/gologger"
	"github.com/spf13/cobra"
	"os"
)

const Version = "0.1.1"

func init() {
	RootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "self update",
	Run: func(cmd *cobra.Command, args []string) {
		err := CheckUpdate()
		if err != nil {
			gologger.Error().Msgf("%v", err)
			return
		}
	},
}

func CheckUpdate() error {
	isLatest, latest, err := selfupdate.CheckVersionIsLatest(Version, "P001water/P1finger")
	if err != nil {
		return err
	}
	if isLatest {
		gologger.Info().Msgf("Current version is the latest")
	} else {
		gologger.Info().Msgf("Current version: %v", Version)
		gologger.Info().Msgf("latest version: %v, Whether to update ? (yes or no)", latest.Version)
		fmt.Printf("Input your choose: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		userInput := scanner.Text()

		switch userInput {
		case "yes":
			gologger.Info().Msg("Updating...")
			exe, err := os.Executable()
			if err != nil {
				gologger.Info().Msg("Could not locate executable path")
				return err
			}
			if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
				gologger.Info().Msgf("Error occurred while updating binary: %v\n", err)
				return err
			}
			gologger.Info().Msgf("Successfully updated to version: %v\n", latest.Version)
		case "no":
			gologger.Info().Msg("Update cancelled.")
			return nil
		default:
			gologger.Info().Msg("Invalid input. Update cancelled.")
			return nil
		}
	}

	return nil
}
