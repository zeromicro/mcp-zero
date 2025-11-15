package docs

import "strings"

// Concept represents a go-zero framework concept with explanation
type Concept struct {
	Name        string
	Category    string
	Description string
	Example     string
	RelatedDocs []string
}

// ConceptDatabase holds all go-zero framework concepts
var ConceptDatabase = map[string]Concept{
	"middleware": {
		Name:        "Middleware",
		Category:    "Core Concepts",
		Description: "Middleware in go-zero is a function that wraps HTTP handlers to add cross-cutting functionality. Middleware can handle authentication, logging, rate limiting, CORS, and more. go-zero supports both global middleware (applied to all routes) and route-specific middleware.\n\nMiddleware functions follow the pattern: func(next http.HandlerFunc) http.HandlerFunc\n\nMiddleware is registered in the .api file using the @server directive with the middleware keyword.",
		Example: `// Define middleware in middleware/auth.go
type AuthMiddleware struct {
    // dependencies
}

func NewAuthMiddleware() *AuthMiddleware {
    return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Pre-processing
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Call next handler
        next(w, r)

        // Post-processing (if needed)
    }
}

// Register in .api file:
@server(
    middleware: Auth
)
service user-api {
    @handler LoginHandler
    post /user/login (LoginReq) returns (LoginResp)
}`,
		RelatedDocs: []string{
			"https://go-zero.dev/docs/concepts/middleware",
			"https://go-zero.dev/docs/tutorials/http/middleware",
		},
	},
	"api-definition": {
		Name:        "API Definition",
		Category:    "Core Concepts",
		Description: "The .api file is go-zero's Domain Specific Language (DSL) for defining HTTP services. It describes routes, request/response types, middleware, and service configuration in a declarative format.\n\nKey sections:\n- info: Service metadata\n- type: Request/response type definitions\n- @server: Route group configuration (prefix, middleware, etc.)\n- service: Service definition with handlers and routes",
		Example: `syntax = "v1"

info(
    title: "User API"
    desc: "User management service"
    author: "developer"
    version: "v1"
)

type (
    LoginReq {
        Username string ` + "`json:\"username\"`" + `
        Password string ` + "`json:\"password\"`" + `
    }

    LoginResp {
        Token string ` + "`json:\"token\"`" + `
        UserId int64 ` + "`json:\"user_id\"`" + `
    }
)

@server(
    prefix: /api/v1
    middleware: Auth
    group: user
)
service user-api {
    @handler LoginHandler
    post /user/login (LoginReq) returns (LoginResp)

    @handler GetUserHandler
    get /user/:id returns (UserInfo)
}`,
		RelatedDocs: []string{
			"https://go-zero.dev/docs/tutorials/http/api-definition",
			"https://go-zero.dev/docs/reference/api-syntax",
		},
	},
	"service-context": {
		Name:        "Service Context",
		Category:    "Core Concepts",
		Description: "ServiceContext is the dependency injection container in go-zero. It holds all shared dependencies like database connections, Redis clients, RPC clients, and configuration. The context is created once at startup and passed to all handlers.\n\nBenefits:\n- Centralized dependency management\n- Easy to test (mock dependencies)\n- Singleton pattern for shared resources\n- Type-safe dependency access",
		Example: `// Define in internal/svc/servicecontext.go
type ServiceContext struct {
    Config config.Config
    UserModel model.UserModel
    RedisClient *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
    conn := sqlx.NewMysql(c.Mysql.DataSource)
    return &ServiceContext{
        Config: c,
        UserModel: model.NewUserModel(conn, c.CacheRedis),
        RedisClient: redis.MustNewRedis(c.Redis),
    }
}

// Use in handler
type LoginHandler struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func (l *LoginHandler) Login(req *types.LoginReq) (*types.LoginResp, error) {
    // Access dependencies through svcCtx
    user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
    if err != nil {
        return nil, err
    }

    return &types.LoginResp{
        Token: "...",
        UserId: user.Id,
    }, nil
}`,
		RelatedDocs: []string{
			"https://go-zero.dev/docs/concepts/service-context",
		},
	},
	"rpc": {
		Name:        "RPC Services",
		Category:    "Core Concepts",
		Description: "go-zero RPC services use gRPC for inter-service communication. Services are defined using Protocol Buffers (.proto files), and go-zero generates both server and client code.\n\nRPC features:\n- Automatic code generation from .proto files\n- Built-in service discovery (etcd, consul, kubernetes)\n- Load balancing\n- Circuit breaker\n- Interceptors for middleware functionality\n- Streaming support",
		Example: `// Define service in user.proto
syntax = "proto3";

package user;
option go_package = "./user";

message LoginRequest {
    string username = 1;
    string password = 2;
}

message LoginResponse {
    int64 user_id = 1;
    string token = 2;
}

service User {
    rpc Login(LoginRequest) returns (LoginResponse);
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
}

// Generate code:
// goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.

// Call RPC from API service:
type ServiceContext struct {
    UserRpc userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
    }
}`,
		RelatedDocs: []string{
			"https://go-zero.dev/docs/tutorials/rpc/server",
			"https://go-zero.dev/docs/tutorials/rpc/client",
		},
	},
	"model": {
		Name:        "Database Models",
		Category:    "Core Concepts",
		Description: "go-zero provides automatic model generation from database schema. Models include CRUD operations, cache management, and connection pooling.\n\nFeatures:\n- Auto-generated from table schema or SQL DDL\n- Redis cache integration\n- Prepared statements\n- Transaction support\n- Customizable query methods\n- Support for MySQL, PostgreSQL, MongoDB",
		Example: `// Generate model from database:
// goctl model mysql datasource -url="user:pass@tcp(127.0.0.1:3306)/database" -table="user" -dir="./model"

// Generated model usage:
type (
    userModel interface {
        Insert(ctx context.Context, data *User) (sql.Result, error)
        FindOne(ctx context.Context, id int64) (*User, error)
        Update(ctx context.Context, data *User) error
        Delete(ctx context.Context, id int64) error
    }

    defaultUserModel struct {
        sqlc.CachedConn
        table string
    }

    User struct {
        Id       int64  ` + "`db:\"id\"`" + `
        Username string ` + "`db:\"username\"`" + `
        Password string ` + "`db:\"password\"`" + `
        CreateTime time.Time ` + "`db:\"create_time\"`" + `
    }
)

// Use in service:
user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
if err != nil {
    return nil, err
}`,
		RelatedDocs: []string{
			"https://go-zero.dev/docs/tutorials/model/mysql",
			"https://go-zero.dev/docs/tutorials/model/cache",
		},
	},
	"configuration": {
		Name:        "Configuration",
		Category:    "Core Concepts",
		Description: "go-zero uses YAML configuration files to manage service settings. Configuration is automatically loaded and validated at startup. Supports environment variable expansion.\n\nCommon config sections:\n- Server: HTTP/RPC server settings (port, timeout, etc.)\n- Database: MySQL, PostgreSQL connection strings\n- Cache: Redis configuration\n- Log: Logging level and format\n- Telemetry: Prometheus, tracing settings",
		Example: `# etc/user-api.yaml
Name: user-api
Host: 0.0.0.0
Port: 8888

# Database configuration
Mysql:
  DataSource: user:password@tcp(localhost:3306)/users?charset=utf8mb4&parseTime=true

# Cache configuration
CacheRedis:
  - Host: localhost:6379
    Pass: ""
    Type: node

# Redis configuration
Redis:
  Host: localhost:6379
  Type: node
  Pass: ""

# Log configuration
Log:
  ServiceName: user-api
  Mode: console
  Level: info

# Telemetry
Prometheus:
  Host: 0.0.0.0
  Port: 9091
  Path: /metrics

# Load in code:
var c config.Config
conf.MustLoad(*configFile, &c)`,
		RelatedDocs: []string{
			"https://go-zero.dev/docs/tutorials/configuration/overview",
		},
	},
	"error-handling": {
		Name:        "Error Handling",
		Category:    "Best Practices",
		Description: "go-zero provides structured error handling with custom error codes and messages. Errors can be returned from handlers and automatically formatted into HTTP responses.\n\nError handling approaches:\n1. Return errors directly - converted to 500 status\n2. Use httpx.Error() for custom status codes\n3. Register custom error handler with httpx.SetErrorHandler()\n4. Define error types for different scenarios",
		Example: `// Custom error handler
func ErrorHandler(err error) (int, interface{}) {
    switch e := err.(type) {
    case *ValidationError:
        return http.StatusBadRequest, map[string]interface{}{
            "code": 400,
            "message": e.Message,
            "field": e.Field,
        }
    case *BusinessError:
        return e.StatusCode, map[string]interface{}{
            "code": e.Code,
            "message": e.Message,
        }
    default:
        return http.StatusInternalServerError, map[string]interface{}{
            "code": 500,
            "message": "Internal server error",
        }
    }
}

// Register in main.go
httpx.SetErrorHandler(ErrorHandler)

// Use in handler
func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
    if req.Username == "" {
        return nil, &ValidationError{
            Field: "username",
            Message: "username is required",
        }
    }
    // ... logic
}`,
		RelatedDocs: []string{
			"https://go-zero.dev/docs/tutorials/http/error-handling",
		},
	},
	"jwt": {
		Name:        "JWT Authentication",
		Category:    "Best Practices",
		Description: "go-zero has built-in JWT support for stateless authentication. JWT middleware can be configured in .api files to protect routes.\n\nFeatures:\n- Automatic token validation\n- Claims extraction\n- Token expiration handling\n- Custom secret key per service\n- Integration with middleware chain",
		Example: `// Configure JWT in .api file
@server(
    jwt: Auth
    prefix: /api/v1
)
service user-api {
    @handler GetUserInfoHandler
    get /user/info returns (UserInfoResp)
}

// Configure in config.yaml
Auth:
  AccessSecret: your-secret-key
  AccessExpire: 7200

// Generate token in login handler
func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
    // Verify credentials...

    now := time.Now().Unix()
    accessExpire := l.svcCtx.Config.Auth.AccessExpire

    token, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret,
        now, accessExpire, userId)
    if err != nil {
        return nil, err
    }

    return &types.LoginResp{
        Token: token,
        Expire: now + accessExpire,
    }, nil
}

// Extract claims in protected handler
userId := l.ctx.Value("userId").(json.Number).Int64()`,
		RelatedDocs: []string{
			"https://go-zero.dev/docs/tutorials/http/jwt",
		},
	},
	"cache": {
		Name:        "Cache",
		Category:    "Performance",
		Description: "go-zero provides automatic Redis caching for database models. Cache is managed transparently with cache-aside pattern.\n\nCache features:\n- Automatic cache key generation\n- Cache miss handling\n- Cache invalidation on updates/deletes\n- Configurable TTL\n- Support for cache tags\n- Distributed cache consistency",
		Example: `// Model with cache
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
    return &customUserModel{
        defaultUserModel: &defaultUserModel{
            CachedConn: sqlc.NewConn(conn, c),
            table:      "` + "`user`" + `",
        },
    }
}

// Cache is used automatically
user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
// First call: queries DB and caches result
// Subsequent calls: returns from cache

// Update invalidates cache
err = l.svcCtx.UserModel.Update(l.ctx, user)
// Cache entry is automatically deleted

// Configure cache in config.yaml
CacheRedis:
  - Host: localhost:6379
    Pass: ""
    Type: node`,
		RelatedDocs: []string{
			"https://go-zero.dev/docs/tutorials/model/cache",
		},
	},
	"validation": {
		Name:        "Request Validation",
		Category:    "Best Practices",
		Description: "go-zero supports automatic request validation using struct tags. Validation is performed before the handler is called.\n\nValidation tags:\n- optional: Field is optional\n- options: Enum validation\n- range: Numeric range validation\n- default: Default value if not provided\n- Custom validators can be added",
		Example: `// Define validation in types
type LoginReq {
    Username string ` + "`json:\"username\" validate:\"required,min=3,max=20\"`" + `
    Password string ` + "`json:\"password\" validate:\"required,min=6\"`" + `
    Email    string ` + "`json:\"email,optional\" validate:\"email\"`" + `
}

// Or use go-zero's built-in validation in .api file
type CreateUserReq {
    Username string ` + "`json:\"username\"`" + ` // required by default
    Age      int    ` + "`json:\"age,optional\"`" + ` // optional field
    Gender   string ` + "`json:\"gender,options=male|female\"`" + ` // enum validation
    Score    int    ` + "`json:\"score,range=[0:100]\"`" + ` // range validation
}

// Custom validation in handler
func (l *CreateUserLogic) CreateUser(req *types.CreateUserReq) error {
    if req.Age != 0 && (req.Age < 0 || req.Age > 150) {
        return errors.New("invalid age")
    }
    // ... logic
}`,
		RelatedDocs: []string{
			"https://go-zero.dev/docs/tutorials/http/parameter",
		},
	},
}

// SearchConcepts searches for concepts matching the query
func SearchConcepts(query string) []Concept {
	query = strings.ToLower(query)
	var results []Concept

	for _, concept := range ConceptDatabase {
		// Check if query matches name, category, or description
		if strings.Contains(strings.ToLower(concept.Name), query) ||
			strings.Contains(strings.ToLower(concept.Category), query) ||
			strings.Contains(strings.ToLower(concept.Description), query) {
			results = append(results, concept)
		}
	}

	return results
}

// GetConceptByName retrieves a specific concept by name
func GetConceptByName(name string) *Concept {
	name = strings.ToLower(name)
	for key, concept := range ConceptDatabase {
		if strings.ToLower(key) == name || strings.ToLower(concept.Name) == name {
			return &concept
		}
	}
	return nil
}
