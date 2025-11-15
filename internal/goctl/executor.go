package goctl

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
)

// Executor handles safe execution of goctl commands
type Executor struct {
	goctlPath string
}

// NewExecutor creates a new goctl executor
// Discovers goctl path on initialization
func NewExecutor() (*Executor, error) {
	goctlPath, err := DiscoverGoctl()
	if err != nil {
		return nil, err
	}

	return &Executor{
		goctlPath: goctlPath,
	}, nil
}

// ExecuteResult contains the result of a goctl command execution
type ExecuteResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Error    error
}

// Execute runs a goctl command with the given arguments
// Uses absolute paths and captures both stdout and stderr
func (e *Executor) Execute(args ...string) *ExecuteResult {
	cmd := exec.Command(e.goctlPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := &ExecuteResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if err != nil {
		result.Error = err
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
	}

	return result
}

// ExecuteInDir runs a goctl command in a specific directory
func (e *Executor) ExecuteInDir(dir string, args ...string) *ExecuteResult {
	// Ensure directory is absolute
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return &ExecuteResult{
			Error:    fmt.Errorf("failed to get absolute path: %w", err),
			ExitCode: -1,
		}
	}

	cmd := exec.Command(e.goctlPath, args...)
	cmd.Dir = absDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	result := &ExecuteResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if err != nil {
		result.Error = err
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
	}

	return result
}

// GetPath returns the discovered goctl path
func (e *Executor) GetPath() string {
	return e.goctlPath
}
