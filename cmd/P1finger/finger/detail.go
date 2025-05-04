package finger

import (
	"encoding/json"
	"fmt"
	"github.com/P001water/P1finger/cmd/vars"
	"github.com/P001water/P1finger/modules/ruleClient"
	"github.com/projectdiscovery/gologger"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

func init() {
	fingersCmd.AddCommand(detailCmd)
}

var detailCmd = &cobra.Command{
	Use:   "detail",
	Short: "查看一个指纹的详细信息",
	Run: func(cmd *cobra.Command, args []string) {

		fingerName := args[0]
		p1ruleClient, err := ruleClient.NewRuleClientBuilder().
			WithCustomizeFingerFile(vars.AppConf.CustomizeFingerFile).
			WithDefaultFingerFiles(vars.AppConf.UseDefaultFingerFiles).
			WithOutputFormat(vars.Options.Output).
			WithTimeout(10 * time.Second).
			Build()
		if err != nil {
			gologger.Error().Msgf("%v", err)
		}

		fingers := p1ruleClient.P1FingerPrints.GetElements()

		for _, finger := range fingers {
			if strings.Contains(finger.Name, fingerName) {
				jsonData, err := json.Marshal(finger)
				if err != nil {
					fmt.Printf("Error converting to JSON: %v\n", err)
					continue
				}
				fmt.Println(string(jsonData))
			}
		}
	},
}
