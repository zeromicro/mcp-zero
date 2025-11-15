package docs

import "strings"

// MigrationGuide represents a migration guide from other frameworks to go-zero
type MigrationGuide struct {
	FromFramework  string
	ToGoZero       string
	Difficulty     string
	KeyDifferences string
	Example        string
	Steps          []string
}

// MigrationDatabase holds migration guides
var MigrationDatabase = map[string]MigrationGuide{
	"gin-to-gozero": {
		FromFramework: "Gin",
		ToGoZero:      "go-zero API service",
		Difficulty:    "Easy",
		KeyDifferences: `Main differences:
1. Route definition: Gin uses code, go-zero uses .api DSL
2. Handler structure: go-zero separates handler and logic layers
3. Dependency injection: go-zero uses ServiceContext
4. Code generation: go-zero auto-generates boilerplate`,
		Example: `// Gin handler
func LoginHandler(c *gin.Context) {
    var req LoginReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // business logic here
    c.JSON(200, gin.H{"token": "xxx"})
}

// go-zero equivalent
// 1. Define in .api file:
type LoginReq {
    Username string ` + "`json:\"username\"`" + `
    Password string ` + "`json:\"password\"`" + `
}

type LoginResp {
    Token string ` + "`json:\"token\"`" + `
}

service user-api {
    @handler LoginHandler
    post /user/login (LoginReq) returns (LoginResp)
}

// 2. Implement in logic file:
func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
    // business logic here
    return &types.LoginResp{
        Token: "xxx",
    }, nil
}`,
		Steps: []string{
			"1. Create .api file to define routes and types (replaces Gin router setup)",
			"2. Run 'goctl api go -api user.api -dir .' to generate code",
			"3. Move business logic from Gin handlers to generated logic files",
			"4. Replace Gin context usage with go-zero types and error handling",
			"5. Move middleware to go-zero middleware pattern",
			"6. Update configuration to use YAML instead of Gin config",
		},
	},
	"echo-to-gozero": {
		FromFramework: "Echo",
		ToGoZero:      "go-zero API service",
		Difficulty:    "Easy",
		KeyDifferences: `Main differences:
1. Echo uses echo.Context, go-zero uses generated types
2. Route definition: Echo uses code, go-zero uses .api DSL
3. Middleware: Similar pattern but different signature
4. Error handling: go-zero has centralized error handler`,
		Example: `// Echo handler
func login(c echo.Context) error {
    req := new(LoginRequest)
    if err := c.Bind(req); err != nil {
        return c.JSON(400, map[string]string{"error": err.Error()})
    }

    return c.JSON(200, map[string]string{"token": "xxx"})
}

// go-zero logic
func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
    // business logic
    return &types.LoginResp{Token: "xxx"}, nil
}`,
		Steps: []string{
			"1. Map Echo routes to .api file definitions",
			"2. Generate go-zero code structure",
			"3. Convert Echo handlers to go-zero logic functions",
			"4. Adapt Echo middleware to go-zero middleware pattern",
			"5. Replace echo.Context usage with request/response types",
			"6. Set up centralized error handler",
		},
	},
	"grpc-to-gozero-rpc": {
		FromFramework: "gRPC (vanilla)",
		ToGoZero:      "go-zero RPC service",
		Difficulty:    "Easy",
		KeyDifferences: `Main differences:
1. go-zero adds automatic service discovery
2. Built-in load balancing and circuit breaker
3. Simplified client creation
4. Integrated with go-zero ecosystem (logging, tracing, metrics)`,
		Example: `// Vanilla gRPC server
type server struct {
    pb.UnimplementedUserServer
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    // logic
    return &pb.GetUserResponse{UserId: 1}, nil
}

// go-zero RPC (proto file is same, but generation and usage differ)
// Generate with: goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.

// Server automatically gets:
// - Service discovery (etcd/consul/k8s)
// - Health checks
// - Metrics
// - Interceptors

// Client usage is simplified:
userRpc := userclient.NewUser(zrpc.MustNewClient(c.UserRpc))
resp, err := userRpc.GetUser(ctx, &user.GetUserRequest{UserId: 1})`,
		Steps: []string{
			"1. Keep existing .proto files (protobuf definitions stay the same)",
			"2. Use goctl to generate go-zero RPC code",
			"3. Migrate server implementation to generated logic files",
			"4. Update client code to use go-zero RPC client",
			"5. Configure service discovery (etcd, consul, or kubernetes)",
			"6. Add configuration files (YAML) for server and clients",
			"7. Enable metrics and tracing if needed",
		},
	},
	"springboot-to-gozero": {
		FromFramework: "Spring Boot",
		ToGoZero:      "go-zero microservices",
		Difficulty:    "Moderate",
		KeyDifferences: `Main differences:
1. Language: Java -> Go
2. Annotations -> .api/.proto files
3. Spring DI -> ServiceContext
4. JPA -> go-zero models
5. Compiled binary vs JVM`,
		Example: `// Spring Boot controller
@RestController
@RequestMapping("/api/user")
public class UserController {
    @Autowired
    private UserService userService;

    @PostMapping("/login")
    public ResponseEntity<LoginResp> login(@RequestBody LoginReq req) {
        return ResponseEntity.ok(userService.login(req));
    }
}

// go-zero equivalent
// .api file:
service user-api {
    @handler LoginHandler
    post /api/user/login (LoginReq) returns (LoginResp)
}

// Logic:
type LoginLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
    // Access dependencies through svcCtx (like Spring DI)
    return l.svcCtx.UserService.Login(req)
}`,
		Steps: []string{
			"1. Map REST controllers to .api files",
			"2. Convert JPA entities to database models (use goctl model)",
			"3. Replace @Autowired with ServiceContext dependency injection",
			"4. Migrate service layer to go-zero logic layer",
			"5. Convert Java exceptions to Go error handling",
			"6. Replace Spring configuration with YAML config files",
			"7. Update deployment from JAR to Go binary",
			"8. Migrate Spring Security to go-zero JWT/middleware",
		},
	},
	"nodejs-express-to-gozero": {
		FromFramework: "Node.js Express",
		ToGoZero:      "go-zero API service",
		Difficulty:    "Moderate",
		KeyDifferences: `Main differences:
1. Language: JavaScript/TypeScript -> Go
2. Async/await -> goroutines and channels
3. Express middleware -> go-zero middleware
4. Route definition: Code -> .api DSL
5. Package.json -> go.mod`,
		Example: `// Express route
app.post('/api/user/login', async (req, res) => {
    try {
        const { username, password } = req.body;
        const result = await userService.login(username, password);
        res.json(result);
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

// go-zero
// .api file:
type LoginReq {
    Username string ` + "`json:\"username\"`" + `
    Password string ` + "`json:\"password\"`" + `
}

service user-api {
    @handler LoginHandler
    post /api/user/login (LoginReq) returns (LoginResp)
}

// Logic:
func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
    result, err := l.svcCtx.UserService.Login(req.Username, req.Password)
    if err != nil {
        return nil, err
    }
    return result, nil
}`,
		Steps: []string{
			"1. Map Express routes to .api file definitions",
			"2. Convert async/await patterns to goroutines where needed",
			"3. Migrate Express middleware to go-zero middleware",
			"4. Replace Mongoose/Sequelize with go-zero models",
			"5. Convert promises to Go error handling",
			"6. Update package dependencies (npm -> go modules)",
			"7. Rewrite JavaScript/TypeScript logic in Go",
			"8. Update deployment from Node.js to Go binary",
		},
	},
}

// SearchMigrationGuides searches for migration guides matching the query
func SearchMigrationGuides(query string) []MigrationGuide {
	query = strings.ToLower(query)
	var results []MigrationGuide

	for _, guide := range MigrationDatabase {
		// Check if query matches framework names or description
		if strings.Contains(strings.ToLower(guide.FromFramework), query) ||
			strings.Contains(strings.ToLower(guide.ToGoZero), query) ||
			strings.Contains(strings.ToLower(guide.KeyDifferences), query) {
			results = append(results, guide)
		}
	}

	return results
}

// GetMigrationGuide retrieves a specific migration guide
func GetMigrationGuide(fromFramework string) *MigrationGuide {
	key := strings.ToLower(fromFramework) + "-to-gozero"
	if guide, exists := MigrationDatabase[key]; exists {
		return &guide
	}

	// Try partial match
	fromFramework = strings.ToLower(fromFramework)
	for _, guide := range MigrationDatabase {
		if strings.Contains(strings.ToLower(guide.FromFramework), fromFramework) {
			return &guide
		}
	}

	return nil
}
