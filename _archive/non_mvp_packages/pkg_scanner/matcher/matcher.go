package matcher

import (
	"regexp"
	"sync"
)

type Matcher struct {
	patterns map[string]*regexp.Regexp
	mu       sync.RWMutex
}

func NewMatcher() *Matcher {
	return &Matcher{
		patterns: make(map[string]*regexp.Regexp),
	}
}

func (m *Matcher) AddPattern(name, pattern string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	
	m.patterns[name] = regex
	return nil
}

func (m *Matcher) RemovePattern(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.patterns, name)
}

func (m *Matcher) Match(input string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	for name, pattern := range m.patterns {
		if pattern.MatchString(input) {
			return name, true
		}
	}
	
	return "", false
}

func (m *Matcher) ListPatterns() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var patterns []string
	for name := range m.patterns {
		patterns = append(patterns, name)
	}
	
	return patterns
}
