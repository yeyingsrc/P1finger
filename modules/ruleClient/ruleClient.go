package ruleClient

import (
	"crypto/tls"
	"net/http"
	"time"
)

type RuleClient struct {
	FingerFilePath  string
	FingersTdSafe   FingerPrintsTdSafeType
	ProxyUrl        string
	ProxyClient     *http.Client
	OutputFormat    string
	DetectRstTdSafe DetectResultTdSafeType

	RstShoot   DetectResultTdSafeType
	RstMiss    DetectResultTdSafeType
	RstReqFail DetectResultTdSafeType
}

// NewHttpClient creates a new HTTP client with optional settings
func (r *RuleClient) NewProxyClient() {

	transCfg := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, // disable verify
		DisableKeepAlives: true,
	}

	httpClient := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transCfg,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	r.ProxyClient = httpClient
}

func NewRuleClient() (*RuleClient, error) {

	r := &RuleClient{
		FingerFilePath: "P1fingersYaml",
	}

	r.NewProxyClient()

	return r, nil
}
