# 定义项目名称
BINARY_NAME=go-ldap-admin

# 定义输出目录
OUTPUT_DIR=bin


.PHONY: default
default: help

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build -o ${BINARY_NAME} main.go

.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME} main.go

.PHONY: lint
lint:
	env GOGC=25 golangci-lint run --fix -j 8 -v ./... --timeout=5m --skip-files="public/client/feishu/feishu.go"

.PHONY: gox-linux
gox-linux:
	CGO_ENABLED=0 gox -osarch="linux/amd64 linux/arm64" -output="${OUTPUT_DIR}/${BINARY_NAME}_{{.OS}}_{{.Arch}}"
	@for b in $$(ls ${OUTPUT_DIR}); do \
		OUTPUT_FILE="${OUTPUT_DIR}/$$b"; \
		upx -9 -q "$$OUTPUT_FILE"; \
	done

.PHONY: gox-all
gox-all:
	CGO_ENABLED=0 gox -osarch="darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 linux/ppc64le windows/amd64" -output="${OUTPUT_DIR}/${BINARY_NAME}_{{.OS}}_{{.Arch}}"
	@for b in $$(ls ${OUTPUT_DIR}); do \
		FILENAME=$$(basename -s .exe "$$b"); \
		GOOS=$$(echo "$$FILENAME" | rev | cut -d'_' -f2 | rev); \
		GOARCH=$$(echo "$$FILENAME" | rev | cut -d'_' -f1 | rev); \
		OUTPUT_FILE="${OUTPUT_DIR}/$$b"; \
		if [ "$$GOOS" = "windows" ] && [ "$$GOARCH" = "arm64" ]; then \
			echo "跳过 $$OUTPUT_FILE (Windows/arm64 不压缩)"; \
		elif [ "$$GOOS" = "darwin" ]; then \
			echo "压缩 macOS 文件: $$OUTPUT_FILE"; \
			upx --force-macos -fq -9 "$$OUTPUT_FILE"; \
		else \
			echo "压缩通用文件: $$OUTPUT_FILE"; \
			upx -q -9 "$$OUTPUT_FILE"; \
		fi; \
	done

.PHONY: clean
clean:
	@rm -rf ${OUTPUT_DIR}

# 帮助信息
.PHONY: help
help:
	@echo "参数:"
	@echo "  run         运行项目"
	@echo "  build       为当前平台构建可执行文件"
	@echo "  gox-linux   为Linux平台构建可执行文件"
	@echo "  gox-all     为所有平台构建可执行文件"
	@echo "  clean       清理生成的可执行文件"
	@echo "  lint        代码格式检查"
	@echo "  help        显示帮助信息"