package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ProjectAnalysis represents a comprehensive analysis of a go-zero project
type ProjectAnalysis struct {
	ProjectPath  string
	Services     []ServiceInfo
	Dependencies []Dependency
	Configs      []ConfigFile
	Summary      ProjectSummary
}

// ServiceInfo represents information about a service in the project
type ServiceInfo struct {
	Name       string
	Type       string // "api", "rpc", or "model"
	Path       string
	SpecFile   string
	Endpoints  []EndpointInfo
	RPCMethods []RPCMethodInfo
}

// EndpointInfo represents an API endpoint
type EndpointInfo struct {
	Method  string
	Path    string
	Handler string
}

// RPCMethodInfo represents an RPC method
type RPCMethodInfo struct {
	Name     string
	Request  string
	Response string
	Stream   bool
}

// Dependency represents a project dependency
type Dependency struct {
	Name    string
	Version string
	Type    string // "direct" or "indirect"
}

// ConfigFile represents a configuration file
type ConfigFile struct {
	Path string
	Type string // "yaml", "json", "toml"
}

// ProjectSummary provides high-level statistics
type ProjectSummary struct {
	TotalServices     int
	APIServices       int
	RPCServices       int
	ModelServices     int
	TotalEndpoints    int
	TotalRPCMethods   int
	TotalDependencies int
	GoZeroVersion     string
}

// ScanProject analyzes a go-zero project directory
func ScanProject(projectPath string) (*ProjectAnalysis, error) {
	if !filepath.IsAbs(projectPath) {
		var err error
		projectPath, err = filepath.Abs(projectPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve project path: %w", err)
		}
	}

	info, err := os.Stat(projectPath)
	if err != nil {
		return nil, fmt.Errorf("project path does not exist: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("project path is not a directory")
	}

	analysis := &ProjectAnalysis{
		ProjectPath:  projectPath,
		Services:     []ServiceInfo{},
		Dependencies: []Dependency{},
		Configs:      []ConfigFile{},
	}

	// Discover API services
	apiFiles, err := discoverAPIFiles(projectPath)
	if err == nil {
		for _, apiFile := range apiFiles {
			service := ServiceInfo{
				Type:      "api",
				Path:      filepath.Dir(apiFile),
				SpecFile:  apiFile,
				Endpoints: []EndpointInfo{},
			}

			// Parse API spec to extract endpoints
			if spec, err := ParseAPISpecification(apiFile); err == nil {
				service.Name = spec.ServiceName
				for _, endpoint := range spec.Endpoints {
					service.Endpoints = append(service.Endpoints, EndpointInfo{
						Method:  endpoint.Method,
						Path:    endpoint.Path,
						Handler: endpoint.Handler,
					})
				}
			}

			analysis.Services = append(analysis.Services, service)
			analysis.Summary.APIServices++
			analysis.Summary.TotalEndpoints += len(service.Endpoints)
		}
	}

	// Discover RPC services
	protoFiles, err := discoverProtoFiles(projectPath)
	if err == nil {
		for _, protoFile := range protoFiles {
			service := ServiceInfo{
				Type:       "rpc",
				Path:       filepath.Dir(protoFile),
				SpecFile:   protoFile,
				RPCMethods: []RPCMethodInfo{},
			}

			// Parse proto spec to extract methods
			if spec, err := ParseProtoSpecification(protoFile); err == nil {
				service.Name = spec.ServiceName
				for _, method := range spec.Methods {
					isStream := strings.Contains(method.Stream, "stream")
					service.RPCMethods = append(service.RPCMethods, RPCMethodInfo{
						Name:     method.Name,
						Request:  method.Request,
						Response: method.Response,
						Stream:   isStream,
					})
				}
			}

			analysis.Services = append(analysis.Services, service)
			analysis.Summary.RPCServices++
			analysis.Summary.TotalRPCMethods += len(service.RPCMethods)
		}
	}

	// Discover config files
	configs, err := discoverConfigFiles(projectPath)
	if err == nil {
		analysis.Configs = configs
	}

	// Parse dependencies
	deps, goZeroVersion, err := parseDependencies(projectPath)
	if err == nil {
		analysis.Dependencies = deps
		analysis.Summary.TotalDependencies = len(deps)
		analysis.Summary.GoZeroVersion = goZeroVersion
	}

	analysis.Summary.TotalServices = len(analysis.Services)
	analysis.Summary.ModelServices = 0 // Models don't have spec files typically

	return analysis, nil
}

// discoverAPIFiles finds all .api files in the project
func discoverAPIFiles(projectPath string) ([]string, error) {
	var apiFiles []string

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue on error
		}

		// Skip hidden directories and vendor/node_modules
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasSuffix(path, ".api") {
			apiFiles = append(apiFiles, path)
		}

		return nil
	})

	return apiFiles, err
}

// discoverProtoFiles finds all .proto files in the project
func discoverProtoFiles(projectPath string) ([]string, error) {
	var protoFiles []string

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasSuffix(path, ".proto") {
			protoFiles = append(protoFiles, path)
		}

		return nil
	})

	return protoFiles, err
}

// discoverConfigFiles finds all config files (yaml, json, toml)
func discoverConfigFiles(projectPath string) ([]ConfigFile, error) {
	var configs []ConfigFile

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		var configType string
		switch ext {
		case ".yaml", ".yml":
			configType = "yaml"
		case ".json":
			configType = "json"
		case ".toml":
			configType = "toml"
		default:
			return nil
		}

		// Only include common config file names
		name := strings.ToLower(filepath.Base(path))
		if strings.Contains(name, "config") || strings.Contains(name, "settings") ||
			name == "etc.yaml" || name == "etc.json" {
			configs = append(configs, ConfigFile{
				Path: path,
				Type: configType,
			})
		}

		return nil
	})

	return configs, err
}

// parseDependencies extracts dependencies from go.mod
func parseDependencies(projectPath string) ([]Dependency, string, error) {
	goModPath := filepath.Join(projectPath, "go.mod")

	content, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	var deps []Dependency
	goZeroVersion := ""

	lines := strings.Split(string(content), "\n")
	inRequire := false

	directDepRegex := regexp.MustCompile(`^\s*([a-zA-Z0-9\-\._/]+)\s+v([0-9\.\-+a-zA-Z]+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "require (") {
			inRequire = true
			continue
		}
		if line == ")" {
			inRequire = false
			continue
		}

		if strings.HasPrefix(line, "require ") {
			// Single line require
			line = strings.TrimPrefix(line, "require ")
			matches := directDepRegex.FindStringSubmatch(line)
			if len(matches) == 3 {
				dep := Dependency{
					Name:    matches[1],
					Version: matches[2],
					Type:    "direct",
				}
				deps = append(deps, dep)

				if strings.Contains(dep.Name, "go-zero") {
					goZeroVersion = dep.Version
				}
			}
			continue
		}

		if inRequire {
			matches := directDepRegex.FindStringSubmatch(line)
			if len(matches) == 3 {
				depType := "direct"
				if strings.Contains(line, "// indirect") {
					depType = "indirect"
				}

				dep := Dependency{
					Name:    matches[1],
					Version: matches[2],
					Type:    depType,
				}
				deps = append(deps, dep)

				if strings.Contains(dep.Name, "go-zero") && goZeroVersion == "" {
					goZeroVersion = dep.Version
				}
			}
		}
	}

	return deps, goZeroVersion, nil
}
