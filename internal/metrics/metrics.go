package metrics

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Metrics tracks performance metrics for MCP tools
type Metrics struct {
	mu            sync.RWMutex
	toolCalls     map[string]int
	toolDurations map[string][]time.Duration
	toolErrors    map[string]int
	enabled       bool
}

// NewMetrics creates a new metrics instance
func NewMetrics(enabled bool) *Metrics {
	return &Metrics{
		toolCalls:     make(map[string]int),
		toolDurations: make(map[string][]time.Duration),
		toolErrors:    make(map[string]int),
		enabled:       enabled,
	}
}

// RecordToolCall records a tool invocation
func (m *Metrics) RecordToolCall(toolName string, duration time.Duration, success bool) {
	if !m.enabled {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.toolCalls[toolName]++
	m.toolDurations[toolName] = append(m.toolDurations[toolName], duration)

	if !success {
		m.toolErrors[toolName]++
	}
}

// GetStats returns formatted statistics
func (m *Metrics) GetStats() string {
	if !m.enabled {
		return "Metrics disabled"
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var result string
	result += "=== MCP Tool Metrics ===\n"

	for tool, count := range m.toolCalls {
		durations := m.toolDurations[tool]
		errors := m.toolErrors[tool]

		if len(durations) == 0 {
			continue
		}

		// Calculate average duration
		var total time.Duration
		for _, d := range durations {
			total += d
		}
		avg := total / time.Duration(len(durations))

		// Calculate min/max
		min := durations[0]
		max := durations[0]
		for _, d := range durations {
			if d < min {
				min = d
			}
			if d > max {
				max = d
			}
		}

		successRate := float64(count-errors) / float64(count) * 100

		result += fmt.Sprintf("\nTool: %s\n", tool)
		result += fmt.Sprintf("  Calls: %d\n", count)
		result += fmt.Sprintf("  Errors: %d\n", errors)
		result += fmt.Sprintf("  Success Rate: %.1f%%\n", successRate)
		result += fmt.Sprintf("  Avg Duration: %v\n", avg)
		result += fmt.Sprintf("  Min Duration: %v\n", min)
		result += fmt.Sprintf("  Max Duration: %v\n", max)
	}

	return result
}

// LogStats logs the current statistics
func (m *Metrics) LogStats() {
	if !m.enabled {
		return
	}
	log.Println(m.GetStats())
}

// Reset clears all metrics
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.toolCalls = make(map[string]int)
	m.toolDurations = make(map[string][]time.Duration)
	m.toolErrors = make(map[string]int)
}
