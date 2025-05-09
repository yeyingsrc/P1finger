package ruleClient

import (
	"crypto/tls"
	"github.com/P001water/P1finger/libs/p1httputils"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type RuleClient struct {
	DefaultFingerPath     string                 // 默认指纹库路径
	P1FingerPrints        FingerPrintsTdSafeType // P1finger指纹库
	CustomizeFingerFiles  []string               // 可选 自定义指纹文件
	UseDefaultFingerFiles bool                   // 可选 自定义指纹文件后是否启用默认指纹库

	ProxyUrl              string // 可选 代理地址
	ProxyClient           *http.Client
	ProxyNoRedirectCilent *http.Client

	OutputFormat string // 可选 输出格式

	DetectRstTdSafe DetectResultTdSafeType
	RstShoot        DetectResultTdSafeType
	RstMiss         DetectResultTdSafeType
	RstReqFail      DetectResultTdSafeType
}

type RuleClientBuilder struct {
	defaultFingerPath     string
	useDefaultFingerFiles bool
	customizeFingerFiles  []string
	outputFormat          string
	timeout               time.Duration
	proxyURL              string
}

func NewRuleClientBuilder() *RuleClientBuilder {
	return &RuleClientBuilder{
		customizeFingerFiles:  []string{},
		useDefaultFingerFiles: true,   // 默认值
		outputFormat:          "json", // 默认值
		timeout:               5 * time.Second,
	}
}

func (b *RuleClientBuilder) Build() (*RuleClient, error) {
	r := &RuleClient{
		DefaultFingerPath:     "P1fingersYaml",
		UseDefaultFingerFiles: b.useDefaultFingerFiles,
		CustomizeFingerFiles:  b.customizeFingerFiles,
		OutputFormat:          b.outputFormat,
		ProxyUrl:              b.proxyURL,
	}

	var err error
	// 加载自定义指纹（如果有）
	if len(r.CustomizeFingerFiles) > 0 {
		err = r.LoadFingersFromFile(filepath.Dir(os.Args[0]), r.CustomizeFingerFiles)
		if err != nil {
			return nil, err
		}
	}

	// 加载默认指纹（如果需要）
	if r.UseDefaultFingerFiles || len(r.CustomizeFingerFiles) == 0 {
		err = r.LoadFingersFromExEfs()
		if err != nil {
			return nil, err
		}
	}

	r.newProxyClientWithTimeout(b.timeout)

	r.ProxyClient = p1httputils.NewHttpClientBuilder().
		WithProxy(b.proxyURL).
		Build()

	r.ProxyNoRedirectCilent = p1httputils.NewHttpClientBuilder().
		WithProxy(b.proxyURL).
		NoRedirect().
		Build()

	return r, nil
}

func (r *RuleClient) newProxyClientWithTimeout(timeout time.Duration) {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	transCfg := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
	}
	r.ProxyClient = &http.Client{
		Timeout:   timeout,
		Transport: transCfg,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func (b *RuleClientBuilder) WithDefaultFingerFiles(val bool) *RuleClientBuilder {
	b.useDefaultFingerFiles = val
	return b
}

func (b *RuleClientBuilder) WithCustomizeFingerFile(files []string) *RuleClientBuilder {
	b.customizeFingerFiles = files
	return b
}

func (b *RuleClientBuilder) WithOutputFormat(format string) *RuleClientBuilder {
	b.outputFormat = format
	return b
}

func (b *RuleClientBuilder) WithTimeout(t time.Duration) *RuleClientBuilder {
	b.timeout = t
	return b
}

func (b *RuleClientBuilder) WithProxyURL(url string) *RuleClientBuilder {
	b.proxyURL = url
	return b
}
