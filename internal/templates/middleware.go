package templates

import "fmt"

var ErrTemplateNotFound = fmt.Errorf("template not found")

// Middleware templates for go-zero services

const AuthMiddlewareTemplate = `package middleware

import (
	"net/http"
	"strings"
)

type {{.MiddlewareName}}Middleware struct {
	// Add your dependencies here
	{{if .SecretKey}}SecretKey string{{end}}
}

func New{{.MiddlewareName}}Middleware() *{{.MiddlewareName}}Middleware {
	return &{{.MiddlewareName}}Middleware{
		{{if .SecretKey}}SecretKey: "{{.SecretKey}}",{{end}}
	}
}

func (m *{{.MiddlewareName}}Middleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Check Bearer token format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// TODO: Validate token here
		// Example: Verify JWT, check database, etc.
		if token == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Pass to next handler
		next(w, r)
	}
}
`

const LoggingMiddlewareTemplate = `package middleware

import (
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type {{.MiddlewareName}}Middleware struct {
}

func New{{.MiddlewareName}}Middleware() *{{.MiddlewareName}}Middleware {
	return &{{.MiddlewareName}}Middleware{}
}

func (m *{{.MiddlewareName}}Middleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Log incoming request
		logx.Infof("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Create response writer wrapper to capture status code
		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call next handler
		next(wrapper, r)

		// Log response
		duration := time.Since(startTime)
		logx.Infof("Response: %s %s - Status: %d - Duration: %v",
			r.Method, r.URL.Path, wrapper.statusCode, duration)
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
`

const RateLimitingMiddlewareTemplate = `package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type {{.MiddlewareName}}Middleware struct {
	limiter *limit.PeriodLimit
}

func New{{.MiddlewareName}}Middleware(redisConf redis.RedisConf) *{{.MiddlewareName}}Middleware {
	// Configure rate limiter
	limiter := limit.NewPeriodLimit(
		{{.RequestsPerPeriod}},  // quota: requests per period
		{{.PeriodSeconds}},       // period in seconds
		redis.MustNewRedis(redisConf),
		"rate-limit",
	)

	return &{{.MiddlewareName}}Middleware{
		limiter: limiter,
	}
}

func (m *{{.MiddlewareName}}Middleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Use client IP as the key for rate limiting
		key := r.RemoteAddr

		// Check rate limit
		code, err := m.limiter.Take(key)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Check if rate limit exceeded
		switch code {
		case limit.OverQuota:
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		case limit.Allowed:
			next(w, r)
		case limit.HitQuota:
			next(w, r)
		default:
			next(w, r)
		}
	}
}
`

// GetMiddlewareTemplate returns a middleware template by name
func GetMiddlewareTemplate(name string) (*Template, error) {
	switch name {
	case "auth":
		return &Template{
			Name:        "auth",
			Type:        "middleware",
			Description: "JWT authentication middleware",
			Content:     AuthMiddlewareTemplate,
			Parameters: []TemplateParameter{
				{
					Name:        "MiddlewareName",
					Type:        "string",
					Description: "Name of the middleware",
					Required:    true,
					Default:     "Auth",
				},
				{
					Name:        "SecretKey",
					Type:        "string",
					Description: "Secret key for JWT validation",
					Required:    false,
					Default:     "your-secret-key",
				},
			},
		}, nil
	case "logging":
		return &Template{
			Name:        "logging",
			Type:        "middleware",
			Description: "Request/response logging middleware",
			Content:     LoggingMiddlewareTemplate,
			Parameters: []TemplateParameter{
				{
					Name:        "MiddlewareName",
					Type:        "string",
					Description: "Name of the middleware",
					Required:    true,
					Default:     "Logging",
				},
			},
		}, nil
	case "rate-limiting":
		return &Template{
			Name:        "rate-limiting",
			Type:        "middleware",
			Description: "Rate limiting middleware using Redis",
			Content:     RateLimitingMiddlewareTemplate,
			Parameters: []TemplateParameter{
				{
					Name:        "MiddlewareName",
					Type:        "string",
					Description: "Name of the middleware",
					Required:    true,
					Default:     "RateLimit",
				},
				{
					Name:        "RequestsPerPeriod",
					Type:        "int",
					Description: "Number of requests allowed per period",
					Required:    false,
					Default:     100,
				},
				{
					Name:        "PeriodSeconds",
					Type:        "int",
					Description: "Period duration in seconds",
					Required:    false,
					Default:     60,
				},
			},
		}, nil
	default:
		return nil, ErrTemplateNotFound
	}
}
