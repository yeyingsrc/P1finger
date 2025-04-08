package ruleClient

import "sync"

type DetectResult struct {
	Host                  string   `json:"host"`
	OriginUrl             string   `json:"origin_target"`
	RedirectUrl           string   `json:"redirect_url"`
	OriginUrlStatusCode   int      `json:"origin_url_status_code"`
	RedirectUrlStatusCode int      `json:"redirect_url_status_code"`
	WebTitle              string   `json:"web_title"` //Important
	SiteUp                string   `json:"site_up"`
	FingerTag             []string `json:"finger_tag"`     //Important
	LastUpdateTime        string   `json:"lastupdatetime"` //Important
}

type DetectResultTdSafeType struct {
	mu           sync.Mutex
	TargetFinger []DetectResult
}

func (s *DetectResultTdSafeType) AddElement(elem DetectResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TargetFinger = append(s.TargetFinger, elem)
}

func (s *DetectResultTdSafeType) GetElements() []DetectResult {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.TargetFinger
}
