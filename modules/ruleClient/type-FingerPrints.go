package ruleClient

import (
	"sync"
)

type FingerprintsType struct {
	ID          string        `yaml:"id"`
	Name        string        `yaml:"name"`
	Author      string        `yaml:"author"`
	Tags        []string      `yaml:"tags"`
	Description string        `yaml:"description"`
	Matchers    []MatcherType `yaml:"matchers"`
}

type MatcherType struct {
	Location        string   `yaml:"location,omitempty"`
	Type            string   `yaml:"type,omitempty"`
	Words           []string `yaml:"words,omitempty"`
	FaviconHash     []string `yaml:"hash,omitempty"`
	BodyHash        string   `yaml:"bodyHash,omitempty"`
	Accuracy        string   `yaml:"accuracy"`
	Condition       string   `yaml:"condition,omitempty"`
	CaseInsensitive bool     `yaml:"case-insensitive,omitempty"`
}

type FingerPrintsTdSafeType struct {
	mu          sync.Mutex
	FingerSlice []FingerprintsType
}

func (s *FingerPrintsTdSafeType) AddElement(elem FingerprintsType) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.FingerSlice = append(s.FingerSlice, elem)
}

func (s *FingerPrintsTdSafeType) GetElements() []FingerprintsType {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.FingerSlice
}
