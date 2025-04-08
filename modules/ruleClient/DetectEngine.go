package ruleClient

import (
	"P1finger/libs/p1httputils"
	"regexp"
	"strings"
)

const (
	Up   = "up"
	Down = "down"
)

func (r *RuleClient) Detect(url string) (err error) {

	var DetectRst DetectResult
	var P1fingerResps []p1httputils.P1fingerHttpResp

	fixedUrl, err := CheckHttpPrefix(url)
	if err != nil {
		DetectRst = DetectResult{
			OriginUrl: url,
			WebTitle:  "无法访问",
			SiteUp:    Down,
			FingerTag: []string{"unknown proto http/https, try manually"},
		}
		r.RstReqFail.AddElement(DetectRst)
		return
	}

	resp, p1fingerResp, err := p1httputils.HttpGet(fixedUrl, r.ProxyClient)
	if err != nil {
		DetectRst = DetectResult{
			OriginUrl:           p1fingerResp.Url,
			OriginUrlStatusCode: p1fingerResp.StatusCode,
			WebTitle:            p1fingerResp.WebTitle,
			SiteUp:              Down,
			FingerTag:           []string{"WebSite down"},
		}
		r.DetectRstTdSafe.AddElement(DetectRst)
		r.RstReqFail.AddElement(DetectRst)
		return
	}
	P1fingerResps = append(P1fingerResps, p1fingerResp)

	var P1fingerRedirectResp p1httputils.P1fingerHttpResp
	redirectType, _, isRedirect := p1httputils.CheckPageRedirect(resp, p1fingerResp)
	if isRedirect {
		switch redirectType {
		case "Location":
			redirectClient := p1httputils.NewRedirectHttpClient()
			_, P1fingerRedirectResp, err = p1httputils.HttpGet(url, redirectClient)
			if err != nil {
				return err
			}
			P1fingerResps = append(P1fingerResps, P1fingerRedirectResp)

		case "jsRedirect":
			// todo
		case "VueRoute":
			// todo
		}
	}

	for _, finger := range r.FingersTdSafe.GetElements() {
		for _, matcher := range finger.Matchers {
			matchFlag := false
			for _, targetInfo := range P1fingerResps {
				var content string

				switch matcher.Location {
				case "title":
					content = targetInfo.WebTitle
				case "header":
					if matcher.Type == "regex" {
						re := regexp.MustCompile(strings.Join(matcher.Words, "|")) // 预编译正则
						if re.MatchString(targetInfo.HeaderStr) {
							matchFlag = true
							finger.Name = re.FindString(targetInfo.HeaderStr)
						}
					} else {
						content = targetInfo.HeaderStr
					}
				case "body":
					content = targetInfo.BodyStr
				default:
					continue
				}

				if matcher.Type != "regex" && matchCondition(content, matcher.Words, matcher.Condition) {
					matchFlag = true
				}

				if matchFlag {
					DetectRst.FingerTag = append(DetectRst.FingerTag, finger.Name)
					break
				}
			}

		}
	}

	if len(DetectRst.FingerTag) <= 0 {
		DetectRst = DetectResult{
			OriginUrl:           p1fingerResp.Url,
			OriginUrlStatusCode: p1fingerResp.StatusCode,
			WebTitle:            p1fingerResp.WebTitle,
			SiteUp:              Up,
			FingerTag:           []string{"Target up but missed."},
		}
		r.RstMiss.AddElement(DetectRst)
	} else {
		DetectRst = DetectResult{
			OriginUrl:           p1fingerResp.Url,
			OriginUrlStatusCode: p1fingerResp.StatusCode,
			WebTitle:            p1fingerResp.WebTitle,
			SiteUp:              Up,
			FingerTag:           DetectRst.FingerTag,
		}

		if isRedirect {
			DetectRst.WebTitle = P1fingerRedirectResp.WebTitle
		}

		r.RstShoot.AddElement(DetectRst)
	}

	r.DetectRstTdSafe.AddElement(DetectRst)
	return nil
}

func matchCondition(content string, words []string, condition string) bool {
	shooted := 0
	loweredContent := strings.ToLower(content)

	for _, word := range words {
		loweredWord := strings.ToLower(word)
		if strings.Contains(loweredContent, loweredWord) {
			shooted++
			if condition == "or" || condition == "" {
				return true
			}
		} else if condition == "and" {
			return false
		}
	}

	return condition == "and" && shooted == len(words)
}
