package detectbyFofa

import (
	"encoding/json"
	"fmt"
	"github.com/P001water/P1finger/libs/p1httputils"
	"github.com/projectdiscovery/gologger"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// params is key=>value for query, auto encoded with uri escape
func (c *FofaClient) buildURL(apiURI string, params map[string]string) string {
	fullURL := fmt.Sprintf("%s/api/%s/%s?", c.Server, c.APIVersion, apiURI)

	query := url.Values{}
	query.Set("email", c.Email)
	query.Set("key", c.Key)

	for k, v := range params {
		query.Set(k, v)
	}
	return fullURL + query.Encode()
}

// just fetch fofa body, no need to unmarshal
func (c *FofaClient) fetchBody(apiURI string, params map[string]string) (bodyRaw []byte, err error) {

	fullURL := c.buildURL(apiURI, params)
	gologger.Debug().Msgf("[fetch fofa]: %s", fullURL)

	req, err := http.NewRequest("GET", fullURL, nil)

	// to let go/http lib auto decode http request
	if _, ok := req.Header["Accept-Encoding"]; !ok {
		req.Header.Del("Accept-Encoding")
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		if !c.accountDebug {
			// 替换账号明文信息
			if e, ok := err.(*url.Error); ok {
				newClient := c
				newClient.Email = "<email>"
				newClient.Key = "<key>"
				e.URL = newClient.buildURL(apiURI, params)
				err = e
			}
		}
		return
	}
	defer resp.Body.Close()
	reqDump, _ := httputil.DumpRequest(req, true)
	gologger.Debug().Msgf("[debug info] Requests Http body: \n%v", string(reqDump))

	// response info for debug
	respDump, _ := httputil.DumpResponse(resp, true)
	gologger.Debug().Msgf("[debug info] Response Http body: \n%v", string(respDump))

	bodyRaw, err = p1httputils.ReadBodyTimeout(resp.Body, time.Second*3)
	if err != nil {
		return
	}

	return
}

// Fetch http request and parse as json return to v
func (c *FofaClient) Fetch(apiURI string, params map[string]string, v interface{}) (err error) {
	content, err := c.fetchBody(apiURI, params)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(content, v); err != nil {
		return err
	}
	return
}
