package vars

import (
	"P1finger/cmd/flag"
	"P1finger/modules/detectbyFofa"
	"P1finger/modules/ruleClient"
	"net/http"
)

var (
	DetectResultTdSafe ruleClient.DetectResultTdSafeType //探测结果需要的数据
)

var (
	Options          = &flag.Options{}
	CustomhttpClient *http.Client

	FofaCli *detectbyFofa.FofaClient
	AppConf P1fingerConf
)
