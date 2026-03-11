// AegisGate Quick Test Runner - Completes in under 3 minutes
package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("AegisGate Quick Test Suite")
	fmt.Println("========================================")
	fmt.Println()

	passed := 0
	failed := 0

	// Test 1: Binary version check (5 seconds)
	fmt.Print("[1/5] Binary version check... ")
	if err := runWithTimeout(5*time.Second, "./bin/aegisgate.exe", "--version"); err != nil {
		fmt.Printf("FAIL: %v\n", err)
		failed++
	} else {
		fmt.Println("PASS")
		passed++
	}

	// Test 2: Help command (5 seconds)
	fmt.Print("[2/5] Help command... ")
	if err := runWithTimeout(5*time.Second, "./bin/aegisgate.exe", "--help"); err != nil {
		fmt.Printf("FAIL: %v\n", err)
		failed++
	} else {
		fmt.Println("PASS")
		passed++
	}

	// Test 3: Quick unit tests (60 seconds max)
	fmt.Print("[3/5] Unit tests (pkg/core, pkg/config)... ")
	if err := runWithTimeout(60*time.Second, "go", "test", "./pkg/core/...", "./pkg/config/...", "-v", "-count=1", "-timeout=30s"); err != nil {
		fmt.Printf("FAIL: %v\n", err)
		failed++
	} else {
		fmt.Println("PASS")
		passed++
	}

	// Test 4: Server starts (10 seconds)
	fmt.Print("[4/5] Server startup... ")
	if err := runWithTimeout(10*time.Second, "./bin/aegisgate.exe", "-bind", "0.0.0.0:9095", "-tier", "community"); err != nil {
		fmt.Printf("FAIL: %v\n", err)
		failed++
	} else {
		fmt.Println("PASS")
		passed++
	}

	// Test 5: Build verification (30 seconds)
	fmt.Print("[5/5] Build verification... ")
	if err := runWithTimeout(30*time.Second, "go", "build", "-o", "./bin/aegisgate_test.exe", "./cmd/aegisgate"); err != nil {
		fmt.Printf("FAIL: %v\n", err)
		failed++
	} else {
		fmt.Println("PASS")
		passed++
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Printf("Results: %d passed, %d failed\n", passed, failed)
	fmt.Println("========================================")

	// Cleanup
	exec.Command("taskkill", "/F", "/IM", "aegisgate.exe").Run()
	os.Remove("./bin/aegisgate_test.exe")

	if failed > 0 {
		os.Exit(1)
	}
}

func runWithTimeout(timeout time.Duration, args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = "C:\\Users\\Administrator\\Desktop\\Testing\\AegisGate"
	
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		cmd.Process.Kill()
		return fmt.Errorf("timeout after %v", timeout)
	}
}