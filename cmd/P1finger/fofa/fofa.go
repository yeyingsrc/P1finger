package fofa

import (
	"github.com/P001water/P1finger/cmd"
	"github.com/P001water/P1finger/cmd/vars"
	"github.com/P001water/P1finger/libs/fileutils"
	"github.com/P001water/P1finger/modules/detectbyFofa"
	"github.com/P001water/P1finger/modules/ruleClient"
	"github.com/projectdiscovery/gologger"
	"github.com/spf13/cobra"
	"strings"
)

var (
	Url     string
	UrlFile string
)

func init() {
	cmd.RootCmd.AddCommand(fofaCmd)
	fofaCmd.Flags().StringVarP(&Url, "url", "u", "", "target url")
	fofaCmd.Flags().StringVarP(&UrlFile, "file", "f", "", "target url file")
}

var fofaCmd = &cobra.Command{
	Use:   "fofa",
	Short: "基于Fofa空间测绘引擎的指纹识别",
	Long:  "基于Fofa的指纹识别",
	Run: func(cmd *cobra.Command, args []string) {
		gologger.Info().Msgf("p1fingeprint detect model: Fofa\n")

		var err error
		// 汇总查询的url
		var targetUrls []string
		if Url != "" {
			targetUrls = append(targetUrls, Url)
		}

		if UrlFile != "" {
			var tmp []string
			tmp, err = fileutils.ReadLinesFromFile(UrlFile)
			if err != nil {
				return
			}
			targetUrls = append(targetUrls, tmp...)
		}

		if len(targetUrls) <= 0 {
			gologger.Error().Msg("input url is null")
			return
		}

		err = fofaAction(targetUrls)
		if err != nil {
			gologger.Error().Msg(err.Error())
			return
		}
	},
}

func fofaAction(urls []string) (err error) {

	vars.FofaCli, err = detectbyFofa.NewClient(
		detectbyFofa.WithURL("https://fofa.info/?email=&key=&version=v1"),
		detectbyFofa.WithAccountDebug(true),
		detectbyFofa.WithDebug(vars.Options.Debug),
		detectbyFofa.WithEmail(vars.AppConf.FofaCredentials.Email),
		detectbyFofa.WithApiKey(vars.AppConf.FofaCredentials.ApiKey),
	)

	// 美化查询语法用于Fofa
	var group []string
	var FinalQuery []string
	var querybeautify []string
	domains, ips := detectbyFofa.SplitDomainsAndIPs(urls)
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
			tmp := ruleClient.DetectResult{
				OriginUrl:      simpleTarget[5] + "://" + simpleTarget[0] + ":" + simpleTarget[1],
				Host:           simpleTarget[6],
				WebTitle:       simpleTarget[2],
				FingerTag:      strings.Split(simpleTarget[3], ","),
				LastUpdateTime: simpleTarget[4],
			}
			vars.DetectResultTdSafe.AddElement(tmp)
		}
	}

	err = ruleClient.SaveToFile(vars.DetectResultTdSafe.GetElements(), vars.Options.Output)
	if err != nil {
		gologger.Error().Msg(err.Error())
		return
	}

	return nil
}
