package logging

import (
	"fmt"
	"log"
	"time"
)

// Logger provides structured logging for MCP tools
type Logger struct {
	enabled bool
}

// NewLogger creates a new logger instance
func NewLogger(enabled bool) *Logger {
	return &Logger{enabled: enabled}
}

// Info logs informational messages
func (l *Logger) Info(tool string, message string, args ...interface{}) {
	if !l.enabled {
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(message, args...)
	log.Printf("[INFO] %s [%s] %s", timestamp, tool, msg)
}

// Error logs error messages
func (l *Logger) Error(tool string, message string, err error, args ...interface{}) {
	if !l.enabled {
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(message, args...)
	log.Printf("[ERROR] %s [%s] %s: %v", timestamp, tool, msg, err)
}

// Debug logs debug messages
func (l *Logger) Debug(tool string, message string, args ...interface{}) {
	if !l.enabled {
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(message, args...)
	log.Printf("[DEBUG] %s [%s] %s", timestamp, tool, msg)
}

// Warning logs warning messages
func (l *Logger) Warning(tool string, message string, args ...interface{}) {
	if !l.enabled {
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(message, args...)
	log.Printf("[WARN] %s [%s] %s", timestamp, tool, msg)
}
