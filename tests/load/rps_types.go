// Copyright 2024 AegisGate
// RPS Load Test Types - Shared types for load testing
// This file contains non-test types used by load testing components

package load

import (
	"time"
)

// RPSLevel defines different RPS test levels
type RPSLevel int

const (
	RPS10K RPSLevel = iota
	RPS25K
	RPS50K
)

// String returns string representation of RPS level
func (r RPSLevel) String() string {
	switch r {
	case RPS10K:
		return "10,000"
	case RPS25K:
		return "25,000"
	case RPS50K:
		return "50,000"
	default:
		return "Unknown"
	}
}

// RPSConfig defines the configuration for an RPS load test
type RPSConfig struct {
	Level             RPSLevel
	TargetRPS         int
	Duration          time.Duration
	WarmupDuration    time.Duration
	CooldownDuration  time.Duration
	Concurrency       int
	RequestTimeout    time.Duration
	TargetHost        string
	TargetPort        int
	UseHTTPS          bool
	KeepAlive         bool
	NumClients        int
}

// RPSResult contains the results of an RPS load test
type RPSResult struct {
	Config           RPSConfig
	StartTime        time.Time
	EndTime          time.Time
	TotalDuration    time.Duration

	// Throughput metrics
	ActualRPS       float64
	TargetRPS       float64
	RPSAchievement  float64 // Percentage of target achieved

	// Request metrics
	TotalRequests   int64
	SuccessCount    int64
	ErrorCount      int64
	TimeoutCount    int64
	ConnectionError int64

	// Latency metrics
	MinLatency      time.Duration
	MaxLatency      time.Duration
	AvgLatency      time.Duration
	P50Latency      time.Duration
	P90Latency      time.Duration
	P95Latency      time.Duration
	P99Latency      time.Duration
	P999Latency     time.Duration
	P9999Latency    time.Duration

	// Latency percentiles (histogram buckets)
	LatencyBuckets  map[string]int64

	// Error breakdown
	ErrorBreakdown  map[string]int64

	// Resource metrics
	AvgMemoryMB     float64
	PeakMemoryMB    float64
	AvgCPUPercent   float64
	NumGC           uint32

	// Scaling metrics (for auto-scaling analysis)
	Samples         []RPSTimeSeriesSample
}

// RPSTimeSeriesSample represents a point in time series data
type RPSTimeSeriesSample struct {
	Timestamp     time.Time
	RPS            float64
	LatencyP50     time.Duration
	LatencyP99     time.Duration
	ErrorRate      float64
	MemoryMB       float64
	CPUPercent     float64
	ActiveConns    int
}

// DefaultRPSConfigs returns default configurations for each RPS level
func DefaultRPSConfigs() map[RPSLevel]RPSConfig {
	return map[RPSLevel]RPSConfig{
		RPS10K: {
			Level:            RPS10K,
			TargetRPS:        10000,
			Duration:         60 * time.Second,
			WarmupDuration:   10 * time.Second,
			CooldownDuration: 5 * time.Second,
			Concurrency:      500,
			RequestTimeout:   10 * time.Second,
			TargetHost:       "localhost",
			TargetPort:       8443,
			UseHTTPS:         true,
			KeepAlive:        true,
			NumClients:       100,
		},
		RPS25K: {
			Level:            RPS25K,
			TargetRPS:        25000,
			Duration:         60 * time.Second,
			WarmupDuration:   15 * time.Second,
			CooldownDuration: 5 * time.Second,
			Concurrency:      1000,
			RequestTimeout:   10 * time.Second,
			TargetHost:       "localhost",
			TargetPort:       8443,
			UseHTTPS:         true,
			KeepAlive:        true,
			NumClients:       200,
		},
		RPS50K: {
			Level:            RPS50K,
			TargetRPS:        50000,
			Duration:         60 * time.Second,
			WarmupDuration:   20 * time.Second,
			CooldownDuration: 5 * time.Second,
			Concurrency:      2000,
			RequestTimeout:   10 * time.Second,
			TargetHost:       "localhost",
			TargetPort:       8443,
			UseHTTPS:         true,
			KeepAlive:        true,
			NumClients:       400,
		},
	}
}

// LatencyBuckets returns predefined latency bucket boundaries
func LatencyBuckets() []string {
	return []string{
		"<1ms",
		"1-5ms",
		"5-10ms",
		"10-50ms",
		"50-100ms",
		"100-500ms",
		">500ms",
	}
}
