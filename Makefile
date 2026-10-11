# Define the project name
BINARY_NAME=go-ldap-admin

# Define the output directory
OUTPUT_DIR=bin

# Version information
VERSION ?= $(shell git tag --points-at HEAD 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date '+%Y-%m-%d %H:%M:%S')
GO_VERSION := $(shell go version | awk '{print $$3}')

# ldflags
LDFLAGS := -X 'github.com/eryajf/go-ldap-admin/public/version.Version=$(VERSION)' \
           -X 'github.com/eryajf/go-ldap-admin/public/version.GitCommit=$(GIT_COMMIT)' \
           -X 'github.com/eryajf/go-ldap-admin/public/version.BuildTime=$(BUILD_TIME)' \
           -X 'github.com/eryajf/go-ldap-admin/public/version.GoVersion=$(GO_VERSION)'

.PHONY: default
default: help

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o ${BINARY_NAME} main.go

.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "$(LDFLAGS)" -o ${BINARY_NAME} main.go

.PHONY: lint
lint:
	env GOGC=25 golangci-lint run --fix -j 8 -v ./... --timeout=5m

.PHONY: gox-linux
gox-linux:
	CGO_ENABLED=0 gox -osarch="linux/amd64 linux/arm64" -ldflags "$(LDFLAGS)" -output="${OUTPUT_DIR}/${BINARY_NAME}_{{.OS}}_{{.Arch}}"
	@for b in $$(ls ${OUTPUT_DIR}); do \
		OUTPUT_FILE="${OUTPUT_DIR}/$$b"; \
		upx -9 -q "$$OUTPUT_FILE"; \
	done

.PHONY: gox-all
gox-all:
	CGO_ENABLED=0 gox -osarch="darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 linux/ppc64le windows/amd64" -ldflags "$(LDFLAGS)" -output="${OUTPUT_DIR}/${BINARY_NAME}_{{.OS}}_{{.Arch}}"
	@for b in $$(ls ${OUTPUT_DIR}); do \
		FILENAME=$$(basename -s .exe "$$b"); \
		GOOS=$$(echo "$$FILENAME" | rev | cut -d'_' -f2 | rev); \
		GOARCH=$$(echo "$$FILENAME" | rev | cut -d'_' -f1 | rev); \
		OUTPUT_FILE="${OUTPUT_DIR}/$$b"; \
		if [ "$$GOOS" = "windows" ] && [ "$$GOARCH" = "arm64" ]; then \
			echo "Skipping $$OUTPUT_FILE (Windows/arm64 is not compressed)"; \
		elif [ "$$GOOS" = "darwin" ]; then \
			echo "Compressing macOS file: $$OUTPUT_FILE"; \
			upx --force-macos -fq -9 "$$OUTPUT_FILE"; \
		else \
			echo "Compressing binary: $$OUTPUT_FILE"; \
			upx -q -9 "$$OUTPUT_FILE"; \
		fi; \
	done

.PHONY: clean
clean:
	@rm -rf ${OUTPUT_DIR}

# Help information
.PHONY: help
help:
	@echo "Arguments:"
	@echo "  run         Run the project"
	@echo "  build       Build an executable for the current platform"
	@echo "  gox-linux   Build executables for Linux platforms"
	@echo "  gox-all     Build executables for all platforms"
	@echo "  clean       Remove generated executables"
	@echo "  lint        Check code formatting"
	@echo "  help        Display this help message"