package logging

import (
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	if logger == nil {
		t.Error("NewLogger() returned nil")
	}
}

func TestSetLevel(t *testing.T) {
	logger := NewLogger()
	logger.SetLevel(LoggerDebug)
	if logger.GetLevel() != LoggerDebug {
		t.Errorf("Expected debug level, got %v", logger.GetLevel())
	}
}

func TestDebug(t *testing.T) {
	logger := NewLogger()
	logger.SetLevel(LoggerDebug)
	logger.Debug("test debug message")
}

func TestInfo(t *testing.T) {
	logger := NewLogger()
	logger.Info("test info message")
}

func TestWarn(t *testing.T) {
	logger := NewLogger()
	logger.Warn("test warn message")
}

func TestError(t *testing.T) {
	logger := NewLogger()
	logger.Error("test error message")
}

func TestFormatTime(t *testing.T) {
	timeStr := FormatTime(time.Now())
	if timeStr == "" {
		t.Error("FormatTime() returned empty string")
	}
}

func TestGetLevelString(t *testing.T) {
	if GetLevelString(LoggerDebug) != "DEBUG" {
		t.Error("GetLevelString(LoggerDebug) != DEBUG")
	}
	if GetLevelString(LoggerInfo) != "INFO" {
		t.Error("GetLevelString(LoggerInfo) != INFO")
	}
	if GetLevelString(LoggerWarn) != "WARN" {
		t.Error("GetLevelString(LoggerWarn) != WARN")
	}
	if GetLevelString(LoggerError) != "ERROR" {
		t.Error("GetLevelString(LoggerError) != ERROR")
	}
}
