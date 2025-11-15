package fixer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zeromicro/mcp-zero/internal/fixer"
)

func TestFixImportsInFile(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "mcp-zero-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test Go file with absolute path imports
	testFile := filepath.Join(tmpDir, "test.go")
	testContent := `package main

import (
	"fmt"
	"` + tmpDir + `/internal/config"
	"` + tmpDir + `/pkg/utils"
)

func main() {
	fmt.Println("test")
}
`

	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Fix imports
	moduleName := "github.com/test/project"
	err = fixer.FixImports(tmpDir, moduleName)
	if err != nil {
		t.Fatalf("FixImports() failed: %v", err)
	}

	// Read modified file
	modifiedContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	modified := string(modifiedContent)

	// Verify absolute paths were replaced
	if strings.Contains(modified, tmpDir) {
		t.Errorf("File still contains absolute path %q:\n%s", tmpDir, modified)
	}

	// Verify module name is present
	if !strings.Contains(modified, moduleName) {
		t.Errorf("File does not contain module name %q:\n%s", moduleName, modified)
	}

	// Verify specific imports were fixed
	expectedImports := []string{
		moduleName + "/internal/config",
		moduleName + "/pkg/utils",
	}

	for _, expectedImport := range expectedImports {
		if !strings.Contains(modified, expectedImport) {
			t.Errorf("Expected import %q not found in file", expectedImport)
		}
	}
}

func TestInitializeGoModule(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mcp-zero-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	moduleName := "github.com/test/testmodule"

	// Initialize module
	err = fixer.InitializeGoModule(tmpDir, moduleName)
	if err != nil {
		t.Fatalf("InitializeGoModule() failed: %v", err)
	}

	// Verify go.mod was created
	goModPath := filepath.Join(tmpDir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Error("go.mod file was not created")
	}

	// Verify go.mod contains module name
	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), moduleName) {
		t.Errorf("go.mod does not contain module name %q", moduleName)
	}

	// Test idempotency - calling again should not fail
	err = fixer.InitializeGoModule(tmpDir, moduleName)
	if err != nil {
		t.Errorf("InitializeGoModule() should be idempotent but failed on second call: %v", err)
	}
}

func TestTidyGoModule(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mcp-zero-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a simple Go module
	moduleName := "github.com/test/tidytest"
	err = fixer.InitializeGoModule(tmpDir, moduleName)
	if err != nil {
		t.Fatal(err)
	}

	// Create a simple main.go
	mainGo := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(mainGo, []byte(`package main

import "fmt"

func main() {
	fmt.Println("test")
}
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Run tidy
	err = fixer.TidyGoModule(tmpDir)
	if err != nil {
		t.Fatalf("TidyGoModule() failed: %v", err)
	}

	// Verify go.mod still exists
	goModPath := filepath.Join(tmpDir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Error("go.mod file disappeared after tidy")
	}
}

func TestVerifyBuild(t *testing.T) {
	// Test case 1: Valid buildable project
	t.Run("valid project", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "mcp-zero-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		// Create a simple valid Go project
		moduleName := "github.com/test/buildtest"
		err = fixer.InitializeGoModule(tmpDir, moduleName)
		if err != nil {
			t.Fatal(err)
		}

		mainGo := filepath.Join(tmpDir, "main.go")
		err = os.WriteFile(mainGo, []byte(`package main

import "fmt"

func main() {
	fmt.Println("test")
}
`), 0644)
		if err != nil {
			t.Fatal(err)
		}

		err = fixer.TidyGoModule(tmpDir)
		if err != nil {
			t.Fatal(err)
		}

		// Verify build succeeds
		err = fixer.VerifyBuild(tmpDir)
		if err != nil {
			t.Errorf("VerifyBuild() should succeed for valid project, got error: %v", err)
		}
	})

	// Test case 2: Invalid project with syntax error
	t.Run("invalid project", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "mcp-zero-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		// Create project with syntax error
		moduleName := "github.com/test/invalidtest"
		err = fixer.InitializeGoModule(tmpDir, moduleName)
		if err != nil {
			t.Fatal(err)
		}

		mainGo := filepath.Join(tmpDir, "main.go")
		err = os.WriteFile(mainGo, []byte(`package main

func main() {
	// Missing closing brace
`), 0644)
		if err != nil {
			t.Fatal(err)
		}

		// Verify build fails
		err = fixer.VerifyBuild(tmpDir)
		if err == nil {
			t.Error("VerifyBuild() should fail for invalid project")
		}
	})
}

func TestGetGoModuleName(t *testing.T) {
	tests := []struct {
		name         string
		goModContent string
		wantModule   string
		wantErr      bool
	}{
		{
			name:         "simple module name",
			goModContent: "module github.com/test/project\n\ngo 1.19\n",
			wantModule:   "github.com/test/project",
			wantErr:      false,
		},
		{
			name:         "module name without newline",
			goModContent: "module mymodule",
			wantModule:   "mymodule",
			wantErr:      false,
		},
		{
			name:         "missing module declaration",
			goModContent: "go 1.19\n",
			wantModule:   "",
			wantErr:      true,
		},
		{
			name:         "empty file",
			goModContent: "",
			wantModule:   "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "mcp-zero-test-*")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			// Create go.mod with test content
			goModPath := filepath.Join(tmpDir, "go.mod")
			err = os.WriteFile(goModPath, []byte(tt.goModContent), 0644)
			if err != nil {
				t.Fatal(err)
			}

			// Test GetGoModuleName
			moduleName, err := fixer.GetGoModuleName(tmpDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGoModuleName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if moduleName != tt.wantModule {
				t.Errorf("GetGoModuleName() = %q, want %q", moduleName, tt.wantModule)
			}
		})
	}

	// Test case: missing go.mod file
	t.Run("missing go.mod", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "mcp-zero-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		_, err = fixer.GetGoModuleName(tmpDir)
		if err == nil {
			t.Error("GetGoModuleName() should fail when go.mod is missing")
		}
	})
}

func TestFixImportsWithNestedDirectories(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mcp-zero-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create nested directory structure
	internalDir := filepath.Join(tmpDir, "internal", "handler")
	err = os.MkdirAll(internalDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create multiple Go files in different directories
	files := map[string]string{
		filepath.Join(tmpDir, "main.go"): `package main

import "` + tmpDir + `/internal/handler"

func main() {}
`,
		filepath.Join(internalDir, "handler.go"): `package handler

import "` + tmpDir + `/internal/config"

func Handle() {}
`,
	}

	for path, content := range files {
		err = os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Fix imports
	moduleName := "example.com/myproject"
	err = fixer.FixImports(tmpDir, moduleName)
	if err != nil {
		t.Fatalf("FixImports() failed: %v", err)
	}

	// Verify all files were fixed
	for path := range files {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}

		if strings.Contains(string(content), tmpDir) {
			t.Errorf("File %s still contains absolute path", path)
		}

		if !strings.Contains(string(content), moduleName) {
			t.Errorf("File %s does not contain module name", path)
		}
	}
}
