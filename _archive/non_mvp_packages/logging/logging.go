package logging

import (
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// LogLevel defines the logging levels
type LogLevel int

const (
	LoggerDebug LogLevel = iota
	LoggerInfo
	LoggerWarn
	LoggerError
)

// Logger represents the application logger
type Logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	level       LogLevel
	mu          sync.RWMutex
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		level:       LoggerInfo,
	}
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLevel returns the current logging level
func (l *Logger) GetLevel() LogLevel {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

// SetOutput sets the output writer for all loggers
func (l *Logger) SetOutput(w io.Writer) {
	l.infoLogger.SetOutput(w)
	l.warnLogger.SetOutput(w)
	l.errorLogger.SetOutput(w)
	l.debugLogger.SetOutput(w)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	l.mu.RLock()
	level := l.level
	l.mu.RUnlock()
	
	if level <= LoggerDebug {
		l.debugLogger.Println(msg)
	}
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.mu.RLock()
	level := l.level
	l.mu.RUnlock()
	
	if level <= LoggerInfo {
		l.infoLogger.Println(msg)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	l.mu.RLock()
	level := l.level
	l.mu.RUnlock()
	
	if level <= LoggerWarn {
		l.warnLogger.Println(msg)
	}
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.errorLogger.Println(msg)
}

// WithFields adds structured fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	// TODO: Implement structured logging with fields
	return l
}

// FormatTime formats a time timestamp
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000")
}

// GetLevelString converts LogLevel to string
func GetLevelString(level LogLevel) string {
	switch level {
	case LoggerDebug:
		return "DEBUG"
	case LoggerInfo:
		return "INFO"
	case LoggerWarn:
		return "WARN"
	case LoggerError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// NewContextLogger creates a logger with context
func NewContextLogger(ctx map[string]interface{}) *Logger {
	l := NewLogger()
	// TODO: Add context fields
	return l
}

// RotateConfig holds rotation configuration
type RotateConfig struct {
	MaxSize    int    // megabytes
	MaxAge     int    // days
	MaxBackups int    // number of backups
	Compress   bool   // compress backups
}

// FileLogger writes logs to a file
type FileLogger struct {
	logger     *Logger
	file       *os.File
	config     RotateConfig
	path       string
	mu         sync.Mutex
}

// NewFileLogger creates a file-based logger
func NewFileLogger(path string, config RotateConfig) (*FileLogger, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	
	return &FileLogger{
		logger: NewLogger(),
		file:   file,
		config: config,
		path:   path,
	}, nil
}

// Close closes the file logger
func (f *FileLogger) Close() error {
	return f.file.Close()
}

// Rotate rotates the log file
func (f *FileLogger) Rotate() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	// TODO: Implement log rotation logic
	return nil
}

// Debug logs a debug message to file
func (f *FileLogger) Debug(msg string) {
	f.logger.Debug(msg)
}

// Info logs an info message to file
func (f *FileLogger) Info(msg string) {
	f.logger.Info(msg)
}

// Warn logs a warning message to file
func (f *FileLogger) Warn(msg string) {
	f.logger.Warn(msg)
}

// Error logs an error message to file
func (f *FileLogger) Error(msg string) {
	f.logger.Error(msg)
}
