package deployment

import (
	"fmt"
	"os"
	"text/template"
)

type DeploymentHelper struct {
	APIName    string
	APIVersion string
	Port       int
}

func NewDeploymentHelper(apiName string, apiVersion string, port int) *DeploymentHelper {
	return &DeploymentHelper{
		APIName:    apiName,
		APIVersion: apiVersion,
		Port:       port,
	}
}

func (dh *DeploymentHelper) GenerateDockerfile() error {
	dockerfileTemplate := `FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /soft-crusher-api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /soft-crusher-api ./

EXPOSE {{.Port}}

CMD ["./soft-crusher-api"]
`

	return dh.generateFile("Dockerfile", dockerfileTemplate)
}

func (dh *DeploymentHelper) GenerateKubernetesManifests() error {
	k8sTemplate := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.APIName}}
spec:
  replicas: 3
  selector:
    matchLabels:
      app: {{.APIName}}
  template:
    metadata:
      labels:
        app: {{.APIName}}
    spec:
      containers:
      - name: {{.APIName}}
        image: {{.APIName}}:{{.APIVersion}}
        ports:
        - containerPort: {{.Port}}
---
apiVersion: v1
kind: Service
metadata:
  name: {{.APIName}}-service
spec:
  selector:
    app: {{.APIName}}
  ports:
    - protocol: TCP
      port: 80
      targetPort: {{.Port}}
  type: LoadBalancer
`

	return dh.generateFile("kubernetes-manifests.yaml", k8sTemplate)
}

func (dh *DeploymentHelper) GenerateDockerComposeFile() error {
	dockerComposeTemplate := `version: '3'
services:
  {{.APIName}}:
    build: .
    ports:
      - "{{.Port}}:{{.Port}}"
    environment:
      - GIN_MODE=release
`

	return dh.generateFile("docker-compose.yaml", dockerComposeTemplate)
}

func (dh *DeploymentHelper) generateFile(filename string, templateContent string) error {
	tmpl, err := template.New(filename).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("error parsing template for %s: %v", filename, err)
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating %s: %v", filename, err)
	}
	defer f.Close()

	err = tmpl.Execute(f, dh)
	if err != nil {
		return fmt.Errorf("error executing template for %s: %v", filename, err)
	}

	fmt.Printf("%s generated successfully\n", filename)
	return nil
}