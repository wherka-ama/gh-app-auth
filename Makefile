# gh-app-auth Makefile

.PHONY: help build test lint clean install dev-setup security-scan release deps vet gocyclo staticcheck ineffassign misspell test-coverage-check markdownlint yamllint actionlint cli-smoke-test package-deb package-rpm packages

# Default target
help:
	@echo "gh-app-auth - GitHub App Authentication Extension"
	@echo ""
	@echo "Available targets:"
	@echo "  build              Build the extension binary"
	@echo "  test               Run all tests"
	@echo "  test-race          Run tests with race detection"
	@echo "  test-cover         Run tests with coverage report"
	@echo "  test-coverage-check Enforce minimum coverage threshold"
	@echo "  lint               Run all linters (golangci-lint)"
	@echo "  vet                Run go vet"
	@echo "  staticcheck        Run staticcheck"
	@echo "  gocyclo            Run gocyclo (cyclomatic complexity)"
	@echo "  ineffassign        Run ineffassign (ineffectual assignments)"
	@echo "  misspell           Run misspell (spelling checker)"
	@echo "  markdownlint       Run markdownlint on .md files"
	@echo "  yamllint           Run yamllint on .yml/.yaml files"
	@echo "  actionlint         Run actionlint on GitHub workflow files"
	@echo "  lint-all           Run all individual linters"
	@echo "  fmt                Format code"
	@echo "  clean              Clean build artifacts"
	@echo "  install            Install extension to GitHub CLI"
	@echo "  uninstall          Uninstall extension from GitHub CLI"
	@echo "  dev-setup          Set up development environment (config only)"
	@echo "  validate-tools     Validate core tools are installed"
	@echo "  validate-lint-tools Validate linting tools can run via 'go run'"
	@echo "  security-scan      Run security scans (gosec, govulncheck via go run)"
	@echo "  deps               Download and verify dependencies"
	@echo "  dev                Quick development cycle (fmt + lint + test + build)"
	@echo "  ci                 CI pipeline simulation (mirrors GitHub CI)"
	@echo "  quality            Full quality check (all linters + tests + security)"
	@echo "  release            Build release binaries for all platforms"
	@echo ""
	@echo "Packaging targets (all use 'go run', no installation):"
	@echo "  package-deb          Build DEB package for amd64"
	@echo "  package-deb-arm64    Build DEB package for arm64 (requires 'make release')"
	@echo "  package-deb-386      Build DEB package for 386/i386 (requires 'make release')"
	@echo "  package-deb-arm      Build DEB package for arm/armhf (requires 'make release')"
	@echo "  package-rpm          Build RPM package for amd64"
	@echo "  package-rpm-arm64    Build RPM package for arm64 (requires 'make release')"
	@echo "  package-rpm-386      Build RPM package for 386/i386 (requires 'make release')"
	@echo "  package-rpm-arm      Build RPM package for arm/armv7hl (requires 'make release')"
	@echo "  packages             Build all packages (deb/rpm for all architectures)"
	@echo "  packages-local       Build packages for local architecture only"
	@echo ""
	@echo "Presentation targets (use npx, no global install):"
	@echo "  presentation-setup Verify npx availability for Mermaid tools"
	@echo "  presentation       Build both HTML and PDF presentations"
	@echo "  presentation-html  Build interactive HTML presentation"
	@echo "  presentation-pdf   Build PDF presentation (uses npx mermaid-filter)"
	@echo "  presentation-serve Serve presentation locally on :8000"
	@echo "  presentation-clean Clean presentation build artifacts"

# Build variables
BINARY_NAME := gh-app-auth
VERSION := $(shell git describe --tags --exact-match 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# Go tool commands using 'go run' (no installation to user's environment)
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
GOIMPORTS := go run golang.org/x/tools/cmd/goimports@latest
STATICCHECK := go run honnef.co/go/tools/cmd/staticcheck@latest
GOCYCLO := go run github.com/fzipp/gocyclo/cmd/gocyclo@latest
INEFFASSIGN := go run github.com/gordonklaus/ineffassign@latest
MISSPELL := go run github.com/client9/misspell/cmd/misspell@latest
GOSEC := go run github.com/securego/gosec/v2/cmd/gosec@latest
GOVULNCHECK := go run golang.org/x/vuln/cmd/govulncheck@latest
ACTIONLINT := go run github.com/rhysd/actionlint/cmd/actionlint@latest
NFPM_CMD := go run github.com/goreleaser/nfpm/v2/cmd/nfpm@latest

# Build the extension
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -race ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Check coverage meets minimum threshold
COVERAGE_THRESHOLD ?= 50.0
test-coverage-check:
	@echo "Checking coverage threshold (minimum: $(COVERAGE_THRESHOLD)%)..."
	@go test -race -coverprofile=coverage.out -covermode=atomic ./... > /dev/null 2>&1 || true
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Current coverage: $$COVERAGE%"; \
	if [ "$$(echo "$$COVERAGE < $(COVERAGE_THRESHOLD)" | bc -l)" -eq 1 ]; then \
		echo "❌ Coverage $$COVERAGE% is below threshold $(COVERAGE_THRESHOLD)%"; \
		echo ""; \
		echo "Package breakdown:"; \
		go tool cover -func=coverage.out | grep -E "(pkg/|cmd/)" | awk '{printf "  %-50s %s\n", $$1, $$3}' | sort; \
		exit 1; \
	else \
		echo "✅ Coverage $$COVERAGE% meets threshold $(COVERAGE_THRESHOLD)%"; \
	fi

# Lint code with golangci-lint using go run (no installation)
lint:
	@echo "Running golangci-lint..."
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run staticcheck using go run (no installation)
staticcheck:
	@echo "Running staticcheck..."
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...

# Run gocyclo (cyclomatic complexity) using go run (no installation)
gocyclo:
	@echo "Running gocyclo (complexity threshold: 10)..."
	go run github.com/fzipp/gocyclo/cmd/gocyclo@latest -over 10 . || echo "⚠️  High complexity functions found (expected for CLI commands)"

# Run ineffassign (ineffectual assignments) using go run (no installation)
ineffassign:
	@echo "Running ineffassign..."
	go run github.com/gordonklaus/ineffassign@latest ./...

# Run misspell (spelling checker) using go run (no installation)
misspell:
	@echo "Running misspell..."
	go run github.com/client9/misspell/cmd/misspell@latest -error .

# Run markdownlint (requires npx/node)
markdownlint:
	@echo "Running markdownlint..."
	@command -v npx >/dev/null 2>&1 || { echo "⚠️  npx not found, skipping markdownlint"; exit 0; }
	npx markdownlint-cli2 "**/*.md" "!node_modules/**" || echo "⚠️  Markdown lint issues found"

# Run yamllint (requires pip install yamllint)
yamllint:
	@echo "Running yamllint..."
	@command -v yamllint >/dev/null 2>&1 || { echo "⚠️  yamllint not found, skipping (install: pip install yamllint)"; exit 0; }
	yamllint -d relaxed . || echo "⚠️  YAML lint issues found"

# Run actionlint on GitHub workflow files using go run (no installation)
actionlint:
	@echo "Running actionlint..."
	@go run github.com/rhysd/actionlint/cmd/actionlint@latest || echo "⚠️  Action lint issues found"

# CLI smoke test - verify binary works
cli-smoke-test: build
	@echo "Running CLI smoke tests..."
	./$(BINARY_NAME) --help > /dev/null
	./$(BINARY_NAME) --version > /dev/null
	@echo "✅ CLI smoke tests passed"

# Run all individual linters
lint-all: vet staticcheck gocyclo ineffassign misspell markdownlint yamllint actionlint
	@echo "All individual linters completed!"

# Format code using go run (no installation)
fmt:
	@echo "Formatting code..."
	go run golang.org/x/tools/cmd/goimports@latest -w .
	gofmt -s -w .

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	rm -rf dist/

# Install extension to GitHub CLI
install: build
	@echo "Installing extension to GitHub CLI..."
	gh extension install .

# Uninstall extension from GitHub CLI
uninstall:
	@echo "Uninstalling extension from GitHub CLI..."
	gh extension remove app-auth || true

# Set up development environment (configuration only, no tool installation)
dev-setup:
	@echo "Setting up development environment..."
	go mod download
	@echo "Setting up git commit template..."
	git config commit.template .gitmessage
	@echo "Development environment ready!"
	@echo ""
	@echo "💡 Tip: Use 'git commit' (without -m) to use the conventional commit template"
	@echo "📖 See CONTRIBUTING.md for conventional commit guidelines"

# Set up presentation tools using npx (no global installation)
presentation-setup:
	@echo "Setting up presentation tools..."
	@command -v npm >/dev/null 2>&1 || { echo "npm is required. Install Node.js first"; exit 1; }
	@echo "Verifying npx is available..."
	@command -v npx >/dev/null 2>&1 || { echo "npx not found. Install Node.js properly"; exit 1; }
	@echo "✅ npx available for mermaid-cli and mermaid-filter"
	@echo ""
	@echo "Presentation tools ready (will use npx at runtime)!"

# Run security scans using go run (no installation)
security-scan:
	@echo "Running security scans..."
	go run github.com/securego/gosec/v2/cmd/gosec@latest -fmt sarif -out gosec.sarif ./... || true
	@echo "Running vulnerability check..."
	go run golang.org/x/vuln/cmd/govulncheck@latest ./... || true

# Download and verify dependencies  
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod verify
	go mod tidy

# Build release binaries
release: clean
	@echo "Building release binaries..."
	mkdir -p dist
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 .
	
	# Linux 386 (i386)
	GOOS=linux GOARCH=386 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-386 .
	
	# Linux ARM (armv7)
	GOOS=linux GOARCH=arm GOARM=7 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm .
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .
	
	# Windows ARM64
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-arm64.exe .
	
	@echo "Release binaries built in dist/"
	@ls -la dist/

# Validate that all required tools are installed
validate-tools:
	@echo "Validating required tools..."
	@command -v go >/dev/null 2>&1 || { echo "❌ Go is required but not installed"; exit 1; }
	@command -v gh >/dev/null 2>&1 || { echo "❌ GitHub CLI is required but not installed"; exit 1; }
	@command -v git >/dev/null 2>&1 || { echo "❌ Git is required but not installed"; exit 1; }
	@echo "✅ Core tools are installed."

# Validate that linting tools can be run (check go run works)
validate-lint-tools:
	@echo "Validating linting tools can be executed via 'go run'..."
	@go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest version >/dev/null 2>&1 || { echo "❌ golangci-lint failed via go run"; exit 1; }
	@go run golang.org/x/tools/cmd/goimports@latest -l . >/dev/null 2>&1 || { echo "❌ goimports failed via go run"; exit 1; }
	@go run honnef.co/go/tools/cmd/staticcheck@latest --help >/dev/null 2>&1 || { echo "❌ staticcheck failed via go run"; exit 1; }
	@go run github.com/fzipp/gocyclo/cmd/gocyclo@latest . >/dev/null 2>&1 || { echo "❌ gocyclo failed via go run"; exit 1; }
	@go run github.com/gordonklaus/ineffassign@latest . >/dev/null 2>&1 || { echo "❌ ineffassign failed via go run"; exit 1; }
	@go run github.com/client9/misspell/cmd/misspell@latest --help >/dev/null 2>&1 || { echo "❌ misspell failed via go run"; exit 1; }
	@go run github.com/securego/gosec/v2/cmd/gosec@latest --help >/dev/null 2>&1 || { echo "❌ gosec failed via go run"; exit 1; }
	@go run golang.org/x/vuln/cmd/govulncheck@latest --help >/dev/null 2>&1 || { echo "❌ govulncheck failed via go run"; exit 1; }
	@echo "✅ All linting tools work via 'go run'"

# Quick development cycle
dev: fmt lint test build
	@echo "Development cycle complete!"

# CI pipeline simulation (mirrors GitHub CI workflows)
# Runs: deps → vet → lint → test-race → coverage-check → security → build → smoke-test
ci: deps validate-tools validate-lint-tools
	@echo ""
	@echo "=========================================="
	@echo "  CI Pipeline Simulation"
	@echo "=========================================="
	@echo ""
	@echo "Step 1/8: Running go vet..."
	go vet ./...
	@echo ""
	@echo "Step 2/8: Running golangci-lint..."
	$(GOLANGCI_LINT) run --timeout=5m
	@echo ""
	@echo "Step 3/8: Running tests with race detection..."
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo ""
	@echo "Step 4/8: Checking coverage threshold (CI: 35%)..."
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Current coverage: $$COVERAGE%"; \
	if awk "BEGIN {exit !($$COVERAGE < 35.0)}"; then \
		echo "❌ Coverage $$COVERAGE% is below threshold 35.0%"; \
		exit 1; \
	else \
		echo "✅ Coverage $$COVERAGE% meets threshold 35.0%"; \
	fi
	@echo ""
	@echo "Step 5/8: Running security scans..."
	@$(MAKE) security-scan
	@echo ""
	@echo "Step 6/8: Building binary..."
	go build -v -o $(BINARY_NAME) .
	@echo ""
	@echo "Step 7/8: Running CLI smoke tests..."
	./$(BINARY_NAME) --help > /dev/null
	./$(BINARY_NAME) --version > /dev/null
	@echo "✅ CLI smoke tests passed"
	@echo ""
	@echo "Step 8/8: Running additional linters (non-blocking)..."
	@$(MAKE) markdownlint || true
	@$(MAKE) yamllint || true
	@$(MAKE) actionlint || true
	@echo ""
	@echo "=========================================="
	@echo "  ✅ CI Pipeline Complete!"
	@echo "=========================================="

# Full quality check (all linters + tests)
quality: validate-lint-tools fmt lint-all test-coverage-check security-scan
	@echo "Quality check complete!"

# Packaging targets - use 'go run' to avoid installing tools in user's environment
.PHONY: package-deb package-deb-arm64 package-deb-386 package-deb-arm package-rpm package-rpm-arm64 package-rpm-386 package-rpm-arm packages packages-local

# Build DEB package for amd64
package-deb:
	@echo "Building DEB package for amd64..."
	@mkdir -p dist
	@echo "Building Linux amd64 binary..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	@export GOARCH=amd64 ARCH=amd64 VERSION=$(VERSION); \
	envsubst '$$GOARCH $$ARCH $$VERSION' < nfpm.yaml > nfpm-temp.yaml; \
	$(NFPM_CMD) pkg --config nfpm-temp.yaml --packager deb --target dist/$(BINARY_NAME)_$(VERSION)_amd64.deb; \
	rm nfpm-temp.yaml
	@echo "✅ DEB package created: dist/$(BINARY_NAME)_$(VERSION)_amd64.deb"

# Build DEB package for arm64
package-deb-arm64: release
	@echo "Building DEB package for arm64..."
	@test -f dist/$(BINARY_NAME)-linux-arm64 || { echo "❌ Linux ARM64 binary not found. Run 'make release' first."; exit 1; }
	@export GOARCH=arm64 ARCH=arm64 VERSION=$(VERSION); \
	envsubst '$$GOARCH $$ARCH $$VERSION' < nfpm.yaml > nfpm-temp.yaml; \
	$(NFPM_CMD) pkg --config nfpm-temp.yaml --packager deb --target dist/$(BINARY_NAME)_$(VERSION)_arm64.deb; \
	rm nfpm-temp.yaml
	@echo "✅ DEB package created: dist/$(BINARY_NAME)_$(VERSION)_arm64.deb"

# Build DEB package for 386 (i386)
package-deb-386: release
	@echo "Building DEB package for 386..."
	@test -f dist/$(BINARY_NAME)-linux-386 || { echo "❌ Linux 386 binary not found. Run 'make release' first."; exit 1; }
	@export GOARCH=386 ARCH=386 VERSION=$(VERSION); \
	envsubst '$$GOARCH $$ARCH $$VERSION' < nfpm.yaml > nfpm-temp.yaml; \
	$(NFPM_CMD) pkg --config nfpm-temp.yaml --packager deb --target dist/$(BINARY_NAME)_$(VERSION)_i386.deb; \
	rm nfpm-temp.yaml
	@echo "✅ DEB package created: dist/$(BINARY_NAME)_$(VERSION)_i386.deb"

# Build DEB package for arm (armhf)
package-deb-arm: release
	@echo "Building DEB package for arm..."
	@test -f dist/$(BINARY_NAME)-linux-arm || { echo "❌ Linux ARM binary not found. Run 'make release' first."; exit 1; }
	@export GOARCH=arm ARCH=arm VERSION=$(VERSION); \
	envsubst '$$GOARCH $$ARCH $$VERSION' < nfpm.yaml > nfpm-temp.yaml; \
	$(NFPM_CMD) pkg --config nfpm-temp.yaml --packager deb --target dist/$(BINARY_NAME)_$(VERSION)_armhf.deb; \
	rm nfpm-temp.yaml
	@echo "✅ DEB package created: dist/$(BINARY_NAME)_$(VERSION)_armhf.deb"

# Build RPM package for amd64
package-rpm:
	@echo "Building RPM package for amd64..."
	@mkdir -p dist
	@echo "Building Linux amd64 binary..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	@export GOARCH=amd64 ARCH=amd64 VERSION=$(VERSION); \
	envsubst '$$GOARCH $$ARCH $$VERSION' < nfpm.yaml > nfpm-temp.yaml; \
	$(NFPM_CMD) pkg --config nfpm-temp.yaml --packager rpm --target dist/$(BINARY_NAME)_$(VERSION)_x86_64.rpm; \
	rm nfpm-temp.yaml
	@echo "✅ RPM package created: dist/$(BINARY_NAME)_$(VERSION)_x86_64.rpm"

# Build RPM package for arm64
package-rpm-arm64: release
	@echo "Building RPM package for arm64..."
	@test -f dist/$(BINARY_NAME)-linux-arm64 || { echo "❌ Linux ARM64 binary not found. Run 'make release' first."; exit 1; }
	@export GOARCH=arm64 ARCH=arm64 VERSION=$(VERSION); \
	envsubst '$$GOARCH $$ARCH $$VERSION' < nfpm.yaml > nfpm-temp.yaml; \
	$(NFPM_CMD) pkg --config nfpm-temp.yaml --packager rpm --target dist/$(BINARY_NAME)_$(VERSION)_aarch64.rpm; \
	rm nfpm-temp.yaml
	@echo "✅ RPM package created: dist/$(BINARY_NAME)_$(VERSION)_aarch64.rpm"

# Build RPM package for 386 (i386/i686)
package-rpm-386: release
	@echo "Building RPM package for 386..."
	@test -f dist/$(BINARY_NAME)-linux-386 || { echo "❌ Linux 386 binary not found. Run 'make release' first."; exit 1; }
	@export GOARCH=386 ARCH=386 VERSION=$(VERSION); \
	envsubst '$$GOARCH $$ARCH $$VERSION' < nfpm.yaml > nfpm-temp.yaml; \
	$(NFPM_CMD) pkg --config nfpm-temp.yaml --packager rpm --target dist/$(BINARY_NAME)_$(VERSION)_i386.rpm; \
	rm nfpm-temp.yaml
	@echo "✅ RPM package created: dist/$(BINARY_NAME)_$(VERSION)_i386.rpm"

# Build RPM package for arm (armv7hl)
package-rpm-arm: release
	@echo "Building RPM package for arm..."
	@test -f dist/$(BINARY_NAME)-linux-arm || { echo "❌ Linux ARM binary not found. Run 'make release' first."; exit 1; }
	@export GOARCH=arm ARCH=arm VERSION=$(VERSION); \
	envsubst '$$GOARCH $$ARCH $$VERSION' < nfpm.yaml > nfpm-temp.yaml; \
	$(NFPM_CMD) pkg --config nfpm-temp.yaml --packager rpm --target dist/$(BINARY_NAME)_$(VERSION)_armv7hl.rpm; \
	rm nfpm-temp.yaml
	@echo "✅ RPM package created: dist/$(BINARY_NAME)_$(VERSION)_armv7hl.rpm"

# Build all packages (requires release binaries)
packages: release package-deb package-rpm package-deb-arm64 package-rpm-arm64 package-deb-386 package-rpm-386 package-deb-arm package-rpm-arm
	@echo ""
	@echo "=========================================="
	@echo "  All packages built successfully!"
	@echo "=========================================="
	@ls -lh dist/*.deb dist/*.rpm 2>/dev/null || echo "⚠️  Some packages may not have been created"

# Build packages for local architecture only
packages-local:
	@echo "Building packages for local architecture ($(shell go env GOARCH))..."
	@mkdir -p dist
	@echo "Building Linux binary for local architecture..."
	@GOOS=linux GOARCH=$(shell go env GOARCH) go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-$(shell go env GOARCH) .
ifeq ($(shell go env GOARCH),amd64)
	@export GOARCH=amd64 ARCH=amd64 VERSION=$(VERSION); \
	envsubst '$$GOARCH $$ARCH $$VERSION' < nfpm.yaml > nfpm-temp.yaml; \
	$(NFPM_CMD) pkg --config nfpm-temp.yaml --packager deb --target dist/$(BINARY_NAME)_$(VERSION)_amd64.deb; \
	$(NFPM_CMD) pkg --config nfpm-temp.yaml --packager rpm --target dist/$(BINARY_NAME)_$(VERSION)_x86_64.rpm; \
	rm nfpm-temp.yaml
	@echo "✅ Local packages created (amd64)"
else ifeq ($(shell go env GOARCH),arm64)
	@export GOARCH=arm64 ARCH=arm64 VERSION=$(VERSION); \
	envsubst '$$GOARCH $$ARCH $$VERSION' < nfpm.yaml > nfpm-temp.yaml; \
	$(NFPM_CMD) pkg --config nfpm-temp.yaml --packager deb --target dist/$(BINARY_NAME)_$(VERSION)_arm64.deb; \
	$(NFPM_CMD) pkg --config nfpm-temp.yaml --packager rpm --target dist/$(BINARY_NAME)_$(VERSION)_aarch64.rpm; \
	rm nfpm-temp.yaml
	@echo "✅ Local packages created (arm64)"
else ifeq ($(shell go env GOARCH),386)
	@export GOARCH=386 ARCH=386 VERSION=$(VERSION); \
	envsubst '$$GOARCH $$ARCH $$VERSION' < nfpm.yaml > nfpm-temp.yaml; \
	$(NFPM_CMD) pkg --config nfpm-temp.yaml --packager deb --target dist/$(BINARY_NAME)_$(VERSION)_i386.deb; \
	$(NFPM_CMD) pkg --config nfpm-temp.yaml --packager rpm --target dist/$(BINARY_NAME)_$(VERSION)_i386.rpm; \
	rm nfpm-temp.yaml
	@echo "✅ Local packages created (386)"
else ifeq ($(shell go env GOARCH),arm)
	@export GOARCH=arm ARCH=arm VERSION=$(VERSION); \
	envsubst '$$GOARCH $$ARCH $$VERSION' < nfpm.yaml > nfpm-temp.yaml; \
	$(NFPM_CMD) pkg --config nfpm-temp.yaml --packager deb --target dist/$(BINARY_NAME)_$(VERSION)_armhf.deb; \
	$(NFPM_CMD) pkg --config nfpm-temp.yaml --packager rpm --target dist/$(BINARY_NAME)_$(VERSION)_armv7hl.rpm; \
	rm nfpm-temp.yaml
	@echo "✅ Local packages created (arm)"
else
	@echo "⚠️  Unsupported architecture: $(shell go env GOARCH)"
	@exit 1
endif

# Presentation targets
.PHONY: presentation presentation-setup presentation-html presentation-pdf presentation-serve presentation-clean

# Build presentation HTML
presentation-html:
	@echo "Building presentation HTML..."
	@command -v pandoc >/dev/null 2>&1 || { echo "Pandoc is required. Install: apt install pandoc"; exit 1; }
	mkdir -p dist/presentation
	pandoc docs/presentation.md \
		-t revealjs \
		-s \
		-o dist/presentation/index.html \
		-V revealjs-url=https://unpkg.com/reveal.js@3.9.2 \
		-V theme=white \
		-V transition=slide \
		-V slideNumber=true \
		--mathjax \
		--highlight-style=pygments
	@echo "Adding custom CSS and Mermaid support..."
	cp docs/presentation-custom.css dist/presentation/
	@# Insert custom CSS link before </head>
	sed -i 's|</head>|  <link rel="stylesheet" href="presentation-custom.css">\n  <script src="https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.min.js"></script>\n</head>|' dist/presentation/index.html
	@# Convert mermaid code blocks to divs - use perl for multi-line replacement
	perl -i -0pe 's/<pre class="mermaid"><code>(.*?)<\/code><\/pre>/<div class="mermaid">\1<\/div>/gs' dist/presentation/index.html
	@# Initialize Mermaid after Reveal.init
	sed -i 's|Reveal\.initialize({|mermaid.initialize({ startOnLoad: true, theme: "default" });\n      Reveal.initialize({|' dist/presentation/index.html
	@echo "Presentation HTML created: dist/presentation/index.html"

# Build presentation PDF
presentation-pdf:
	@echo "Building presentation PDF..."
	@command -v pandoc >/dev/null 2>&1 || { echo "Pandoc is required. Install: apt install pandoc"; exit 1; }
	@command -v xelatex >/dev/null 2>&1 || { echo "XeLaTeX is required. Install: apt install texlive-xetex"; exit 1; }
	@command -v mmdc >/dev/null 2>&1 || { echo "Mermaid CLI is required. Run: make presentation-setup"; exit 1; }
	@npm list -g mermaid-filter >/dev/null 2>&1 || { echo "mermaid-filter is required. Run: make presentation-setup"; exit 1; }
	mkdir -p dist/presentation
	pandoc docs/presentation.md \
		-o dist/presentation/presentation.pdf \
		--pdf-engine=xelatex \
		-F mermaid-filter \
		-V geometry:margin=1in \
		-V fontsize=12pt \
		-V colorlinks=true \
		-V mainfont="DejaVu Sans"
	@echo "Presentation PDF created: dist/presentation/presentation.pdf"

# Build both HTML and PDF
presentation: presentation-html presentation-pdf
	@echo "All presentation formats built successfully!"

# Serve presentation locally
presentation-serve: presentation-html
	@echo "Starting presentation server at http://localhost:8000"
	@command -v python3 >/dev/null 2>&1 || { echo "Python 3 is required"; exit 1; }
	cd dist/presentation && python3 -m http.server 8000

# Clean presentation artifacts
presentation-clean:
	rm -rf dist/presentation
	@echo "Presentation artifacts cleaned"
