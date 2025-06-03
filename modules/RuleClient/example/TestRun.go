package example

import (
	"fmt"
	"github.com/P001water/P1finger/modules/RuleClient"
	"github.com/projectdiscovery/gologger"
)

func TestRun(url string) {
	// 创建一个新的 RuleClient 实例
	client, _ := RuleClient.NewRuleClient()
	// 加载指纹数据
	err := client.LoadFingersFromExEfs()
	if err != nil {
		gologger.Error().Msgf("%v", err)
		return
	}

	err = client.Detect(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = RuleClient.SaveToFile(client.DetectRstTdSafe.GetElements(), client.OutputFormat)
	if err != nil {
		gologger.Error().Msg(err.Error())
		return
	}

}
