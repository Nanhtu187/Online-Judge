# Root Makefile for Online Judge project

.PHONY: all proto backend-build backend-run docker-build docker-load frontend-install frontend-build frontend-dev clean help

SERVICES := iam server runner

# Default target
all: frontend-install proto backend-build frontend-build

# Proto generation
# NOTE: Requires frontend-install to be run first for buf and plugins
proto:
	@echo "Generating proto files..."
	@cd src/packages/proto && \
	 buf dep update && \
	 buf generate

# Backend targets
backend-build:
	@echo "Building backend services..."
	bazel build //src/apps/backend/...

# Container targets
docker-build:
	@echo "Building OCI images via Bazel for linux/arm64..."
	bazel build --platforms=@rules_go//go/toolchain:linux_arm64 //src/apps/backend/iam:tarball
	bazel build --platforms=@rules_go//go/toolchain:linux_arm64 //src/apps/backend/server:tarball
	bazel build --platforms=@rules_go//go/toolchain:linux_arm64 //src/apps/backend/runner:tarball
	@echo "Building frontend image via Docker..."
	docker build -t online-judge/frontend:latest -f src/apps/frontend/Dockerfile .

docker-load: docker-build
	@echo "Loading OCI images into Docker daemon..."
	docker load -i bazel-bin/src/apps/backend/iam/tarball/tarball.tar
	docker load -i bazel-bin/src/apps/backend/server/tarball/tarball.tar
	docker load -i bazel-bin/src/apps/backend/runner/tarball/tarball.tar

backend-run-%:
	bazel run //src/apps/backend/$*

# Run all backend services in background
backend-run:
	@for service in $(SERVICES); do \
		echo "Starting $$service..."; \
		bazel run //src/apps/backend/$$service & \
	done

# Frontend targets
frontend-install:
	@echo "Installing frontend dependencies..."
	@cd src/apps/frontend && yarn install

frontend-build:
	@echo "Building frontend..."
	@cd src/apps/frontend && yarn build

frontend-dev:
	@echo "Starting frontend dev server..."
	@cd src/apps/frontend && yarn dev

# Convenience targets
dev:
	@echo "Starting development environment..."
	@$(MAKE) -j 4 backend-run-iam backend-run-server backend-run-runner frontend-dev

clean:
	@echo "Cleaning up..."
	bazel clean
	rm -rf src/apps/frontend/dist src/apps/frontend/.turbo

help:
	@echo "Available targets:"
	@echo "  proto               Generate proto files"
	@echo "  backend-build       Build all Go services via Bazel"
	@echo "  docker-build        Build OCI container images via Bazel"
	@echo "  docker-load         Build and load OCI images into local Docker"
	@echo "  backend-run-[name]  Run a specific service via Bazel"
	@echo "  backend-run         Run all backend services"
	@echo "  frontend-install    Install frontend dependencies"
	@echo "  frontend-build      Build the frontend"
	@echo "  frontend-dev        Run the frontend dev server"
	@echo "  dev                 Run everything for development"
	@echo "  clean               Clean build artifacts"
	@echo "  all                 Full build pipeline"
