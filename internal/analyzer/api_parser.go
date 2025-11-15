package analyzer

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type APISpecification struct {
	FilePath    string
	ServiceName string
	Endpoints   []Endpoint
	Types       []string
}

type Endpoint struct {
	Method   string
	Path     string
	Handler  string
	Request  string
	Response string
}

func ParseAPISpecification(apiFile string) (*APISpecification, error) {
	content, err := os.ReadFile(apiFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read API file: %w", err)
	}
	spec := &APISpecification{FilePath: apiFile}
	fileContent := string(content)
	spec.ServiceName = extractServiceName(fileContent)
	if spec.ServiceName == "" {
		return nil, fmt.Errorf("no service name found")
	}
	spec.Endpoints = extractEndpoints(fileContent)
	spec.Types = extractTypes(fileContent)
	return spec, nil
}

func extractServiceName(content string) string {
	serviceRegex := regexp.MustCompile(`service\s+([a-zA-Z0-9_-]+)\s*{`)
	matches := serviceRegex.FindStringSubmatch(content)
	if len(matches) > 1 {
		return strings.TrimSuffix(matches[1], "-api")
	}
	return ""
}

func extractEndpoints(content string) []Endpoint {
	var endpoints []Endpoint
	serviceRegex := regexp.MustCompile(`service\s+[a-zA-Z0-9_-]+\s*{([^}]+)}`)
	serviceMatches := serviceRegex.FindStringSubmatch(content)
	if len(serviceMatches) < 2 {
		return endpoints
	}
	serviceBlock := serviceMatches[1]
	endpointRegex := regexp.MustCompile(`@handler\s+(\w+)\s+(\w+)\s+(/[^\s(]*)\s*(?:\(([^)]+)\))?\s*(?:returns\s*\(([^)]+)\))?`)
	matches := endpointRegex.FindAllStringSubmatch(serviceBlock, -1)
	for _, match := range matches {
		if len(match) > 3 {
			endpoint := Endpoint{
				Handler: match[1],
				Method:  strings.ToUpper(match[2]),
				Path:    match[3],
			}
			if len(match) > 4 {
				endpoint.Request = strings.TrimSpace(match[4])
			}
			if len(match) > 5 {
				endpoint.Response = strings.TrimSpace(match[5])
			}
			endpoints = append(endpoints, endpoint)
		}
	}
	return endpoints
}

func extractTypes(content string) []string {
	var types []string
	typeRegex := regexp.MustCompile(`type\s+(\w+)\s+(?:struct\s*)?{`)
	matches := typeRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			types = append(types, match[1])
		}
	}
	return types
}
