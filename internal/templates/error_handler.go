package templates

import "fmt"

const BasicErrorHandlerTemplate = `package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type ErrorResponse struct {
	Code    int    ` + "`json:\"code\"`" + `
	Message string ` + "`json:\"message\"`" + `
}

func ErrorHandler(err error) (int, interface{}) {
	switch err.(type) {
	default:
		return http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
}
`

const DetailedErrorHandlerTemplate = `package handler

import (
	"net/http"
	"time"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type DetailedErrorResponse struct {
	Code      int    ` + "`json:\"code\"`" + `
	Message   string ` + "`json:\"message\"`" + `
	Details   string ` + "`json:\"details,omitempty\"`" + `
	Timestamp int64  ` + "`json:\"timestamp\"`" + `
	Path      string ` + "`json:\"path,omitempty\"`" + `
}

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type BusinessError struct {
	Code    int
	Message string
}

func (e BusinessError) Error() string {
	return e.Message
}

func ErrorHandler(err error) (int, interface{}) {
	logx.Errorf("Error occurred: %v", err)
	timestamp := time.Now().Unix()

	switch e := err.(type) {
	case ValidationError:
		return http.StatusBadRequest, DetailedErrorResponse{
			Code:      http.StatusBadRequest,
			Message:   "Validation failed",
			Details:   e.Error(),
			Timestamp: timestamp,
		}
	case BusinessError:
		return e.Code, DetailedErrorResponse{
			Code:      e.Code,
			Message:   e.Message,
			Timestamp: timestamp,
		}
	default:
		return http.StatusInternalServerError, DetailedErrorResponse{
			Code:      http.StatusInternalServerError,
			Message:   "Internal server error",
			Timestamp: timestamp,
		}
	}
}
`

func GetErrorHandlerTemplate(name string) (*Template, error) {
	switch name {
	case "basic":
		return &Template{
			Name:        "basic",
			Type:        "error_handler",
			Description: "Basic error handler with simple error responses",
			Content:     BasicErrorHandlerTemplate,
			Parameters:  []TemplateParameter{},
		}, nil
	case "detailed":
		return &Template{
			Name:        "detailed",
			Type:        "error_handler",
			Description: "Detailed error handler with structured error responses",
			Content:     DetailedErrorHandlerTemplate,
			Parameters:  []TemplateParameter{},
		}, nil
	default:
		return nil, fmt.Errorf("error handler template not found: %s", name)
	}
}
