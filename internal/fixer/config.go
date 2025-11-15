package fixer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// UpdateConfigFile updates configuration files with correct settings
func UpdateConfigFile(projectPath string, serviceName string, port int) error {
	// Find config file in etc/ directory
	etcDir := filepath.Join(projectPath, "etc")
	configFiles, err := filepath.Glob(filepath.Join(etcDir, "*.yaml"))
	if err != nil {
		return fmt.Errorf("failed to search for config files: %w", err)
	}

	if len(configFiles) == 0 {
		// No config file found, skip
		return nil
	}

	// Update first config file found
	configFile := configFiles[0]
	content, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	configContent := string(content)

	// Update service name if needed
	if strings.Contains(configContent, "Name:") && !strings.Contains(configContent, fmt.Sprintf("Name: %s", serviceName)) {
		configContent = updateYAMLField(configContent, "Name", serviceName)
	}

	// Update port if specified and field exists
	if port > 0 && strings.Contains(configContent, "Port:") {
		configContent = updateYAMLField(configContent, "Port", fmt.Sprintf("%d", port))
	}

	// Write updated config
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// updateYAMLField updates a field in YAML content
func updateYAMLField(content string, fieldName string, newValue string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, fieldName+":") {
			// Preserve indentation
			indent := ""
			for _, ch := range line {
				if ch == ' ' || ch == '\t' {
					indent += string(ch)
				} else {
					break
				}
			}
			lines[i] = fmt.Sprintf("%s%s: %s", indent, fieldName, newValue)
			break
		}
	}
	return strings.Join(lines, "\n")
}
