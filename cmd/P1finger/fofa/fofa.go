package fofa

import (
	"github.com/P001water/P1finger/cmd"
	"github.com/P001water/P1finger/cmd/vars"
	"github.com/P001water/P1finger/libs/fileutils"
	"github.com/P001water/P1finger/modules/FofaClient"
	"github.com/P001water/P1finger/modules/RuleClient"
	"github.com/projectdiscovery/gologger"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

var (
	Url     string
	UrlFile string
)

func init() {
	cmd.RootCmd.AddCommand(fofaCmd)
	fofaCmd.Flags().StringVarP(&vars.Options.Url, "url", "u", "", "target url")
	fofaCmd.Flags().StringVarP(&vars.Options.UrlFile, "file", "f", "", "target url file")
}

var fofaCmd = &cobra.Command{
	Use:   "fofa",
	Short: "基于Fofa空间测绘引擎的指纹识别",
	Long:  "基于Fofa的指纹识别",
	Run: func(cmd *cobra.Command, args []string) {
		gologger.Info().Msgf("p1fingeprint detect model: Fofa\n")

		err := fofaAction()
		if err != nil {
			gologger.Error().Msg(err.Error())
			return
		}
	},
}

func fofaAction() (err error) {

	// 整合目标输入
	var targetUrls []string
	if vars.Options.Url != "" {
		targetUrls = append(targetUrls, vars.Options.Url)
	}

	if vars.Options.UrlFile != "" {
		var urlsFromFile []string
		filePath := filepath.Join(vars.ExecDir, vars.Options.UrlFile)
		urlsFromFile, err = fileutils.ReadLinesFromFile(filePath)
		if err != nil {
			gologger.Error().Msgf("%v", err)
			return
		}
		targetUrls = append(targetUrls, urlsFromFile...)
	}

	if len(targetUrls) <= 0 {
		gologger.Error().Msg("input url is null")
		return
	}

	vars.FofaCli, err = FofaClient.NewClient(
		FofaClient.WithURL("https://fofa.info/?email=&key=&version=v1"),
		FofaClient.WithAccountDebug(true),
		FofaClient.WithDebug(vars.Options.Debug),
		FofaClient.WithEmail(vars.AppConf.FofaCredentials.Email),
		FofaClient.WithApiKey(vars.AppConf.FofaCredentials.ApiKey),
	)

	// 美化查询语法用于Fofa
	var group []string
	var FinalQuery []string
	var querybeautify []string
	domains, ips := FofaClient.SplitDomainsAndIPs(targetUrls)
	querybeautify = append(querybeautify, domains...)
	querybeautify = append(querybeautify, ips...)
	for i, simpleQuery := range querybeautify {
		group = append(group, simpleQuery)
		if (i+1)%50 == 0 || i == len(querybeautify)-1 {
			FinalQuery = append(FinalQuery, strings.Join(group, " || "))
			group = nil
		}
	}

	// 开始查询
	queryFields := []string{"ip", "port", "title", "product", "lastupdatetime", "protocol", "host"}
	for _, item := range FinalQuery {
		res, err := vars.FofaCli.HostSearch(item, -1, queryFields)
		if err != nil {
			gologger.Error().Msgf("%v", err)
			return err
		}

		for _, simpleTarget := range res {
			tmp := RuleClient.DetectResult{
				OriginUrl:      simpleTarget[5] + "://" + simpleTarget[0] + ":" + simpleTarget[1],
				Host:           simpleTarget[6],
				WebTitle:       simpleTarget[2],
				FingerTag:      strings.Split(simpleTarget[3], ","),
				LastUpdateTime: simpleTarget[4],
			}
			vars.DetectResultTdSafe.AddElement(tmp)
		}
	}

	err = RuleClient.SaveToFile(vars.DetectResultTdSafe.GetElements(), vars.Options.Output)
	if err != nil {
		gologger.Error().Msg(err.Error())
		return
	}

	return nil
}
