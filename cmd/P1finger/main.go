package main

import (
	"P1finger/cmd/vars"
	"P1finger/libs/fileutils"
	"P1finger/libs/goflags"
	"P1finger/modules/detectbyFofa"
	"P1finger/modules/p1fmt"
	"P1finger/modules/ruleClient"
	"fmt"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/projectdiscovery/gologger"
	"github.com/schollz/progressbar/v3"
	"os"
	"strings"
	"sync"
)

const (
	P1fingerConfFile = "p1fingerConf.yaml"
)

func ReadCmdflag() error {
	flagSet := goflags.NewFlagSet()

	flagSet.CreateGroup("scanConf", "scanConf",
		flagSet.StringVar(&vars.Options.Url, "u", "", "target url, eg: [-u www.example.com]"),
		flagSet.StringVar(&vars.Options.UrlFile, "uf", "", "urls in file, eg: [-uf demo.txt]"),
		//flagSet.StringVar(&vars.Options.P1fingerFile, "fingerfile", "p1finger.json", "fingers file, eg: [-fingerfile p1finger.json]"),
		flagSet.IntVar(&vars.Options.Rate, "rate", 500, "Concurrency quantity, eg: [-rate 500] default: 500"),
		flagSet.BoolVar(&vars.Options.Update, "update", false, "update tools, eg: [-update]"),
	)

	flagSet.CreateGroup("collectMode", "collectMode",
		flagSet.StringVar(&vars.Options.CollectMode, "m", "rule", "Fingerprint recognition mode: [-m rule] / [-o fofa]"),
	)

	flagSet.CreateGroup("Proxy", "Proxy",
		flagSet.StringVar(&vars.Options.SocksProxy, "socks", "", "socks proxy"),
		flagSet.StringVar(&vars.Options.HttpProxy, "httpproxy", "", "http proxy"),
	)

	flagSet.CreateGroup("output", "output format",
		flagSet.StringVar(&vars.Options.Output, "o", "p1finger.json", "output file name: [-o p1finger.xlsx] / [-o p1finger.json]"),
	)

	flagSet.CreateGroup("debug", "debug",
		flagSet.BoolVar(&vars.Options.Debug, "debug", false, "http debug info, eg:[-debug]"),
	)

	if err := flagSet.Parse(); err != nil {
		return err
	}

	return nil
}

func main() {

	vars.Banner()
	err := ReadCmdflag()
	if err != nil {
		gologger.Error().Msgf("Could not parse flags: %s\n", err)
		return
	}

	//检查更新
	if vars.Options.Update {
		err = vars.CheckUpdate()
		if err != nil {
			gologger.Error().Msgf("%v", err)
			return
		}
		return
	} else {
		gologger.Info().Msgf("默认不进行新版本检查，请用 -update 参数手动检查工具更新")
	}

	// load P1finger config from p1fingerConf.yaml file
	err = vars.LoadAppConf(P1fingerConfFile, &vars.AppConf)
	if err != nil {
		gologger.Error().Msgf("%v", err)
		return
	}

	err = vars.LoadHttpClient()
	if err != nil {
		gologger.Error().Msgf("%v", err)
		return
	}

	if vars.Options.Url == "" && vars.Options.UrlFile == "" {
		gologger.Info().Msg("Pls input url")
		os.Exit(1)
	}

	switch vars.Options.CollectMode {
	case "rule":
		gologger.Info().Msgf("p1fingeprint detect model: [%v]\n", vars.Options.CollectMode)
		err = ruleAction()
		if err != nil {
			gologger.Error().Msg(err.Error())
			return
		}
	case "fofa":
		gologger.Info().Msgf("p1fingeprint detect model: [%v]\n", vars.Options.CollectMode)
		err = fofaAction()
		if err != nil {
			gologger.Error().Msg(err.Error())
			return
		}
	default:
		gologger.Info().Msgf("p1fingeprint detect model: [rule]\n")
		err := ruleAction()
		if err != nil {
			gologger.Error().Msg(err.Error())
			return
		}
	}
}

func ruleAction() (err error) {

	client, _ := ruleClient.NewRuleClient()

	// 加载指纹数据
	if len(vars.AppConf.CustomizeFingerFile) > 0 {
		err := client.LoadFingersFromFile(vars.AppConf.CustomizeFingerFile)
		if err != nil {
			return err
		}
	}

	if vars.AppConf.UseDefaultFingerFils || len(vars.AppConf.CustomizeFingerFile) == 0 {
		err = client.LoadFingersFromExEfs()
		if err != nil {
			return
		}
	}

	client.OutputFormat = vars.Options.Output

	var targetUrls []string
	if vars.Options.Url != "" {
		targetUrls = append(targetUrls, vars.Options.Url)
	}

	if vars.Options.UrlFile != "" {
		var urlsFromFile []string
		urlsFromFile, err = fileutils.ReadLinesFromFile(vars.Options.UrlFile)
		if err != nil {
			gologger.Error().Msgf("%v", err)
			return
		}
		targetUrls = append(targetUrls, urlsFromFile...)
	}

	bar := progressbar.NewOptions(len(targetUrls),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()), //you should install "github.com/k0kubun/go-ansi"
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription("[cyan][P1finger][reset] 检测进度..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	var workWg sync.WaitGroup
	concurrency := vars.Options.Rate
	urlChan := make(chan string, len(targetUrls))

	// 启动固定数量的 worker
	for i := 0; i < concurrency; i++ {
		go func() {
			for url := range urlChan {
				defer workWg.Done()
				_ = client.Detect(url)
				bar.Add(1)
			}
		}()
	}

	// 发送任务到 channel
	for _, url := range targetUrls {
		workWg.Add(1)
		urlChan <- url
	}

	close(urlChan)
	workWg.Wait()
	bar.Finish()
	fmt.Println()

	for _, shoot := range client.RstShoot.GetElements() {
		prefix := color.New(color.FgGreen).Add(color.Bold).Sprintf("[已命中]")
		// P1finger的指纹还在清洗阶段，区分生产和测试模式，清洗指纹
		if vars.Options.Debug {
			p1fmt.PrintfShoot(prefix, shoot.OriginUrl, shoot.WebTitle, shoot.FingerTag, shoot.OriginUrlStatusCode)
		} else {
			p1fmt.PrintfShoot(prefix, shoot.OriginUrl, shoot.WebTitle, ruleClient.SliceRmDuplication(shoot.FingerTag), shoot.OriginUrlStatusCode)
		}

	}

	for _, miss := range client.RstMiss.GetElements() {
		prefix := color.New(color.FgRed).Add(color.Bold).Sprintf("[未命中]")
		if vars.Options.Debug {
			p1fmt.PrintMiss(prefix, miss.OriginUrl, miss.WebTitle, miss.FingerTag, miss.OriginUrlStatusCode)
		} else {
			p1fmt.PrintMiss(prefix, miss.OriginUrl, miss.WebTitle, ruleClient.SliceRmDuplication(miss.FingerTag), miss.OriginUrlStatusCode)
		}
	}

	for _, reqFail := range client.RstReqFail.GetElements() {
		prefix := color.New(color.FgBlue).Add(color.Bold).Sprintf("[无法访问]")

		if vars.Options.Debug {
			p1fmt.PrintReqFail(prefix, reqFail.OriginUrl, reqFail.FingerTag)
		} else {
			p1fmt.PrintReqFail(prefix, reqFail.OriginUrl, ruleClient.SliceRmDuplication(reqFail.FingerTag))
		}
	}

	err = ruleClient.SaveToFile(client.DetectRstTdSafe.GetElements(), client.OutputFormat)
	if err != nil {
		gologger.Error().Msg(err.Error())
		return
	}

	return
}

func fofaAction() (err error) {

	vars.FofaCli, err = detectbyFofa.NewClient(
		detectbyFofa.WithURL("https://fofa.info/?email=&key=&version=v1"),
		detectbyFofa.WithAccountDebug(true),
		detectbyFofa.WithDebug(vars.Options.Debug),
		detectbyFofa.WithEmail(vars.AppConf.FofaCredentials.Email),
		detectbyFofa.WithApiKey(vars.AppConf.FofaCredentials.ApiKey),
	)

	// 汇总查询的url
	var urls []string
	if vars.Options.Url != "" {
		urls = append(urls, vars.Options.Url)
	}

	if vars.Options.UrlFile != "" {
		var tmp []string
		tmp, err = fileutils.ReadLinesFromFile(vars.Options.UrlFile)
		if err != nil {
			return
		}
		urls = append(urls, tmp...)
	}

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
