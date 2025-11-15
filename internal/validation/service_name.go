package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// ValidateServiceName validates service name against go-zero requirements
// Service names must be valid Go identifiers: start with letter, no hyphens
func ValidateServiceName(name string) error {
	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	// Check if starts with letter
	firstRune := rune(name[0])
	if !unicode.IsLetter(firstRune) {
		return fmt.Errorf("service name must start with a letter, got '%c'", firstRune)
	}

	// Check for hyphens (common mistake)
	if strings.Contains(name, "-") {
		suggestion := strings.ReplaceAll(name, "-", "_")
		return fmt.Errorf("service name cannot contain hyphens, try: %s", suggestion)
	}

	// Check if valid Go identifier
	validIdentifier := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	if !validIdentifier.MatchString(name) {
		return fmt.Errorf("service name must be a valid Go identifier (letters, numbers, underscores only)")
	}

	return nil
}

// SuggestServiceName suggests a valid service name from an invalid one
func SuggestServiceName(name string) string {
	// Replace hyphens with underscores
	suggested := strings.ReplaceAll(name, "-", "_")

	// Remove invalid characters
	var result strings.Builder
	for i, r := range suggested {
		if i == 0 {
			if unicode.IsLetter(r) {
				result.WriteRune(r)
			} else {
				result.WriteRune('s')
				if unicode.IsDigit(r) || r == '_' {
					result.WriteRune(r)
				}
			}
		} else {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				result.WriteRune(r)
			}
		}
	}

	return result.String()
}
