package templates

import "fmt"

var ErrDeploymentTemplateNotFound = fmt.Errorf("deployment template not found")

const DockerfileTemplate = `FROM golang:1.19-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o {{.ServiceName}} .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/{{.ServiceName}} .
COPY --from=builder /build/etc /app/etc

EXPOSE {{.Port}}

CMD ["./{{.ServiceName}}"]
`

const KubernetesTemplate = `apiVersion: v1
kind: Service
metadata:
  name: {{.ServiceName}}
  labels:
    app: {{.ServiceName}}
spec:
  type: ClusterIP
  ports:
    - port: {{.Port}}
      targetPort: {{.Port}}
      protocol: TCP
      name: http
  selector:
    app: {{.ServiceName}}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ServiceName}}
  labels:
    app: {{.ServiceName}}
spec:
  replicas: {{.Replicas}}
  selector:
    matchLabels:
      app: {{.ServiceName}}
  template:
    metadata:
      labels:
        app: {{.ServiceName}}
    spec:
      containers:
      - name: {{.ServiceName}}
        image: {{.ImageName}}:{{.ImageTag}}
        ports:
        - containerPort: {{.Port}}
          name: http
        env:
        - name: SERVICE_NAME
          value: "{{.ServiceName}}"
        resources:
          limits:
            cpu: "{{.CPULimit}}"
            memory: "{{.MemoryLimit}}"
          requests:
            cpu: "{{.CPURequest}}"
            memory: "{{.MemoryRequest}}"
        livenessProbe:
          httpGet:
            path: /healthz
            port: {{.Port}}
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: {{.Port}}
          initialDelaySeconds: 5
          periodSeconds: 5
`

const SystemdTemplate = `[Unit]
Description={{.ServiceName}} service
After=network.target

[Service]
Type=simple
User={{.User}}
WorkingDirectory={{.WorkDir}}
ExecStart={{.WorkDir}}/{{.ServiceName}} -f {{.WorkDir}}/etc/{{.ServiceName}}.yaml
Restart=on-failure
RestartSec=5s

NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths={{.WorkDir}}/logs

LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
`

func GetDeploymentTemplate(name string) (*Template, error) {
	switch name {
	case "docker":
		return &Template{
			Name:        "docker",
			Type:        "deployment",
			Description: "Multi-stage Dockerfile for go-zero service",
			Content:     DockerfileTemplate,
			Parameters: []TemplateParameter{
				{
					Name:        "ServiceName",
					Type:        "string",
					Description: "Name of the service",
					Required:    true,
				},
				{
					Name:        "Port",
					Type:        "int",
					Description: "Service port",
					Required:    false,
					Default:     8888,
				},
			},
		}, nil
	case "kubernetes", "k8s":
		return &Template{
			Name:        "kubernetes",
			Type:        "deployment",
			Description: "Kubernetes deployment and service manifest",
			Content:     KubernetesTemplate,
			Parameters: []TemplateParameter{
				{Name: "ServiceName", Type: "string", Description: "Name of the service", Required: true},
				{Name: "Port", Type: "int", Description: "Service port", Required: false, Default: 8888},
				{Name: "Replicas", Type: "int", Description: "Number of replicas", Required: false, Default: 3},
				{Name: "ImageName", Type: "string", Description: "Docker image name", Required: false, Default: "my-service"},
				{Name: "ImageTag", Type: "string", Description: "Docker image tag", Required: false, Default: "latest"},
				{Name: "CPULimit", Type: "string", Description: "CPU limit", Required: false, Default: "1000m"},
				{Name: "MemoryLimit", Type: "string", Description: "Memory limit", Required: false, Default: "512Mi"},
				{Name: "CPURequest", Type: "string", Description: "CPU request", Required: false, Default: "100m"},
				{Name: "MemoryRequest", Type: "string", Description: "Memory request", Required: false, Default: "128Mi"},
			},
		}, nil
	case "systemd":
		return &Template{
			Name:        "systemd",
			Type:        "deployment",
			Description: "Systemd service unit file",
			Content:     SystemdTemplate,
			Parameters: []TemplateParameter{
				{Name: "ServiceName", Type: "string", Description: "Name of the service", Required: true},
				{Name: "User", Type: "string", Description: "User to run service as", Required: false, Default: "app"},
				{Name: "WorkDir", Type: "string", Description: "Working directory", Required: false, Default: "/opt/app"},
			},
		}, nil
	default:
		return nil, ErrDeploymentTemplateNotFound
	}
}
