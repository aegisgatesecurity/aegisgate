package metrics

import (
	"testing"
	"time"
)

func TestNewMetrics(t *testing.T) {
	m := NewMetrics()
	
	if m == nil {
		t.Error("NewMetrics should return non-nil metrics")
	}
	
	if !m.IsEnabled() {
		t.Error("Metrics should be enabled by default")
	}
	
	if m.bufferSize != 1000 {
		t.Errorf("Expected buffer size 1000, got %d", m.bufferSize)
	}
}

func TestIncrementCounter(t *testing.T) {
	m := NewMetrics()
	
	m.IncrementCounter("test", 1)
	m.IncrementCounter("test", 1)
	
	val := m.GetCounter("test")
	if val != 2 {
		t.Errorf("Expected counter value 2, got %d", val)
	}
}

func TestDecrementCounter(t *testing.T) {
	m := NewMetrics()
	
	m.IncrementCounter("test", 10)
	m.DecrementCounter("test", 3)
	
	val := m.GetCounter("test")
	if val != 7 {
		t.Errorf("Expected counter value 7, got %d", val)
	}
}

func TestRecordDuration(t *testing.T) {
	m := NewMetrics()
	
	m.RecordDuration("test", 100*time.Millisecond)
	m.RecordDuration("test", 200*time.Millisecond)
	
	count, _, _, _ := m.GetDurationStats("test")
	if count != 2 {
		t.Errorf("Expected 2 durations, got %d", count)
	}
}

func TestRecordBytes(t *testing.T) {
	m := NewMetrics()
	
	m.RecordBytes("test", 100)
	m.RecordBytes("test", 200)
	
	count, total, _, _, _ := m.GetByteStats("test")
	
	if count != 2 {
		t.Errorf("Expected 2 byte records, got %d", count)
	}
	if total != 300 {
		t.Errorf("Expected total 300, got %d", total)
	}
}

func TestGetDurationStats(t *testing.T) {
	m := NewMetrics()
	
	m.RecordDuration("test", 100*time.Millisecond)
	m.RecordDuration("test", 200*time.Millisecond)
	m.RecordDuration("test", 300*time.Millisecond)
	
	count, avg, min, max := m.GetDurationStats("test")
	
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
	
	if avg != 200*time.Millisecond {
		t.Errorf("Expected avg 200ms, got %v", avg)
	}
	
	if min != 100*time.Millisecond {
		t.Errorf("Expected min 100ms, got %v", min)
	}
	
	if max != 300*time.Millisecond {
		t.Errorf("Expected max 300ms, got %v", max)
	}
}

func TestEnableDisable(t *testing.T) {
	m := NewMetrics()
	
	m.Disable()
	if m.IsEnabled() {
		t.Error("Metrics should be disabled after Disable()")
	}
	
	m.Enable()
	if !m.IsEnabled() {
		t.Error("Metrics should be enabled after Enable()")
	}
}

func TestReset(t *testing.T) {
	m := NewMetrics()
	
	m.IncrementCounter("test", 1)
	m.RecordDuration("test", 100*time.Millisecond)
	m.RecordBytes("test", 100)
	
	m.Reset()
	
	if m.GetCounter("test") != 0 {
		t.Error("Counter should be reset to 0")
	}
}

func TestGetCounters(t *testing.T) {
	m := NewMetrics()
	
	m.IncrementCounter("test1", 1)
	m.IncrementCounter("test2", 2)
	
	counters := m.GetCounters()
	
	if len(counters) != 2 {
		t.Errorf("Expected 2 counters, got %d", len(counters))
	}
}

func TestGetSummary(t *testing.T) {
	m := NewMetrics()
	
	m.IncrementCounter("test", 1)
	
	summary := m.GetSummary()
	
	if !summary["enabled"].(bool) {
		t.Error("Summary should show enabled")
	}
	
	if summary["counters"].(int) != 1 {
		t.Error("Summary should show 1 counter")
	}
}

func TestBufferLimit(t *testing.T) {
	m := NewMetrics()
	
	// Change buffer size for testing
	m.bufferSize = 3
	
	for i := 0; i < 5; i++ {
		m.RecordDuration("test", time.Duration(i)*time.Millisecond)
	}
	
	count, _, _, _ := m.GetDurationStats("test")
	
	if count != 3 {
		t.Errorf("Expected buffer limit of 3, got %d", count)
	}
}