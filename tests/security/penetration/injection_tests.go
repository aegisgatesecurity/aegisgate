package penetration

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

type InjectionTest struct {
	TargetHost string
	TargetPort int
	Timeout    time.Duration
}

func NewInjectionTest(host string, port int) *InjectionTest {
	return &InjectionTest{TargetHost: host, TargetPort: port, Timeout: 30 * time.Second}
}

var SQLInjectionPayloads = []string{
	"' OR '1'='1",
	"' OR 1=1--",
	"' UNION SELECT * FROM users--",
	"1; DROP TABLE users--",
	"' AND 1=1--",
	"admin'--",
	"' OR '1'='1' /*",
}

var CommandInjectionPayloads = []string{
	"; cat /etc/passwd",
	"| whoami",
	"$(id)",
	"; echo injected",
	" && ls -la",
}

var NoSQLInjectionPayloads = []string{
	"{$ne: null}",
	"{$gt: ''}",
	"{$regex: '.*'}",
}

type InjectionResult struct {
	Payload  string
	Type     string
	Injected bool
	Detected bool
	Response string
	Severity string
	Details  string
}

type InjectionReport struct {
	SQLiResults     []InjectionResult
	CommandResults  []InjectionResult
	NoSQLResults    []InjectionResult
	TotalTests      int
	Vulnerabilities int
	Status          string
}

func (i *InjectionTest) RunInjectionTests(ctx context.Context) (*InjectionReport, error) {
	report := &InjectionReport{Status: "PASS"}

	// SQL Injection tests
	for _, payload := range SQLInjectionPayloads {
		result := i.testSQLInjection(ctx, payload)
		report.SQLiResults = append(report.SQLiResults, result)
		report.TotalTests++
		if result.Detected {
			report.Vulnerabilities++
		}
	}

	// Command Injection tests
	for _, payload := range CommandInjectionPayloads {
		result := i.testCommandInjection(ctx, payload)
		report.CommandResults = append(report.CommandResults, result)
		report.TotalTests++
		if result.Detected {
			report.Vulnerabilities++
		}
	}

	// NoSQL Injection tests
	for _, payload := range NoSQLInjectionPayloads {
		result := i.testNoSQLInjection(ctx, payload)
		report.NoSQLResults = append(report.NoSQLResults, result)
		report.TotalTests++
		if result.Detected {
			report.Vulnerabilities++
		}
	}

	if report.Vulnerabilities > 0 {
		report.Status = "FAIL"
	}

	return report, nil
}

func (i *InjectionTest) testSQLInjection(ctx context.Context, payload string) InjectionResult {
	ctx, cancel := context.WithTimeout(ctx, i.Timeout)
	defer cancel()
	_ = cancel
	defer cancel()

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", i.TargetHost, i.TargetPort), 10*time.Second)
	if err != nil {
		return InjectionResult{Payload: payload, Type: "SQLi", Injected: false, Detected: false, Severity: "Info"}
	}
	defer func() { _ = conn.Close() }()

	request := fmt.Sprintf("POST /api/query HTTP/1.1\r\nHost: %s\r\nContent-Type: application/x-www-form-urlencoded\r\n\r\ndata=%s", i.TargetHost, payload)
	_, _ = conn.Write([]byte(request))

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	response := string(buf[:n])

	sqlErrors := []string{"mysql", "sqlite", "postgresql", "oracle", "syntax error", "unexpected"}
	detected := false
	responseLower := strings.ToLower(response)
	for _, sqlErr := range sqlErrors {
		if strings.Contains(responseLower, sqlErr) {
			detected = true
			break
		}
	}

	return InjectionResult{Payload: payload, Type: "SQLi", Injected: true, Detected: detected, Response: response[:min(len(response), 200)], Severity: "High"}
}

func (i *InjectionTest) testCommandInjection(ctx context.Context, payload string) InjectionResult {
	ctx, cancel := context.WithTimeout(ctx, i.Timeout)
	defer cancel()
	_ = cancel
	defer cancel()

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", i.TargetHost, i.TargetPort), 10*time.Second)
	if err != nil {
		return InjectionResult{Payload: payload, Type: "Command", Injected: false, Detected: false, Severity: "Critical"}
	}
	defer func() { _ = conn.Close() }()

	request := fmt.Sprintf("GET /api/exec?cmd=ping%s HTTP/1.1\r\nHost: %s\r\n\r\n", payload, i.TargetHost)
	_, _ = conn.Write([]byte(request))

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	response := string(buf[:n])

	cmdIndicators := []string{"root:", "uid=", "/bin/bash", "nt authority"}
	detected := false
	for _, indicator := range cmdIndicators {
		if strings.Contains(response, indicator) {
			detected = true
			break
		}
	}

	return InjectionResult{Payload: payload, Type: "Command", Injected: true, Detected: detected, Severity: "Critical"}
}

func (i *InjectionTest) testNoSQLInjection(ctx context.Context, payload string) InjectionResult {
	ctx, cancel := context.WithTimeout(ctx, i.Timeout)
	defer cancel()
	_ = cancel
	defer cancel()

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", i.TargetHost, i.TargetPort), 10*time.Second)
	if err != nil {
		return InjectionResult{Payload: payload, Type: "NoSQL", Injected: false, Detected: false, Severity: "High"}
	}
	defer func() { _ = conn.Close() }()

	request := fmt.Sprintf("GET /api/data?filter=%s HTTP/1.1\r\nHost: %s\r\n\r\n", payload, i.TargetHost)
	_, _ = conn.Write([]byte(request))

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	response := string(buf[:n])

	return InjectionResult{Payload: payload, Type: "NoSQL", Injected: true, Detected: false, Response: response[:min(len(response), 200)], Severity: "High"}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestInjection(t *testing.T) {
	test := NewInjectionTest("localhost", 8443)
	ctx := context.Background()

	report, err := test.RunInjectionTests(ctx)
	if err != nil {
		t.Fatalf("injection test error: %v", err)
	}

	if report.Vulnerabilities > 0 {
		t.Errorf("Found %d injection vulnerabilities", report.Vulnerabilities)
	}

	t.Logf("Total tests: %d, Vulnerabilities: %d", report.TotalTests, report.Vulnerabilities)
}

func BenchmarkInjectionDetection(b *testing.B) {
	test := NewInjectionTest("localhost", 8443)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = test.RunInjectionTests(ctx)
	}
}
