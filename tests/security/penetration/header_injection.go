package penetration

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

type HeaderInjectionTest struct {
	TargetHost string
	TargetPort int
	Timeout    time.Duration
}

func NewHeaderInjectionTest(host string, port int) *HeaderInjectionTest {
	return &HeaderInjectionTest{TargetHost: host, TargetPort: port, Timeout: 30 * time.Second}
}

var MaliciousHeaders = []struct {
	Name  string
	Value string
}{
	{"X-Forwarded-For", "127.0.0.1"},
	{"Host", "attacker.com"},
	{"Referer", "javascript:alert(1)"},
	{"User-Agent", "<script>alert(1)</script>"},
	{"X-HTTP-Method-Override", "DELETE"},
}

var MissingSecurityHeaders = []string{
	"Strict-Transport-Security",
	"X-Content-Type-Options",
	"X-Frame-Options",
}

type HeaderResult struct {
	HeaderName     string
	HeaderValue    string
	Injected       bool
	Accepted       bool
	SecurityHeader string
	Present        bool
	Recommendation string
}

type HeaderReport struct {
	InjectionTests  []HeaderResult
	MissingHeaders  []HeaderResult
	TotalTests      int
	Vulnerabilities int
	Status          string
}

func (h *HeaderInjectionTest) RunHeaderInjectionTests(ctx context.Context) (*HeaderReport, error) {
	report := &HeaderReport{InjectionTests: make([]HeaderResult, 0), Status: "PASS"}

	for _, header := range MaliciousHeaders {
		result := h.testHeaderInjection(ctx, header.Name, header.Value)
		report.InjectionTests = append(report.InjectionTests, result)
		report.TotalTests++
		if result.Accepted {
			report.Vulnerabilities++
		}
	}

	if report.Vulnerabilities > 0 {
		report.Status = "FAIL"
	}
	return report, nil
}

func (h *HeaderInjectionTest) testHeaderInjection(ctx context.Context, name, value string) HeaderResult {
	ctx, cancel := context.WithTimeout(ctx, h.Timeout)
	defer cancel()
	_ = cancel // ignore cancel error in this context
	defer cancel()

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", h.TargetHost, h.TargetPort), 10*time.Second)
	if err != nil {
		return HeaderResult{HeaderName: name, HeaderValue: value, Injected: false, Accepted: false}
	}
	defer func() { _ = conn.Close() }()

	request := fmt.Sprintf("GET / HTTP/1.1\r\nHost: %s\r\n%s: %s\r\n\r\n", h.TargetHost, name, value)
	_, err = conn.Write([]byte(request))
	if err != nil {
		return HeaderResult{HeaderName: name, HeaderValue: value, Injected: false, Accepted: false}
	}

	buf := make([]byte, 4096)
	n, _ := conn.Read(buf)
	response := string(buf[:n])
	accepted := !strings.Contains(response, "400") && !strings.Contains(response, "error")

	return HeaderResult{HeaderName: name, HeaderValue: value, Injected: true, Accepted: accepted}
}

func TestHeaderInjection(t *testing.T) {
	test := NewHeaderInjectionTest("localhost", 8443)
	ctx := context.Background()
	report, err := test.RunHeaderInjectionTests(ctx)
	if err != nil {
		t.Fatalf("header injection test error: %v", err)
	}
	if report.Vulnerabilities > 0 {
		t.Errorf("Header injection vulnerabilities: %d", report.Vulnerabilities)
	}
}
