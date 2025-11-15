package validation

import (
	"fmt"
	"strings"
)

// ConfigValidationResult contains validation results for a config file
type ConfigValidationResult struct {
	Valid    bool
	Errors   []ConfigError
	Warnings []ConfigWarning
}

// ConfigError represents a configuration error
type ConfigError struct {
	Field   string
	Message string
	Value   interface{}
}

// ConfigWarning represents a configuration warning
type ConfigWarning struct {
	Field      string
	Message    string
	Suggestion string
}

// ValidateAPIConfig validates go-zero API service configuration
func ValidateAPIConfig(config map[string]interface{}) *ConfigValidationResult {
	result := &ConfigValidationResult{
		Valid:    true,
		Errors:   []ConfigError{},
		Warnings: []ConfigWarning{},
	}

	// Check required fields
	requiredFields := []string{"Name", "Host", "Port"}
	for _, field := range requiredFields {
		if _, exists := config[field]; !exists {
			result.Valid = false
			result.Errors = append(result.Errors, ConfigError{
				Field:   field,
				Message: fmt.Sprintf("required field '%s' is missing", field),
			})
		}
	}

	// Validate Name
	if name, ok := config["Name"].(string); ok {
		if name == "" {
			result.Valid = false
			result.Errors = append(result.Errors, ConfigError{
				Field:   "Name",
				Message: "service name cannot be empty",
			})
		}
	}

	// Validate Host
	if host, ok := config["Host"].(string); ok {
		if host == "" {
			result.Warnings = append(result.Warnings, ConfigWarning{
				Field:      "Host",
				Message:    "host is empty, will default to 0.0.0.0",
				Suggestion: "Set explicit host value for production",
			})
		}
	}

	// Validate Port
	if port, ok := config["Port"].(float64); ok {
		if port < 1 || port > 65535 {
			result.Valid = false
			result.Errors = append(result.Errors, ConfigError{
				Field:   "Port",
				Message: "port must be between 1 and 65535",
				Value:   port,
			})
		} else if port < 1024 {
			result.Warnings = append(result.Warnings, ConfigWarning{
				Field:      "Port",
				Message:    "using privileged port (< 1024)",
				Suggestion: "Consider using port > 1024 for non-root services",
			})
		}
	} else if portInt, ok := config["Port"].(int); ok {
		if portInt < 1 || portInt > 65535 {
			result.Valid = false
			result.Errors = append(result.Errors, ConfigError{
				Field:   "Port",
				Message: "port must be between 1 and 65535",
				Value:   portInt,
			})
		}
	}

	// Check Timeout settings
	if timeout, ok := config["Timeout"].(float64); ok {
		if timeout <= 0 {
			result.Warnings = append(result.Warnings, ConfigWarning{
				Field:      "Timeout",
				Message:    "timeout is zero or negative",
				Suggestion: "Set a reasonable timeout (e.g., 3000ms)",
			})
		}
	}

	// Check Log settings
	if log, ok := config["Log"].(map[string]interface{}); ok {
		if mode, ok := log["Mode"].(string); ok {
			validModes := []string{"console", "file", "volume"}
			if !contains(validModes, mode) {
				result.Warnings = append(result.Warnings, ConfigWarning{
					Field:      "Log.Mode",
					Message:    fmt.Sprintf("log mode '%s' may not be valid", mode),
					Suggestion: "Use 'console', 'file', or 'volume'",
				})
			}
		}
	}

	// Production environment checks
	if env, ok := config["Environment"].(string); ok {
		if strings.ToLower(env) == "production" || strings.ToLower(env) == "prod" {
			// Check for development settings in production
			if mode, ok := config["Mode"].(string); ok && mode == "dev" {
				result.Warnings = append(result.Warnings, ConfigWarning{
					Field:      "Mode",
					Message:    "development mode in production environment",
					Suggestion: "Set Mode to 'prod' for production",
				})
			}

			// Check for verbose logging
			if log, ok := config["Log"].(map[string]interface{}); ok {
				if level, ok := log["Level"].(string); ok && level == "debug" {
					result.Warnings = append(result.Warnings, ConfigWarning{
						Field:      "Log.Level",
						Message:    "debug logging in production",
						Suggestion: "Use 'info' or 'error' level for production",
					})
				}
			}
		}
	}

	return result
}

// ValidateRPCConfig validates go-zero RPC service configuration
func ValidateRPCConfig(config map[string]interface{}) *ConfigValidationResult {
	result := &ConfigValidationResult{
		Valid:    true,
		Errors:   []ConfigError{},
		Warnings: []ConfigWarning{},
	}

	// Check required fields for RPC
	requiredFields := []string{"Name", "ListenOn"}
	for _, field := range requiredFields {
		if _, exists := config[field]; !exists {
			result.Valid = false
			result.Errors = append(result.Errors, ConfigError{
				Field:   field,
				Message: fmt.Sprintf("required field '%s' is missing", field),
			})
		}
	}

	// Validate ListenOn format (host:port)
	if listenOn, ok := config["ListenOn"].(string); ok {
		if !strings.Contains(listenOn, ":") {
			result.Valid = false
			result.Errors = append(result.Errors, ConfigError{
				Field:   "ListenOn",
				Message: "ListenOn must be in format 'host:port'",
				Value:   listenOn,
			})
		}
	}

	// Check Etcd settings if present
	if etcd, ok := config["Etcd"].(map[string]interface{}); ok {
		if hosts, ok := etcd["Hosts"].([]interface{}); ok {
			if len(hosts) == 0 {
				result.Warnings = append(result.Warnings, ConfigWarning{
					Field:      "Etcd.Hosts",
					Message:    "etcd hosts list is empty",
					Suggestion: "Provide at least one etcd host for service discovery",
				})
			}
		}
	}

	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
