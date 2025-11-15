package analyzer

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type RPCService struct {
	FilePath    string
	ServiceName string
	Methods     []RPCMethod
	Messages    []string
}

type RPCMethod struct {
	Name     string
	Request  string
	Response string
	Stream   string
}

func ParseProtoSpecification(protoFile string) (*RPCService, error) {
	content, err := os.ReadFile(protoFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read proto file: %w", err)
	}

	service := &RPCService{
		FilePath: protoFile,
	}

	fileContent := string(content)

	service.ServiceName = extractRPCServiceName(fileContent)
	if service.ServiceName == "" {
		return nil, fmt.Errorf("no service name found")
	}

	service.Methods = extractRPCMethods(fileContent)
	service.Messages = extractMessages(fileContent)

	return service, nil
}

func extractRPCServiceName(content string) string {
	serviceRegex := regexp.MustCompile(`service\s+([A-Z][a-zA-Z0-9_]*)\s*{`)
	matches := serviceRegex.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractRPCMethods(content string) []RPCMethod {
	var methods []RPCMethod

	serviceRegex := regexp.MustCompile(`service\s+[A-Z][a-zA-Z0-9_]*\s*{([^}]+)}`)
	serviceMatches := serviceRegex.FindStringSubmatch(content)
	if len(serviceMatches) < 2 {
		return methods
	}

	serviceBlock := serviceMatches[1]
	methodRegex := regexp.MustCompile(`rpc\s+([A-Z][a-zA-Z0-9_]*)\s*\(([^)]+)\)\s*returns\s*\(([^)]+)\)`)
	matches := methodRegex.FindAllStringSubmatch(serviceBlock, -1)

	for _, match := range matches {
		if len(match) > 3 {
			method := RPCMethod{
				Name:     match[1],
				Request:  strings.TrimSpace(match[2]),
				Response: strings.TrimSpace(match[3]),
			}

			if strings.Contains(match[2], "stream") {
				method.Stream = "request"
			}
			if strings.Contains(match[3], "stream") {
				if method.Stream == "request" {
					method.Stream = "bidirectional"
				} else {
					method.Stream = "response"
				}
			}

			methods = append(methods, method)
		}
	}

	return methods
}

func extractMessages(content string) []string {
	var messages []string

	messageRegex := regexp.MustCompile(`message\s+([A-Z][a-zA-Z0-9_]*)\s*{`)
	matches := messageRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			messages = append(messages, match[1])
		}
	}

	return messages
}
