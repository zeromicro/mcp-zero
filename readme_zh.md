# Go-Zero MCP 工具

一个基于 Model Context Protocol (MCP) 的工具，帮助开发者快速构建和生成 go-zero 项目。

## 快速开始

**第一次使用 mcp-zero？** 查看我们的[快速入门指南](QUICKSTART.md)获取详细的分步教程！

快速入门涵盖：
- 安装和配置
- 创建你的第一个 API 服务
- 常见用例和工作流程
- 与 Claude Desktop 集成

## 功能特性

### 核心服务生成
- **创建 API 服务**：生成新的 REST API 服务，支持自定义端口和代码风格
- **创建 RPC 服务**：从 protobuf 定义生成 gRPC 服务
- **生成 API 代码**：将 API 规范文件转换为 Go 代码
- **生成数据模型**：从多种数据源创建数据库模型（MySQL、PostgreSQL、MongoDB、DDL）
- **创建 API 规范**：生成示例 API 规范文件

### 高级功能
- **项目分析**：分析现有 go-zero 项目以了解结构和依赖关系
- **配置管理**：生成具有正确结构验证的配置文件
- **生成模板**：创建中间件、错误处理器和部署模板
- **文档查询**：访问 go-zero 概念和从其他框架迁移的指南
- **输入验证**：对 API 规范、protobuf 定义和配置进行全面验证

## 前置要求

1. **Go** (1.19 或更高版本)
2. **go-zero CLI (goctl)**：通过 `go install github.com/zeromicro/go-zero/tools/goctl@latest` 安装
3. **Claude Desktop**（或其他兼容 MCP 的客户端）

详细安装说明请参阅[快速入门指南](QUICKSTART.md)。

## 安装

1. 创建新目录用于 MCP 工具：
```bash
mkdir go-zero-mcp && cd go-zero-mcp
```

2. 初始化 Go 模块：
```bash
go mod init go-zero-mcp
```

3. 安装依赖：
```bash
go get github.com/mark3labs/mcp-go/mcp
go get github.com/mark3labs/mcp-go/server
```

4. 将主工具代码保存为 `main.go`

5. 构建工具：
```bash
go build -o go-zero-mcp main.go
```

## 配置 Claude Desktop

将此配置添加到你的 Claude Desktop MCP 设置：

### macOS/Linux
编辑 `~/.config/claude/mcp_settings.json`：

```json
{
  "mcpServers": {
    "go-zero-mcp": {
      "command": "/path/to/your/go-zero-mcp",
      "args": [],
      "env": {
        "PATH": "/usr/local/bin:/usr/bin:/bin"
      }
    }
  }
}
```

### Windows
编辑 `%APPDATA%/Claude/mcp_settings.json`：

```json
{
  "mcpServers": {
    "go-zero-mcp": {
      "command": "C:\\path\\to\\your\\go-zero-mcp.exe",
      "args": [],
      "env": {
        "PATH": "C:\\Go\\bin;C:\\Program Files\\Go\\bin"
      }
    }
  }
}
```

## 可用工具

### 1. create_api_service
创建新的 go-zero API 服务。

**参数：**
- `service_name`（必需）：API 服务名称
- `port`（可选）：端口号（默认：8888）
- `style`（可选）：代码风格 - "go_zero" 或 "gozero"（默认："go_zero"）
- `output_dir`（可选）：输出目录（默认：当前目录）

### 2. create_rpc_service
从 protobuf 定义创建新的 go-zero RPC 服务。

**参数：**
- `service_name`（必需）：RPC 服务名称
- `proto_content`（必需）：Protobuf 定义内容
- `output_dir`（可选）：输出目录（默认：当前目录）

### 3. generate_api_from_spec
从 API 规范文件生成 go-zero API 代码。

**参数：**
- `api_file`（必需）：.api 规范文件路径
- `output_dir`（可选）：输出目录（默认：当前目录）
- `style`（可选）：代码风格 - "go_zero" 或 "gozero"（默认："go_zero"）

### 4. generate_model
从数据库模式或 DDL 生成 go-zero 模型代码。

**参数：**
- `source_type`（必需）：源类型 - "mysql"、"postgresql"、"mongo" 或 "ddl"
- `source`（必需）：数据库连接字符串或 DDL 文件路径
- `table`（可选）：特定表名（用于数据库源）
- `output_dir`（可选）：输出目录（默认："./model"）

### 5. create_api_spec
创建示例 API 规范文件。

**参数：**
- `service_name`（必需）：API 服务名称
- `endpoints`（必需）：端点对象数组，包含方法、路径和处理器
- `output_file`（可选）：输出文件路径（默认：service_name.api）

### 6. analyze_project
分析现有 go-zero 项目结构和依赖关系。

**参数：**
- `project_dir`（必需）：项目目录路径
- `analysis_type`（可选）：分析类型 - "api"、"rpc"、"model" 或 "full"（默认："full"）

### 7. generate_config
为 go-zero 服务生成配置文件。

**参数：**
- `service_name`（必需）：服务名称
- `service_type`（必需）：服务类型 - "api" 或 "rpc"
- `config_type`（可选）：配置类型 - "dev"、"test" 或 "prod"（默认："dev"）
- `output_file`（可选）：输出文件路径（默认：etc/{service_name}.yaml）

### 8. generate_template
为 go-zero 服务生成常用代码模板。

**参数：**
- `template_type`（必需）：模板类型 - "middleware"、"error_handler"、"dockerfile"、"docker_compose" 或 "kubernetes"
- `service_name`（必需）：服务名称
- `output_path`（可选）：输出文件路径（根据模板类型使用默认值）

### 9. query_docs
查询 go-zero 文档和迁移指南。

**参数：**
- `query`（必需）：关于 go-zero 概念或迁移的自然语言查询
- `doc_type`（可选）：文档类型 - "concept"、"migration" 或 "both"（默认："both"）

### 10. validate_input
验证 API 规范、protobuf 定义或配置文件。

**参数：**
- `input_type`（必需）：输入类型 - "api_spec"、"proto" 或 "config"
- `content`（必需）：要验证的内容
- `strict`（可选）：启用严格验证模式（默认：false）

## 使用示例

### 创建新的 API 服务
```
请创建一个名为 "user-service" 的 go-zero API 服务，端口为 8080
```

### 创建 RPC 服务
```
创建一个名为 "auth-service" 的 go-zero RPC 服务，使用以下 protobuf 定义：

syntax = "proto3";

package auth;

option go_package = "./auth";

service AuthService {
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc Logout(LogoutRequest) returns (LogoutResponse);
}

message LoginRequest {
  string username = 1;
  string password = 2;
}

message LoginResponse {
  string token = 1;
  int64 expires_at = 2;
}

message LogoutRequest {
  string token = 1;
}

message LogoutResponse {
  bool success = 1;
}
```

### 从数据库生成模型
```
使用连接字符串 "user:password@tcp(localhost:3306)/mydb" 从我的 MySQL 数据库生成 go-zero 模型
```

### 创建 API 规范
```
为 "blog-service" 创建一个 API 规范，包含以下端点：
- GET /api/posts（处理器：GetPostsHandler）
- POST /api/posts（处理器：CreatePostHandler）
- GET /api/posts/:id（处理器：GetPostHandler）
- PUT /api/posts/:id（处理器：UpdatePostHandler）
- DELETE /api/posts/:id（处理器：DeletePostHandler）
```

### 分析项目
```
分析我在 /path/to/myproject 的 go-zero 项目以了解其结构
```

### 生成配置
```
为我的 "order-service" API 服务生成生产环境配置文件
```

### 生成模板
```
为我的 "auth-service" 生成中间件模板
```

### 查询文档
```
如何在 go-zero 中实现 JWT 认证？
```

```
如何从 Express.js 迁移到 go-zero？
```

### 验证输入
```
在严格模式下验证位于 /path/to/service.api 的 API 规范文件
```

## 项目结构

构建后，你的 MCP 服务器将具有以下结构：

```
mcp-zero/
├── main.go                    # 入口点和工具注册
├── tools/                     # 工具实现
│   ├── create_api_service.go
│   ├── create_rpc_service.go
│   ├── generate_api.go
│   ├── generate_model.go
│   ├── create_api_spec.go
│   ├── analyze_project.go
│   ├── generate_config.go
│   ├── generate_template.go
│   ├── query_docs.go
│   └── validate_input.go
├── internal/                  # 内部包
│   ├── analyzer/             # 项目分析
│   ├── validation/           # 输入验证
│   ├── security/             # 凭证处理
│   ├── templates/            # 代码模板
│   ├── docs/                 # 文档数据库
│   ├── logging/              # 结构化日志
│   └── metrics/              # 性能指标
└── tests/                     # 测试套件
    ├── integration/
    └── unit/
```

## 架构

MCP 服务器基于以下技术构建：

- **MCP SDK**：使用 github.com/modelcontextprotocol/go-sdk 实现协议
- **传输**：基于 stdio 与 Claude Desktop 通信
- **代码生成**：利用 go-zero 的 goctl CLI 工具生成生产就绪代码
- **验证**：全面的输入验证以确保安全性和正确性
- **安全**：使用环境变量替换进行安全的凭证处理
- **可观测性**：内置日志和指标用于监控工具性能

## 最佳实践

1. **服务命名**：使用小写字母和连字符（例如："user-service"、"auth-api"）
2. **端口配置**：为每个服务选择唯一端口（建议 8080-8090 范围）
3. **代码风格**：坚持使用 "go_zero" 风格以与官方约定保持一致
4. **配置**：使用特定环境的配置（dev、test、prod）
5. **文档**：定期查询文档以与 go-zero 最佳实践保持一致
6. **验证**：在生成之前始终验证输入以尽早发现错误

## 故障排除

### 常见问题

1. **找不到 goctl 命令**：确保已安装 goctl 并在 PATH 中
2. **权限被拒绝**：确保 MCP 工具可执行文件具有适当的权限
3. **数据库连接错误**：验证连接字符串和数据库可访问性

### 调试模式

要启用调试日志，设置环境变量：
```bash
export MCP_DEBUG=1
```

## 演示项目

我们提供了一个完整的书店管理系统演示，展示了所有 mcp-zero 工具的功能：

### 演示特性

- ✅ **REST API 生成** - 完整的书籍和订单 CRUD 操作（10 个端点）
- ✅ **数据库模型生成** - 从 SQL 模式自动生成的模型
- ✅ **RPC 服务创建** - 订单处理微服务
- ✅ **API 规范** - 类型安全的 API 定义
- ✅ **项目分析** - 理解 go-zero 项目结构

### 运行演示

```bash
# 运行自动化测试套件
python3 demo/test_bookstore.py

# 构建服务
cd demo/bookstore
go build -o bookstore-server

# 查看生成的代码
ls -R internal/
```

### 演示结果

```
✅ 所有测试通过！（5/5）
✅ 二进制构建：23MB
✅ 节省时间：99.2%
```

**代码放大倍数：6.5x**（202 行输入 → 1,300+ 行输出）

详细信息请查看 `demo/README.md`。

## 贡献

欢迎扩展此工具，添加更多 go-zero 功能，例如：
- Dockerfile 生成
- Kubernetes 清单生成
- Docker Compose 文件创建
- API 文档生成
- 测试模板创建

## 文档

- [快速入门指南](QUICKSTART.md) - 新手入门教程
- [贡献指南](CONTRIBUTING.md) - 开发者贡献指南
- [错误恢复](ERROR_RECOVERY.md) - 故障排除和错误恢复
- [安全审计](SECURITY_AUDIT.md) - 安全最佳实践
- [发布流程](RELEASE.md) - 版本发布说明

## 测试

项目包含全面的测试套件：

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./tests/unit/...
go test ./tests/integration/...

# 运行带覆盖率的测试
go test -cover ./...
```

## CI/CD

项目配置了 GitHub Actions 工作流：

- **CI 流水线**：在每次提交时运行测试、代码检查和构建
- **发布自动化**：自动构建多平台二进制文件并创建 GitHub 发布

## 性能指标

- **启动时间**：< 100ms
- **内存占用**：< 50MB
- **代码生成速度**：1-5 秒（取决于项目大小）
- **验证速度**：< 100ms

## 许可证

本工具按"原样"提供，用于教育和开发目的。

## 致谢

- [go-zero](https://github.com/zeromicro/go-zero) - 优秀的微服务框架
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) - Model Context Protocol Go 实现
- [Claude](https://claude.ai) - AI 助手平台

## 联系方式

如有问题或建议，请访问 [GitHub Issues](https://github.com/zeromicro/mcp-zero/issues)。

---

**开始使用 mcp-zero，让 AI 助力你的 go-zero 开发！** 🚀
