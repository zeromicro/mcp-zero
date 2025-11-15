package analyzer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zeromicro/mcp-zero/internal/analyzer"
)

func TestScanProjectValidProject(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mcp-zero-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	goModContent := "module github.com/test/project\n\ngo 1.19\n"
	err = os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	analysis, err := analyzer.ScanProject(tmpDir)
	if err != nil {
		t.Fatalf("ScanProject() failed: %v", err)
	}

	if analysis.ProjectPath != tmpDir {
		t.Errorf("ProjectPath = %q, want %q", analysis.ProjectPath, tmpDir)
	}
}

func TestScanProjectInvalidPath(t *testing.T) {
	_, err := analyzer.ScanProject("/nonexistent/path")
	if err == nil {
		t.Error("ScanProject() should fail for non-existent path")
	}
}

func TestScanProjectWithAPIService(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mcp-zero-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	apiDir := filepath.Join(tmpDir, "userservice")
	err = os.MkdirAll(apiDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	apiFile := filepath.Join(apiDir, "user.api")
	apiContent := `syntax = "v1"

service UserService {
	@handler GetUser
	get /api/user/:id
}
`
	err = os.WriteFile(apiFile, []byte(apiContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	analysis, err := analyzer.ScanProject(tmpDir)
	if err != nil {
		t.Fatalf("ScanProject() failed: %v", err)
	}

	if analysis.Summary.APIServices < 1 {
		t.Error("Expected at least 1 API service to be discovered")
	}
}

func TestScanProjectWithConfigFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mcp-zero-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configFiles := map[string]string{
		"config.yaml":   "server:\n  port: 8080\n",
		"settings.json": `{"database": "mysql"}`,
		"etc.yaml":      "timeout: 30s\n",
	}

	for filename, content := range configFiles {
		path := filepath.Join(tmpDir, filename)
		err = os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	analysis, err := analyzer.ScanProject(tmpDir)
	if err != nil {
		t.Fatalf("ScanProject() failed: %v", err)
	}

	if len(analysis.Configs) != 3 {
		t.Errorf("Expected 3 config files, got %d", len(analysis.Configs))
	}
}

func TestScanProjectWithDependencies(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mcp-zero-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	goModContent := `module github.com/test/project

go 1.19

require (
	github.com/zeromicro/go-zero v1.4.0
	github.com/golang/protobuf v1.5.2
)
`
	goModPath := filepath.Join(tmpDir, "go.mod")
	err = os.WriteFile(goModPath, []byte(goModContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	analysis, err := analyzer.ScanProject(tmpDir)
	if err != nil {
		t.Fatalf("ScanProject() failed: %v", err)
	}

	if len(analysis.Dependencies) < 2 {
		t.Errorf("Expected at least 2 dependencies, got %d", len(analysis.Dependencies))
	}

	if analysis.Summary.GoZeroVersion == "" {
		t.Error("go-zero version should be detected")
	}
}
