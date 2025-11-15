package validation_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zeromicro/mcp-zero/internal/validation"
)

func TestValidateServiceName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid names
		{"simple lowercase", "userservice", false},
		{"with uppercase", "UserService", false},
		{"with underscore", "user_service", false},
		{"with numbers", "user123", false},
		{"mixed", "myService_v2", false},

		// Invalid names
		{"empty name", "", true},
		{"starts with number", "123service", true},
		{"contains hyphen", "user-service", true},
		{"multiple hyphens", "my-user-service", true},
		{"starts with underscore", "_userservice", true},
		{"contains space", "user service", true},
		{"special characters", "user@service", true},
		{"dots", "user.service", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateServiceName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateServiceName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestSuggestServiceName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"replace hyphen", "user-service", "user_service"},
		{"multiple hyphens", "my-user-service", "my_user_service"},
		{"starts with number", "123service", "s123service"},
		{"special chars", "user@service!", "userservice"},
		{"mixed invalid", "1user-service!", "s1user_service"},
		{"already valid", "userservice", "userservice"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validation.SuggestServiceName(tt.input)
			if result != tt.expected {
				t.Errorf("SuggestServiceName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		// Valid ports
		{"standard port", 8080, false},
		{"high port", 9000, false},
		{"minimum valid", 1024, false},
		{"maximum valid", 65535, false},

		// Invalid ports
		{"too low", 80, true},
		{"privileged", 443, true},
		{"below minimum", 1023, true},
		{"above maximum", 65536, true},
		{"way too high", 100000, true},
		{"negative", -1, true},
		{"zero", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidatePort(tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePort(%d) error = %v, wantErr %v", tt.port, err, tt.wantErr)
			}
		})
	}
}

func TestSuggestAvailablePort(t *testing.T) {
	// Test that it returns a port in the expected range
	port, err := validation.SuggestAvailablePort(9000)
	if err != nil {
		t.Fatalf("SuggestAvailablePort(9000) failed: %v", err)
	}

	if port < 9000 || port > 9100 {
		t.Errorf("SuggestAvailablePort(9000) = %d, want port in range [9000, 9100]", port)
	}
}

func TestValidatePath(t *testing.T) {
	// Create a temp directory for testing
	tmpDir, err := os.MkdirTemp("", "mcp-zero-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid absolute path", filepath.Join(tmpDir, "testfile"), false},
		{"relative path", "relative/path", true},
		{"dot path", "./file", true},
		{"double dot path", "../file", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidatePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}

func TestValidateOutputDir(t *testing.T) {
	// Create a temp directory for testing
	tmpDir, err := os.MkdirTemp("", "mcp-zero-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
	}{
		{
			name: "existing writable directory",
			setup: func() string {
				return tmpDir
			},
			wantErr: false,
		},
		{
			name: "non-existent directory (should create)",
			setup: func() string {
				return filepath.Join(tmpDir, "newdir")
			},
			wantErr: false,
		},
		{
			name: "relative path",
			setup: func() string {
				return "relative/path"
			},
			wantErr: true,
		},
		{
			name: "file instead of directory",
			setup: func() string {
				file := filepath.Join(tmpDir, "testfile")
				os.WriteFile(file, []byte("test"), 0644)
				return file
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setup()
			err := validation.ValidateOutputDir(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOutputDir(%q) error = %v, wantErr %v", dir, err, tt.wantErr)
			}

			// If successful and directory was created, verify it exists
			if err == nil && tt.name == "non-existent directory (should create)" {
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					t.Errorf("ValidateOutputDir should have created directory %q", dir)
				}
			}
		})
	}
}

func TestEnsureDirectoryExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mcp-zero-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
	}{
		{
			name: "existing directory",
			setup: func() string {
				dir := filepath.Join(tmpDir, "existing")
				os.MkdirAll(dir, 0755)
				return dir
			},
			wantErr: false,
		},
		{
			name: "non-existent directory",
			setup: func() string {
				return filepath.Join(tmpDir, "new", "nested", "dir")
			},
			wantErr: false,
		},
		{
			name: "file instead of directory",
			setup: func() string {
				file := filepath.Join(tmpDir, "file")
				os.WriteFile(file, []byte("test"), 0644)
				return file
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setup()
			err := validation.EnsureDirectoryExists(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureDirectoryExists(%q) error = %v, wantErr %v", dir, err, tt.wantErr)
			}

			// If successful, verify directory exists
			if err == nil {
				info, statErr := os.Stat(dir)
				if statErr != nil {
					t.Errorf("Directory %q should exist but stat failed: %v", dir, statErr)
				} else if !info.IsDir() {
					t.Errorf("Path %q should be a directory", dir)
				}
			}
		})
	}
}
