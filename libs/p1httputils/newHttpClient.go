package p1httputils

import (
	"context"
	"crypto/tls"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"time"
)

// HttpClientOption defines the type for setting optional parameters
type HttpClientOption func(*http.Transport)

// NewNoRedirectHttpClient creates a new HTTP client with optional settings
func NewNoRedirectHttpClient(options ...HttpClientOption) *http.Client {
	transCfg := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, // disable verify
		DisableKeepAlives: true,
	}

	for _, opt := range options {
		opt(transCfg)
	}

	httpClient := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transCfg,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 禁止跟随重定向，手动处理重定向逻辑
		},
	}

	return httpClient
}

// NewNoRedirectHttpClient creates a new HTTP client with optional settings
func NewRedirectHttpClient(options ...HttpClientOption) *http.Client {
	transCfg := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, // disable verify
		DisableKeepAlives: true,
	}

	for _, opt := range options {
		opt(transCfg)
	}

	httpClient := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transCfg,
	}

	return httpClient
}

// WithSocks5Proxy sets the SOCKS5 proxy
func WithSocks5Proxy(socks5Proxy string) HttpClientOption {
	return func(transCfg *http.Transport) {
		if socks5Proxy != "" {
			dialer, err := proxy.SOCKS5("tcp", socks5Proxy, nil, proxy.Direct)
			if err == nil {
				transCfg.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
					return dialer.Dial(network, addr)
				}
			}
		}
	}
}

// WithHttpProxy sets the HTTP proxy
func WithHttpProxy(httpProxy string) HttpClientOption {
	return func(transCfg *http.Transport) {
		if httpProxy != "" {
			transCfg.Proxy = func(_ *http.Request) (*url.URL, error) {
				return url.Parse(httpProxy)
			}
		}
	}
}
