.PHONY: all build run clean test deps fmt lint help \
	docker-build docker-push docker-run docker-stop docker-compose-up docker-compose-down \
	k8s-deploy k8s-delete k8s-status k8s-logs

# ============================================
# Variables
# ============================================
APP_NAME := aiflowy-go
BUILD_DIR := build
MAIN_FILE := cmd/server/main.go
CONFIG_FILE := configs/config.yaml

# Version info
VERSION := $(shell cat VERSION 2>/dev/null || git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_SHA := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Docker settings
DOCKER_REGISTRY ?= docker.io
DOCKER_IMAGE := $(DOCKER_REGISTRY)/aiflowy/$(APP_NAME)
DOCKER_TAG := $(VERSION)

# Kubernetes settings
K8S_NAMESPACE := aiflowy
K8S_DEPLOY_DIR := deploy/k8s

# Build flags
LDFLAGS := -ldflags "-X 'main.version=$(VERSION)' -X 'main.buildTime=$(BUILD_TIME)' -X 'main.commitSHA=$(COMMIT_SHA)' -s -w"

# ============================================
# Development
# ============================================

# Default target
all: deps build

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Build the application
build:
	@echo "Building $(APP_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"
	@ls -lh $(BUILD_DIR)/$(APP_NAME)

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	go run $(MAIN_FILE) -config $(CONFIG_FILE)

# Run with hot reload (requires air)
dev:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	air

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR) 2>/dev/null || true
	go clean

# ============================================
# Testing
# ============================================

# Run all tests
test:
	@echo "Running tests..."
	go test -v -race -cover ./...

# Run tests with coverage report
test-cover:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run benchmark tests
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# ============================================
# Code Quality
# ============================================

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@which goimports > /dev/null && goimports -w . || true

# Lint code
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Generate code
generate:
	@echo "Generating code..."
	go generate ./...

# ============================================
# Docker
# ============================================

# Build Docker image
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg COMMIT_SHA=$(COMMIT_SHA) \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_IMAGE):latest \
		.
	@echo "Image size:"
	@docker images $(DOCKER_IMAGE):$(DOCKER_TAG) --format "{{.Size}}"

# Push Docker image
docker-push:
	@echo "Pushing Docker image..."
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_IMAGE):latest

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -d --name $(APP_NAME) \
		-p 8213:8213 \
		-e DB_HOST=host.docker.internal \
		-e DB_PASSWORD=aiflowy123 \
		-e REDIS_HOST=host.docker.internal \
		$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "Container started: $(APP_NAME)"
	@echo "Health check: curl http://localhost:8213/health"

# Stop Docker container
docker-stop:
	@echo "Stopping Docker container..."
	docker stop $(APP_NAME) 2>/dev/null || true
	docker rm $(APP_NAME) 2>/dev/null || true

# Docker Compose up
docker-compose-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d
	@echo "Services started. Waiting for health check..."
	@sleep 10
	@curl -s http://localhost:8213/health | head -1

# Docker Compose down
docker-compose-down:
	@echo "Stopping services..."
	docker-compose down

# Docker Compose logs
docker-compose-logs:
	docker-compose logs -f $(APP_NAME)

# ============================================
# Kubernetes
# ============================================

# Deploy to Kubernetes
k8s-deploy:
	@echo "Deploying to Kubernetes namespace $(K8S_NAMESPACE)..."
	kubectl apply -k $(K8S_DEPLOY_DIR)
	@echo "Waiting for deployment..."
	kubectl -n $(K8S_NAMESPACE) rollout status deployment/$(APP_NAME) --timeout=120s

# Delete from Kubernetes
k8s-delete:
	@echo "Deleting from Kubernetes..."
	kubectl delete -k $(K8S_DEPLOY_DIR) --ignore-not-found

# Show Kubernetes status
k8s-status:
	@echo "=== Pods ==="
	kubectl -n $(K8S_NAMESPACE) get pods -l app=$(APP_NAME)
	@echo ""
	@echo "=== Service ==="
	kubectl -n $(K8S_NAMESPACE) get svc $(APP_NAME)
	@echo ""
	@echo "=== HPA ==="
	kubectl -n $(K8S_NAMESPACE) get hpa $(APP_NAME) 2>/dev/null || echo "No HPA configured"

# Show Kubernetes logs
k8s-logs:
	kubectl -n $(K8S_NAMESPACE) logs -l app=$(APP_NAME) -f --tail=100

# Port forward for local testing
k8s-port-forward:
	@echo "Port forwarding to localhost:8213..."
	kubectl -n $(K8S_NAMESPACE) port-forward svc/$(APP_NAME) 8213:8213

# Restart deployment
k8s-restart:
	@echo "Restarting deployment..."
	kubectl -n $(K8S_NAMESPACE) rollout restart deployment/$(APP_NAME)

# Scale deployment
k8s-scale:
	@echo "Scaling deployment to $(REPLICAS) replicas..."
	kubectl -n $(K8S_NAMESPACE) scale deployment/$(APP_NAME) --replicas=$(REPLICAS)

# ============================================
# Version
# ============================================

# Show version info
version:
	@echo "Version:     $(VERSION)"
	@echo "Build Time:  $(BUILD_TIME)"
	@echo "Commit SHA:  $(COMMIT_SHA)"
	@echo "Go Version:  $(shell go version | cut -d' ' -f3)"

# ============================================
# Help
# ============================================

help:
	@echo ""
	@echo "AIFlowy Go Backend - Makefile"
	@echo ""
	@echo "Development:"
	@echo "  make deps          - Install dependencies"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make dev           - Run with hot reload"
	@echo "  make clean         - Clean build artifacts"
	@echo ""
	@echo "Testing:"
	@echo "  make test          - Run all tests"
	@echo "  make test-cover    - Run tests with coverage"
	@echo "  make bench         - Run benchmarks"
	@echo ""
	@echo "Code Quality:"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Lint code"
	@echo "  make generate      - Generate code"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-push   - Push image to registry"
	@echo "  make docker-run    - Run container"
	@echo "  make docker-stop   - Stop container"
	@echo "  make docker-compose-up   - Start with Docker Compose"
	@echo "  make docker-compose-down - Stop Docker Compose"
	@echo ""
	@echo "Kubernetes:"
	@echo "  make k8s-deploy    - Deploy to Kubernetes"
	@echo "  make k8s-delete    - Delete from Kubernetes"
	@echo "  make k8s-status    - Show deployment status"
	@echo "  make k8s-logs      - Show logs"
	@echo "  make k8s-restart   - Restart deployment"
	@echo "  make k8s-port-forward - Port forward to localhost"
	@echo ""
	@echo "Other:"
	@echo "  make version       - Show version info"
	@echo "  make help          - Show this help"
	@echo ""
