package fixer

import (
"fmt"
"os"
"path/filepath"
"regexp"
"strings"
)

// FixImports fixes import paths in generated code
// Replaces absolute paths with local module names
func FixImports(projectPath string, moduleName string) error {
	// Find all .go files
	var goFiles []string
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
if err != nil {
return err
}
if !info.IsDir() && strings.HasSuffix(path, ".go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk project directory: %w", err)
	}

	// Fix imports in each file
	for _, goFile := range goFiles {
		if err := fixImportsInFile(goFile, projectPath, moduleName); err != nil {
			return fmt.Errorf("failed to fix imports in %s: %w", goFile, err)
		}
	}

	return nil
}

// fixImportsInFile fixes imports in a single Go file
func fixImportsInFile(filePath string, projectPath string, moduleName string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	originalContent := string(content)
	modifiedContent := originalContent

	// Fix absolute path imports
	// Pattern: import "absolute/path/to/project/internal/..."
	// Replace with: import "moduleName/internal/..."
	absPathPattern := regexp.MustCompile(regexp.QuoteMeta(projectPath) + `(/[^"]+)`)
	modifiedContent = absPathPattern.ReplaceAllString(modifiedContent, moduleName+"$1")

	// Only write if content changed
	if modifiedContent != originalContent {
		if err := os.WriteFile(filePath, []byte(modifiedContent), 0644); err != nil {
			return err
		}
	}

	return nil
}
