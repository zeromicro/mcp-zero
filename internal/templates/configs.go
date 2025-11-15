package templates

// ConfigTemplate represents a configuration template
type ConfigTemplate struct {
	Type        string // "api" or "rpc"
	Environment string // "development", "production", "test"
	Content     string
}

// APIConfigDevelopment is the development config template for API services
const APIConfigDevelopment = `Name: {{.ServiceName}}
Host: 0.0.0.0
Port: {{.Port}}

Log:
  Mode: console
  Level: info
  Encoding: plain

Timeout: 3000
MaxConns: 10000
`

// APIConfigProduction is the production config template for API services
const APIConfigProduction = `Name: {{.ServiceName}}
Host: 0.0.0.0
Port: {{.Port}}
Mode: prod

Log:
  Mode: file
  Path: logs
  Level: error
  Encoding: json
  KeepDays: 7
  Compress: true

Timeout: 5000
MaxConns: 50000

# Prometheus metrics
Prometheus:
  Host: 0.0.0.0
  Port: {{.MetricsPort}}
  Path: /metrics

# Service registration
ServiceConf:
  Name: {{.ServiceName}}
  # Add your service discovery config here
`

// RPCConfigDevelopment is the development config template for RPC services
const RPCConfigDevelopment = `Name: {{.ServiceName}}
ListenOn: 0.0.0.0:{{.Port}}

Log:
  Mode: console
  Level: info
  Encoding: plain

Timeout: 3000
`

// RPCConfigProduction is the production config template for RPC services
const RPCConfigProduction = `Name: {{.ServiceName}}
ListenOn: 0.0.0.0:{{.Port}}
Mode: prod

Log:
  Mode: file
  Path: logs
  Level: error
  Encoding: json
  KeepDays: 7
  Compress: true

Timeout: 5000

# Service discovery with etcd
Etcd:
  Hosts:
    - etcd-host-1:2379
    - etcd-host-2:2379
    - etcd-host-3:2379
  Key: {{.ServiceName}}

# Prometheus metrics
Prometheus:
  Host: 0.0.0.0
  Port: {{.MetricsPort}}
  Path: /metrics

# Redis cache (optional)
# Cache:
#   - Host: redis-host:6379
#     Type: node
`

// APIConfigTest is the test config template for API services
const APIConfigTest = `Name: {{.ServiceName}}-test
Host: 127.0.0.1
Port: {{.Port}}

Log:
  Mode: console
  Level: debug
  Encoding: plain

Timeout: 1000
MaxConns: 100
`

// RPCConfigTest is the test config template for RPC services
const RPCConfigTest = `Name: {{.ServiceName}}-test
ListenOn: 127.0.0.1:{{.Port}}

Log:
  Mode: console
  Level: debug
  Encoding: plain

Timeout: 1000
`

// GetConfigTemplate returns the appropriate config template
func GetConfigTemplate(serviceType, environment string) string {
	switch serviceType {
	case "api":
		switch environment {
		case "development", "dev":
			return APIConfigDevelopment
		case "production", "prod":
			return APIConfigProduction
		case "test":
			return APIConfigTest
		default:
			return APIConfigDevelopment
		}
	case "rpc":
		switch environment {
		case "development", "dev":
			return RPCConfigDevelopment
		case "production", "prod":
			return RPCConfigProduction
		case "test":
			return RPCConfigTest
		default:
			return RPCConfigDevelopment
		}
	}
	return ""
}
