package vars

import (
	"github.com/P001water/P1finger/cmd/flag"
	"github.com/P001water/P1finger/modules/detectbyFofa"
	"github.com/P001water/P1finger/modules/ruleClient"
	"os"
	"path/filepath"
)

var (
	DetectResultTdSafe ruleClient.DetectResultTdSafeType //探测结果需要的数据
)

var (
	Options = &flag.Options{}

	ExecDir = filepath.Dir(os.Args[0])
	FofaCli *detectbyFofa.FofaClient
	AppConf P1fingerConf
)
