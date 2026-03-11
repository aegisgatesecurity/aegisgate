package llm

import (
	"regexp"
	"strings"
	"sync"
)

type Scanner struct {
	patterns map[string]*regexp.Regexp
	mu       sync.RWMutex
}

func NewScanner() *Scanner {
	s := &Scanner{
		patterns: make(map[string]*regexp.Regexp),
	}
	s.initPatterns()
	return s
}

func (s *Scanner) initPatterns() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.patterns["prompt_injection"] = regexp.MustCompile(`ignore\s+(previous|all|prior|above)\s+(instruction|prompt|command)`)
	s.patterns["jailbreak"] = regexp.MustCompile(`disregard|bypass|override\s+(your|system|constraints)`)
}

func (s *Scanner) Scan(text string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var threats []string
	
	for name, pattern := range s.patterns {
		if pattern.MatchString(strings.ToLower(text)) {
			threats = append(threats, name)
		}
	}
	
	return threats
}

func (s *Scanner) HasThreat(text string) bool {
	return len(s.Scan(text)) > 0
}
