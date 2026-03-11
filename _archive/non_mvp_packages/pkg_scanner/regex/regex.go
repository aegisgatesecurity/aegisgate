package regex

import (
	"regexp"
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
	
	s.patterns["ssn"] = regexp.MustCompile("[0-9]{3}-[0-9]{2}-[0-9]{4}")
	s.patterns["phone"] = regexp.MustCompile("[0-9]{3}-[0-9]{3}-[0-9]{4}")
}

func (s *Scanner) AddPattern(name, pattern string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	
	s.patterns[name] = regex
	return nil
}

func (s *Scanner) RemovePattern(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.patterns, name)
}

func (s *Scanner) FindMatches(pattern, text string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	regex, exists := s.patterns[pattern]
	if !exists {
		return nil
	}
	
	return regex.FindAllString(text, -1)
}

func (s *Scanner) ListPatterns() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var patterns []string
	for name := range s.patterns {
		patterns = append(patterns, name)
	}
	
	return patterns
}
