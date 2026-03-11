package metrics

import (
	"sync"
	"time"
)

// Metrics collects proxy metrics
type Metrics struct {
	counters    map[string]int64
	durations   map[string][]time.Duration
	bytes       map[string][]int64
	mu          sync.RWMutex
	enabled     bool
	bufferSize  int
}

// NewMetrics creates a new metrics collector
func NewMetrics() *Metrics {
	return &Metrics{
		counters:   make(map[string]int64),
		durations:  make(map[string][]time.Duration),
		bytes:      make(map[string][]int64),
		enabled:    true,
		bufferSize: 1000,
	}
}

// IncrementCounter increments a counter
func (m *Metrics) IncrementCounter(name string, value int64) {
	if !m.enabled {
		return
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.counters[name] += value
}

// DecrementCounter decrements a counter
func (m *Metrics) DecrementCounter(name string, value int64) {
	if !m.enabled {
		return
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.counters[name] -= value
}

// RecordDuration records a duration
func (m *Metrics) RecordDuration(name string, d time.Duration) {
	if !m.enabled {
		return
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if len(m.durations[name]) >= m.bufferSize {
		m.durations[name] = m.durations[name][1:]
	}
	m.durations[name] = append(m.durations[name], d)
}

// RecordBytes records byte count
func (m *Metrics) RecordBytes(name string, bytes int64) {
	if !m.enabled {
		return
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if len(m.bytes[name]) >= m.bufferSize {
		m.bytes[name] = m.bytes[name][1:]
	}
	m.bytes[name] = append(m.bytes[name], bytes)
}

// GetCounter returns counter value
func (m *Metrics) GetCounter(name string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if val, exists := m.counters[name]; exists {
		return val
	}
	return 0
}

// GetDurationStats returns duration statistics
func (m *Metrics) GetDurationStats(name string) (count int, avg, min, max time.Duration) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	durations := m.durations[name]
	if len(durations) == 0 {
		return 0, 0, 0, 0
	}
	
	count = len(durations)
	for _, d := range durations {
		avg += d
		if min == 0 || d < min {
			min = d
		}
		if d > max {
			max = d
		}
	}
	avg /= time.Duration(count)
	
	return count, avg, min, max
}

// GetByteStats returns byte statistics
func (m *Metrics) GetByteStats(name string) (count int, total, avg, min, max int64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	bytes := m.bytes[name]
	if len(bytes) == 0 {
		return 0, 0, 0, 0, 0
	}
	
	count = len(bytes)
	for _, b := range bytes {
		total += b
		if min == 0 || b < min {
			min = b
		}
		if b > max {
			max = b
		}
	}
	avg = total / int64(count)
	
	return count, total, avg, min, max
}

// Enable enables metrics collection
func (m *Metrics) Enable() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.enabled = true
}

// Disable disables metrics collection
func (m *Metrics) Disable() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.enabled = false
}

// IsEnabled returns whether metrics are enabled
func (m *Metrics) IsEnabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.enabled
}

// Reset resets all metrics
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.counters = make(map[string]int64)
	m.durations = make(map[string][]time.Duration)
	m.bytes = make(map[string][]int64)
}

// GetCounters returns all counters
func (m *Metrics) GetCounters() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	counters := make(map[string]int64, len(m.counters))
	for k, v := range m.counters {
		counters[k] = v
	}
	return counters
}

// GetDurations returns all durations
func (m *Metrics) GetDurations() map[string][]time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	durations := make(map[string][]time.Duration)
	for k, v := range m.durations {
		durations[k] = v
	}
	return durations
}

// GetBytes returns all byte counts
func (m *Metrics) GetBytes() map[string][]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	bytes := make(map[string][]int64)
	for k, v := range m.bytes {
		bytes[k] = v
	}
	return bytes
}

// GetSummary returns a summary of all metrics
func (m *Metrics) GetSummary() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return map[string]interface{}{
		"enabled":   m.enabled,
		"counters":  len(m.counters),
		"durations": len(m.durations),
		"bytes":     len(m.bytes),
	}
}