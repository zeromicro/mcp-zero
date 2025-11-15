package errors

import (
	"errors"
	"fmt"
)

// Common error types for mcp-zero project

var (
	ErrGoctlNotFound      = errors.New("goctl not found")
	ErrGoctlExecution     = errors.New("goctl execution failed")
	ErrInvalidServiceName = errors.New("invalid service name")
	ErrInvalidPort        = errors.New("invalid port")
	ErrInvalidPath        = errors.New("invalid path")
	ErrValidationFailed   = errors.New("validation failed")
	ErrProjectStructure   = errors.New("invalid project structure")
	ErrModuleInit         = errors.New("module initialization failed")
	ErrImportFix          = errors.New("import path fix failed")
	ErrConfigUpdate       = errors.New("config update failed")
)

// ValidationError represents a validation failure
type ValidationError struct {
	Field   string
	Value   string
	Message string
	Err     error
}

func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (value: %s): %v", e.Field, e.Message, e.Value, e.Err)
	}
	return fmt.Sprintf("%s: %s (value: %s)", e.Field, e.Message, e.Value)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

func NewValidationError(field, value, message string, err error) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
		Err:     err,
	}
}

// ExecutionError represents a goctl execution failure
type ExecutionError struct {
	Command  string
	Args     []string
	Stdout   string
	Stderr   string
	ExitCode int
	Err      error
}

func (e *ExecutionError) Error() string {
	msg := fmt.Sprintf("goctl execution failed: %s %v (exit code: %d)", e.Command, e.Args, e.ExitCode)
	if e.Stderr != "" {
		msg += fmt.Sprintf("\nstderr: %s", e.Stderr)
	}
	if e.Err != nil {
		msg += fmt.Sprintf("\nerror: %v", e.Err)
	}
	return msg
}

func (e *ExecutionError) Unwrap() error {
	return e.Err
}

func NewExecutionError(command string, args []string, stdout, stderr string, exitCode int, err error) *ExecutionError {
	return &ExecutionError{
		Command:  command,
		Args:     args,
		Stdout:   stdout,
		Stderr:   stderr,
		ExitCode: exitCode,
		Err:      err,
	}
}

// PathError represents a file/directory path error
type PathError struct {
	Path    string
	Op      string
	Message string
	Err     error
}

func (e *PathError) Error() string {
	msg := fmt.Sprintf("path error (%s): %s - %s", e.Op, e.Path, e.Message)
	if e.Err != nil {
		msg += fmt.Sprintf(": %v", e.Err)
	}
	return msg
}

func (e *PathError) Unwrap() error {
	return e.Err
}

func NewPathError(path, op, message string, err error) *PathError {
	return &PathError{
		Path:    path,
		Op:      op,
		Message: message,
		Err:     err,
	}
}

func IsNotFound(err error) bool {
	var pathErr *PathError
	if errors.As(err, &pathErr) {
		return pathErr.Op == "stat" || pathErr.Op == "read"
	}
	return errors.Is(err, ErrGoctlNotFound)
}

func IsValidationError(err error) bool {
	var valErr *ValidationError
	return errors.As(err, &valErr) || errors.Is(err, ErrValidationFailed)
}

func IsExecutionError(err error) bool {
	var execErr *ExecutionError
	return errors.As(err, &execErr) || errors.Is(err, ErrGoctlExecution)
}
