# Makefile for Family Service
# This Makefile provides targets for building, testing, and deploying the Family Service application.

#################################################
# VARIABLES
#################################################

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GORUN=$(GOCMD) run
GOLINT=golangci-lint
GOVULNCHECK=$(GOCMD) run golang.org/x/vuln/cmd/govulncheck
GO_VERSION=1.24.4

# Application parameters
BINARY_NAME=family-service
MAIN_PATH=./cmd/server/graphql
BINARY_OUTPUT=./bin/$(BINARY_NAME)
DOCKER_IMAGE=family-service
DOCKER_TAG=latest

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# Cross-compilation parameters
PLATFORMS=linux darwin windows
ARCHITECTURES=amd64 arm64

# GraphQL parameters
GQLGEN=go run github.com/99designs/gqlgen
GQLGEN_CONFIG=interface/adapters/graphql/gqlgen.yml

#################################################
# DEFAULT AND HELP TARGETS
#################################################

# Default target
.PHONY: all
all: help

# Show version information
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"

# Show help
.PHONY: help
help:
	@echo "Family Service Makefile"
	@echo "Usage:"
	@echo ""
	@echo "#################################################"
	@echo "# DEFAULT AND HELP TARGETS"
	@echo "#################################################"
	@echo "  make help              - Show this help message"
	@echo "  make version           - Show version information"
	@echo ""
	@echo "#################################################"
	@echo "# BUILD AND DEVELOPMENT TARGETS"
	@echo "#################################################"
	@echo "  make build             - Build the application"
	@echo "  make build-all         - Build the application for all platforms and architectures"
	@echo "  make dev               - Run the application with hot reloading"
	@echo "  make graphql-gen       - Generate GraphQL code"
	@echo "  make init              - Initialize development environment"
	@echo "  make run               - Run the application locally"
	@echo ""
	@echo "#################################################"
	@echo "# TESTING TARGETS"
	@echo "#################################################"
	@echo "  make test              - Run all tests"
	@echo "  make test-all          - Run all tests with coverage and generate a combined report"
	@echo "  make test-bench        - Run benchmarks"
	@echo "  make test-integration  - Run integration tests only"
	@echo "  make test-package      - Run tests for a specific package (PKG=./path/to/package)"
	@echo "  make test-package-coverage - Run tests with coverage for a specific package (PKG=./path/to/package)"
	@echo "  make test-race         - Run tests with race detection"
	@echo "  make test-run          - Run tests matching a specific pattern (PATTERN=TestName)"
	@echo "  make test-run-coverage - Run tests with coverage matching a specific pattern (PATTERN=TestName)"
	@echo "  make test-timeout      - Run tests with timeout"
	@echo "  make test-timeout-behavior - Run tests that verify timeout behavior"
	@echo "  make test-unit         - Run unit tests only"
	@echo ""
	@echo "#################################################"
	@echo "# COVERAGE TARGETS"
	@echo "#################################################"
	@echo "  make test-coverage     - Generate test coverage report"
	@echo "  make test-coverage-func - Show test coverage by function"
	@echo "  make test-coverage-func-sorted - Show test coverage by function, sorted"
	@echo "  make test-coverage-summary - Show test coverage summary"
	@echo "  make test-coverage-view - View test coverage in browser"
	@echo ""
	@echo "#################################################"
	@echo "# PROFILING TARGETS"
	@echo "#################################################"
	@echo "  make profile-all       - Run all profiling"
	@echo "  make profile-block     - Run block profiling"
	@echo "  make profile-cpu       - Run CPU profiling"
	@echo "  make profile-mem       - Run memory profiling"
	@echo "  make profile-mutex     - Run mutex profiling"
	@echo "  make profile-trace     - Run execution tracing"
	@echo ""
	@echo "#################################################"
	@echo "# CODE QUALITY TARGETS"
	@echo "#################################################"
	@echo "  make check-go-version  - Check Go version consistency across files"
	@echo "  make deps              - Download dependencies"
	@echo "  make deps-graph        - Generate dependency graph"
	@echo "  make deps-upgrade      - Upgrade dependencies"
	@echo "  make fmt               - Format code"
	@echo "  make lint              - Run linters"
	@echo "  make pre-commit        - Run all pre-commit checks"
	@echo "  make tidy              - Tidy and verify Go modules"
	@echo "  make vuln-check        - Check for vulnerabilities in dependencies"
	@echo ""
	@echo "#################################################"
	@echo "# DOCUMENTATION TARGETS"
	@echo "#################################################"
	@echo "  make clean             - Remove build artifacts"
	@echo "  make docs              - Show documentation information"
	@echo "  make docs-pkgsite      - Generate and serve documentation with pkgsite"
	@echo "  make docs-static       - Generate static documentation"
	@echo "  make validate-readme   - Validate README.md files against the template"
	@echo ""
	@echo "#################################################"
	@echo "# DOCKER TARGETS"
	@echo "#################################################"
	@echo "  make airconfig         - Create a basic .air.toml configuration file if one doesn't exist"
	@echo "  make docker-build      - Build Docker image"
	@echo "  make docker-compose-build - Build all services"
	@echo "  make docker-compose-down - Stop all services with Docker Compose"
	@echo "  make docker-compose-logs - Show logs for all services"
	@echo "  make docker-compose-ps - Show status of all services"
	@echo "  make docker-compose-restart - Restart all services"
	@echo "  make docker-compose-up - Start all services with Docker Compose"
	@echo "  make docker-run        - Run Docker container"
	@echo "  make docker-stop       - Stop Docker container"
	@echo "  make dockerfile        - Create a basic Dockerfile"
	@echo ""
	@echo "#################################################"
	@echo "# DATABASE TARGETS"
	@echo "#################################################"
	@echo "  make db-init           - Initialize database based on DB_DRIVER environment variable"
	@echo ""
	@echo "#################################################"
	@echo "# DIAGRAM TARGETS"
	@echo "#################################################"
	@echo "  make plantuml          - Generate SVG files from PlantUML files"
	@echo "  make plantuml-deployment-container - Regenerate Deployment Container Diagram SVG file"

#################################################
# BUILD AND DEVELOPMENT TARGETS
#################################################

# Build the application
.PHONY: build
build: graphql-gen
	@echo "Building $(BINARY_NAME)..."
	mkdir -p bin
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_OUTPUT) $(MAIN_PATH)
	@echo "Build successful: $(BINARY_OUTPUT)"

# Build the application for all platforms and architectures
.PHONY: build-all
build-all: graphql-gen
	@echo "Building $(BINARY_NAME) for all platforms..."
	$(foreach platform,$(PLATFORMS),\
		$(foreach arch,$(ARCHITECTURES),\
			$(eval os := $(platform))\
			$(eval ext := $(if $(filter windows,$(platform)),.exe,))\
			mkdir -p bin/$(os)_$(arch) && \
			GOOS=$(os) GOARCH=$(arch) $(GOBUILD) $(LDFLAGS) -o bin/$(os)_$(arch)/$(BINARY_NAME)$(ext) $(MAIN_PATH) && \
			echo "Build successful: bin/$(os)_$(arch)/$(BINARY_NAME)$(ext)" ; \
		)\
	)
	@echo "All builds completed successfully"

# Run the application with hot reloading for development
.PHONY: dev
dev: airconfig
	@echo "Running $(BINARY_NAME) in development mode with hot reloading..."
	air -c .air.toml

# Generate GraphQL code
.PHONY: graphql-gen
graphql-gen:
	@echo "Generating GraphQL code..."
	@go get github.com/99designs/gqlgen/codegen/config@v0.17.75
	@go get github.com/99designs/gqlgen/internal/imports@v0.17.75
	@go get github.com/99designs/gqlgen/api@v0.17.75
	@go get github.com/99designs/gqlgen@v0.17.75
	@go get github.com/urfave/cli/v2
	@go run github.com/99designs/gqlgen generate --config interface/adapters/graphql/gqlgen.yml --verbose
	@echo "GraphQL code generated successfully"

# Initialize development environment
.PHONY: init
init:
	@echo "Initializing development environment..."
	$(GOGET) -u github.com/99designs/gqlgen
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GOGET) -u github.com/cosmtrek/air
	$(GOGET) -u golang.org/x/vuln/cmd/govulncheck
	$(GOMOD) download
	@echo "Development environment initialized successfully"

# Run the application locally
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	$(GORUN) $(MAIN_PATH)

#################################################
# TESTING TARGETS
#################################################

# Run all tests with coverage and generate a combined report
.PHONY: test-all
test-all:
	@echo "Running all tests with coverage..."
	mkdir -p ./coverage
	$(GOTEST) -v -coverprofile=./coverage/all.out ./...
	$(GOCMD) tool cover -html=./coverage/all.out -o ./coverage/all.html
	$(GOCMD) tool cover -func=./coverage/all.out
	@echo "All tests completed and coverage report generated: ./coverage/all.html"

# Run benchmarks
.PHONY: test-bench
test-bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...
	@echo "Benchmarks completed"

# Run integration tests only
.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -run Integration ./...
	@echo "Integration tests completed"

# Run tests for a specific package
.PHONY: test-package
test-package:
	@echo "Running tests for a specific package..."
	@echo "Usage: make test-package PKG=./path/to/package"
	@if [ "$(PKG)" = "" ]; then \
		echo "Error: PKG is required. Example: make test-package PKG=./core/domain"; \
		exit 1; \
	fi
	$(GOTEST) -v $(PKG)
	@echo "Package tests completed"

# Run tests with coverage for a specific package
.PHONY: test-package-coverage
test-package-coverage:
	@echo "Running tests with coverage for a specific package..."
	@echo "Usage: make test-package-coverage PKG=./path/to/package"
	@if [ "$(PKG)" = "" ]; then \
		echo "Error: PKG is required. Example: make test-package-coverage PKG=./core/domain"; \
		exit 1; \
	fi
	mkdir -p ./coverage
	$(GOTEST) -v -cover -coverprofile=./coverage/$(shell basename $(PKG)).out $(PKG)
	$(GOCMD) tool cover -html=./coverage/$(shell basename $(PKG)).out -o ./coverage/$(shell basename $(PKG)).html
	@echo "Package coverage report generated: ./coverage/$(shell basename $(PKG)).html"

# Run tests with race detection
.PHONY: test-race
test-race:
	@echo "Running tests with race detection..."
	$(GOTEST) -race -v ./...
	@echo "Race detection tests completed"

# Run tests matching a specific pattern
.PHONY: test-run
test-run:
	@echo "Running tests matching a specific pattern..."
	@echo "Usage: make test-run PATTERN=TestName"
	@if [ "$(PATTERN)" = "" ]; then \
		echo "Error: PATTERN is required. Example: make test-run PATTERN=TestCreateFamily"; \
		exit 1; \
	fi
	$(GOTEST) -v ./... -run $(PATTERN)
	@echo "Pattern tests completed"

# Run tests with coverage matching a specific pattern
.PHONY: test-run-coverage
test-run-coverage:
	@echo "Running tests with coverage matching a specific pattern..."
	@echo "Usage: make test-run-coverage PATTERN=TestName"
	@if [ "$(PATTERN)" = "" ]; then \
		echo "Error: PATTERN is required. Example: make test-run-coverage PATTERN=TestCreateFamily"; \
		exit 1; \
	fi
	mkdir -p ./coverage
	$(GOTEST) -v -coverprofile=./coverage/pattern.out ./... -run $(PATTERN)
	$(GOCMD) tool cover -html=./coverage/pattern.out -o ./coverage/pattern.html
	@echo "Pattern coverage report generated: ./coverage/pattern.html"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./... -coverprofile=coverage.out
	@echo "Tests completed"

# Run tests with timeout
.PHONY: test-timeout
test-timeout:
	@echo "Running tests with 30s timeout..."
	$(GOTEST) -timeout 30s ./...
	@echo "Timeout tests completed"

# Run specific timeout behavior tests
.PHONY: test-timeout-behavior
test-timeout-behavior:
	@echo "Running timeout behavior tests..."
	$(GOTEST) -v ./... -run "Test.*Timeout"
	@echo "Timeout behavior tests completed"

# Run unit tests only
.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	$(GOTEST) -v -short ./...
	@echo "Unit tests completed"

#################################################
# COVERAGE TARGETS
#################################################

# Generate test coverage report
.PHONY: test-coverage
test-coverage:
	@echo "Generating test coverage report..."
	mkdir -p ./coverage
	$(GOTEST) -coverprofile=./coverage/coverage.out ./...
	$(GOCMD) tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html
	@echo "Coverage report generated: ./coverage/coverage.html"

# Show test coverage by function
.PHONY: test-coverage-func
test-coverage-func:
	@echo "Generating test coverage by function..."
	mkdir -p ./coverage
	$(GOTEST) -coverprofile=./coverage/coverage.out ./...
	$(GOCMD) tool cover -func=./coverage/coverage.out
	@echo "Function coverage completed"

# Show test coverage by function, sorted by coverage percentage
.PHONY: test-coverage-func-sorted
test-coverage-func-sorted:
	@echo "Generating test coverage by function (sorted)..."
	mkdir -p ./coverage
	$(GOTEST) -coverprofile=./coverage/coverage.out ./...
	$(GOCMD) tool cover -func=./coverage/coverage.out | sort -k 3 -n
	@echo "Sorted function coverage completed"

# Show test coverage summary
.PHONY: test-coverage-summary
test-coverage-summary:
	@echo "Generating test coverage summary..."
	$(GOTEST) -cover ./...
	@echo "Coverage summary completed"

# View test coverage in browser
.PHONY: test-coverage-view
test-coverage-view: test-coverage
	@echo "Opening coverage report in browser..."
	$(GOCMD) tool cover -html=./coverage/coverage.out
	@echo "Coverage report opened in browser"

#################################################
# PROFILING TARGETS
#################################################

.PHONY: profile-all
profile-all: profile-cpu profile-mem profile-block profile-mutex profile-trace
	@echo "All profiling completed"

.PHONY: profile-block
profile-block:
	@echo "Running block profiling..."
	mkdir -p ./profiles
	$(GOTEST) -blockprofile=./profiles/block.prof -bench=. ./...
	@echo "Block profiling completed. Results in ./profiles/block.prof"
	@echo "View the profile with: go tool pprof ./profiles/block.prof"

.PHONY: profile-cpu
profile-cpu:
	@echo "Running CPU profiling..."
	mkdir -p ./profiles
	$(GOTEST) -cpuprofile=./profiles/cpu.prof -bench=. ./...
	@echo "CPU profiling completed. Results in ./profiles/cpu.prof"
	@echo "View the profile with: go tool pprof ./profiles/cpu.prof"

.PHONY: profile-mem
profile-mem:
	@echo "Running memory profiling..."
	mkdir -p ./profiles
	$(GOTEST) -memprofile=./profiles/mem.prof -bench=. ./...
	@echo "Memory profiling completed. Results in ./profiles/mem.prof"
	@echo "View the profile with: go tool pprof ./profiles/mem.prof"

.PHONY: profile-mutex
profile-mutex:
	@echo "Running mutex profiling..."
	mkdir -p ./profiles
	$(GOTEST) -mutexprofile=./profiles/mutex.prof -bench=. ./...
	@echo "Mutex profiling completed. Results in ./profiles/mutex.prof"
	@echo "View the profile with: go tool pprof ./profiles/mutex.prof"

.PHONY: profile-trace
profile-trace:
	@echo "Running execution tracing..."
	mkdir -p ./profiles
	$(GOTEST) -trace=./profiles/trace.out -bench=. ./...
	@echo "Execution tracing completed. Results in ./profiles/trace.out"
	@echo "View the trace with: go tool trace ./profiles/trace.out"

#################################################
# CODE QUALITY TARGETS
#################################################

# Check Go version consistency across files
.PHONY: check-go-version
check-go-version:
	@echo "Checking Go version consistency..."
	@echo "Expected Go version: $(GO_VERSION)"
	@MOD_VERSION=$$(grep -E "^go [0-9]+\.[0-9]+(\.[0-9]+)?" go.mod | awk '{print $$2}'); \
	WORK_VERSION=$$(grep -E "^go [0-9]+\.[0-9]+(\.[0-9]+)?" go.work | awk '{print $$2}'); \
	DOCKERFILE_VERSION=$$(grep -E "FROM golang:[0-9]+\.[0-9]+(\.[0-9]+)?-alpine" Dockerfile | sed -E 's/.*golang:([0-9]+\.[0-9]+(\.[0-9]+)?).*/\1/'); \
	echo "Go version in go.mod: $$MOD_VERSION"; \
	echo "Go version in go.work: $$WORK_VERSION"; \
	echo "Go version in Dockerfile: $$DOCKERFILE_VERSION"; \
	if [ "$$MOD_VERSION" != "$(GO_VERSION)" ] || [ "$$WORK_VERSION" != "$(GO_VERSION)" ] || [ "$$DOCKERFILE_VERSION" != "$(GO_VERSION)" ]; then \
		echo "WARNING: Go version inconsistency detected!"; \
		echo "Consider updating the files to use the same Go version."; \
		exit 1; \
	else \
		echo "Go version is consistent across all files."; \
	fi

# Dependency management targets
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	@echo "Dependencies downloaded"

.PHONY: deps-graph
deps-graph:
	@echo "Generating dependency graph..."
	go install github.com/kisielk/godepgraph@latest
	mkdir -p ./docs/deps
	godepgraph -s github.com/abitofhelp/family-service-graphql | dot -Tpng -o ./docs/deps/dependency-graph.png
	@echo "Dependency graph generated at ./docs/deps/dependency-graph.png"

.PHONY: deps-upgrade
deps-upgrade:
	@echo "Upgrading dependencies..."
	go get -u ./...
	$(GOMOD) tidy
	@echo "Dependencies upgraded"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	gofmt -s -w .
	@echo "Formatting completed"

# Run linters
.PHONY: lint
lint:
	@echo "Running linters..."
	$(GOLINT) run
	@echo "Linting completed"

# Pre-commit checks
.PHONY: pre-commit
pre-commit: tidy fmt lint test vuln-check validate-readme check-go-version enforce-naming-convention
	@echo "All pre-commit checks passed!"

# Enforce naming conventions for .puml and .svg files
.PHONY: enforce-naming-convention
enforce-naming-convention:
	@echo "Enforcing naming conventions for .puml and .svg files..."
	./DOCS/tools/scripts/enforce_naming_convention.sh
	@echo "Naming conventions enforced"

# Tidy and verify Go modules
.PHONY: tidy
tidy:
	@echo "Tidying Go modules..."
	$(GOMOD) tidy
	$(GOMOD) verify
	@echo "Go modules tidied and verified"

# Check for vulnerabilities in dependencies
.PHONY: vuln-check
vuln-check:
	@echo "Checking for vulnerabilities in dependencies..."
	$(GOVULNCHECK) ./...
	@echo "Vulnerability check completed"

#################################################
# DOCUMENTATION TARGETS
#################################################

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf bin
	rm -rf coverage
	rm -rf profiles
	rm -f $(BINARY_NAME)
	rm -f family-service
	rm -f coverage.out
	@echo "Clean completed"

# Documentation
.PHONY: docs
docs:
	@echo "Documentation is available in the docs directory:"
	@echo "- Software Requirements Specification (SRS): docs/SRS_FamilyService.md"
	@echo "- Software Design Document (SDD): docs/SDD_FamilyService.md"
	@echo "- Software Test Plan (STP): docs/STP_FamilyService.md"
	@echo "- Deployment Document: docs/Deployment_FamilyService.md"
	@echo "- Component README Template: COMPONENT_README_TEMPLATE.md"

.PHONY: docs-pkgsite
docs-pkgsite:
	@echo "Generating documentation with pkgsite..."
	go install golang.org/x/pkgsite/cmd/pkgsite@latest
	pkgsite -http=:6060 &
	@echo "Documentation server started at http://localhost:6060"

.PHONY: docs-static
docs-static:
	@echo "Generating static documentation..."
	mkdir -p ./docs/godoc
	go install golang.org/x/tools/cmd/godoc@latest
	godoc -url=/pkg/github.com/abitofhelp/family-service-graphql/ > ./docs/godoc/index.html
	@echo "Static documentation generated at ./docs/godoc/index.html"

# Validate README.md files against the template
.PHONY: validate-readme
validate-readme:
	@echo "Validating README.md files against the template..."
	$(GORUN) tools/readme_validator/main.go
	@echo "README validation completed successfully"

#################################################
# DOCKER TARGETS
#################################################

# Create a basic .air.toml configuration file if one doesn't exist
.PHONY: airconfig
airconfig:
	@echo "Creating .air.toml configuration file..."
	@if [ -f .air.toml ]; then \
		echo ".air.toml already exists. Skipping."; \
	else \
		echo "# .air.toml configuration file" > .air.toml; \
		echo "root = \"./\"" >> .air.toml; \
		echo "tmp_dir = \"tmp\"" >> .air.toml; \
		echo "" >> .air.toml; \
		echo "[build]" >> .air.toml; \
		echo "  cmd = \"go build -o ./tmp/main ./cmd/server/graphql\"" >> .air.toml; \
		echo "  bin = \"./tmp/main\"" >> .air.toml; \
		echo "  delay = 1000" >> .air.toml; \
		echo "  exclude_dir = [\"assets\", \"tmp\", \"vendor\", \"bin\"]" >> .air.toml; \
		echo "  include_ext = [\"go\", \"yaml\", \"yml\", \"graphql\"]" >> .air.toml; \
		echo "  exclude_regex = [\"_test\\.go\"]" >> .air.toml; \
		echo "" >> .air.toml; \
		echo "[log]" >> .air.toml; \
		echo "  time = true" >> .air.toml; \
		echo "" >> .air.toml; \
		echo "[color]" >> .air.toml; \
		echo "  main = \"magenta\"" >> .air.toml; \
		echo "  watcher = \"cyan\"" >> .air.toml; \
		echo "  build = \"yellow\"" >> .air.toml; \
		echo "  runner = \"green\"" >> .air.toml; \
		echo "" >> .air.toml; \
		echo "[screen]" >> .air.toml; \
		echo "  clear_on_rebuild = true" >> .air.toml; \
		echo ".air.toml created successfully."; \
	fi

# Docker Compose targets
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker-compose build
	@echo "Docker image built successfully"

.PHONY: docker-compose-build
docker-compose-build:
	@echo "Building all services..."
	docker-compose build
	@echo "All services built"

.PHONY: docker-compose-down
docker-compose-down:
	@echo "Stopping all services with Docker Compose..."
	docker-compose down
	@echo "All services stopped"

.PHONY: docker-compose-logs
docker-compose-logs:
	@echo "Showing logs for all services..."
	docker-compose logs -f

.PHONY: docker-compose-ps
docker-compose-ps:
	@echo "Showing status of all services..."
	docker-compose ps

.PHONY: docker-compose-restart
docker-compose-restart:
	@echo "Restarting all services..."
	docker-compose restart
	@echo "All services restarted"

.PHONY: docker-compose-up
docker-compose-up:
	@echo "Starting all services with Docker Compose..."
	docker-compose up -d
	@echo "All services started"

.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker-compose up -d
	@echo "Docker container started"

.PHONY: docker-stop
docker-stop:
	@echo "Stopping Docker container..."
	docker-compose down
	@echo "Docker container stopped"

# Create a basic Dockerfile if one doesn't exist
.PHONY: dockerfile
dockerfile:
	@echo "Creating Dockerfile..."
	@if [ -f Dockerfile ]; then \
		echo "WARNING: Dockerfile already exists. This command would create a simplified version."; \
		echo "The existing Dockerfile might have specific configurations that would be lost."; \
		echo "If you want to recreate the Dockerfile, delete it first with 'rm Dockerfile'."; \
	else \
		echo "FROM golang:$(GO_VERSION)-alpine3.21 AS builder" > Dockerfile; \
		echo "ENV CGO_ENABLED=1" >> Dockerfile; \
		echo "" >> Dockerfile; \
		echo "# Install build dependencies for CGO" >> Dockerfile; \
		echo "RUN apk add --no-cache gcc musl-dev" >> Dockerfile; \
		echo "" >> Dockerfile; \
		echo "WORKDIR /app" >> Dockerfile; \
		echo "COPY go.mod go.sum ./" >> Dockerfile; \
		echo "RUN go mod download" >> Dockerfile; \
		echo "COPY . ." >> Dockerfile; \
		echo "RUN go build -o family_service ./cmd/server/graphql" >> Dockerfile; \
		echo "" >> Dockerfile; \
		echo "FROM alpine:3.19" >> Dockerfile; \
		echo "RUN apk --no-cache add ca-certificates sqlite-libs" >> Dockerfile; \
		echo "WORKDIR /app" >> Dockerfile; \
		echo "COPY --from=builder /app/family_service ." >> Dockerfile; \
		echo "COPY --from=builder /app/config ./config" >> Dockerfile; \
		echo "COPY entrypoint.sh ." >> Dockerfile; \
		echo "COPY secrets ./secrets" >> Dockerfile; \
		echo "RUN chmod +x /app/entrypoint.sh" >> Dockerfile; \
		echo "" >> Dockerfile; \
		echo "EXPOSE 8089" >> Dockerfile; \
		echo "ENTRYPOINT [\"/app/entrypoint.sh\"]" >> Dockerfile; \
		echo "CMD [\"./family_service\"]" >> Dockerfile; \
		echo "" >> Dockerfile; \
		echo "HEALTHCHECK --interval=30s --timeout=3s \\" >> Dockerfile; \
		echo "  CMD wget --quiet --tries=1 --spider http://localhost:8089/health || exit 1" >> Dockerfile; \
		echo "Dockerfile created successfully."; \
	fi

#################################################
# DATABASE TARGETS
#################################################

# Database initialization
.PHONY: db-init
db-init:
	@echo "Initializing database based on DB_DRIVER environment variable..."
	@if [ "$(DB_DRIVER)" = "sqlite" ] || [ -z "$(DB_DRIVER)" ]; then \
		echo "Initializing SQLite database..."; \
		go run data/dev/sqlite/sqlite_init.go; \
	elif [ "$(DB_DRIVER)" = "postgres" ]; then \
		echo "Initializing PostgreSQL database..."; \
		echo "Make sure PostgreSQL is running and the database exists."; \
		echo "Run: psql -d familydb -f data/dev/postgres/postgresql_init.sql"; \
	elif [ "$(DB_DRIVER)" = "mongo" ]; then \
		echo "Initializing MongoDB database..."; \
		echo "Make sure MongoDB is running."; \
		echo "Run: mongosh < data/dev/mongo/mongodb_init.js"; \
	else \
		echo "Unknown DB_DRIVER: $(DB_DRIVER). Supported values are: sqlite, postgres, mongo"; \
		exit 1; \
	fi

#################################################
# DIAGRAM TARGETS
#################################################

# PlantUML targets
.PHONY: plantuml
plantuml:
	@echo "Generating SVG files from PlantUML files..."
	@if [ ! -f DOCS/tools/plantuml.jar ]; then \
		echo "Downloading PlantUML JAR file..."; \
		curl -L https://sourceforge.net/projects/plantuml/files/plantuml.jar/download -o DOCS/tools/plantuml.jar; \
	fi
	@echo "Generating SVG files..."
	@for file in DOCS/diagrams/*.puml; do \
		echo "Processing $$file..."; \
		java -jar DOCS/tools/plantuml.jar -tsvg $$file; \
	done
	@echo "SVG files generated successfully."

.PHONY: plantuml-deployment-container
plantuml-deployment-container:
	@echo "Regenerating Deployment Container Diagram SVG file..."
	@if [ ! -f DOCS/tools/plantuml.jar ]; then \
		echo "Downloading PlantUML JAR file..."; \
		curl -L https://sourceforge.net/projects/plantuml/files/plantuml.jar/download -o DOCS/tools/plantuml.jar; \
	fi
	@echo "Processing DOCS/diagrams/deployment_container_diagram.puml..."
	@java -jar DOCS/tools/plantuml.jar -tsvg DOCS/diagrams/deployment_container_diagram.puml
	@echo "Deployment Container Diagram SVG file regenerated successfully."
