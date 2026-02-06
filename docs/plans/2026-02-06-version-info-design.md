# 版本信息接口设计文档

## 概述

为项目添加版本信息返回接口，用于查询系统的版本号、Git 提交哈希、构建时间、Go 版本和编译平台信息。

## 需求

返回以下版本信息：
- Version: 版本号（如 v1.0.0）
- GitCommit: Git 提交哈希（短格式）
- BuildTime: 构建时间
- GoVersion: Go 编译器版本
- Platform: 编译平台（如 linux/amd64）

## 技术方案

### 1. 版本信息注入方式

采用 **编译时注入** 方式，通过 `go build -ldflags` 在编译时将版本信息注入到代码中。

**优势：**
- 自动化程度高，无需手动维护
- Git commit 和构建时间自动获取，确保准确性
- 符合 Go 社区最佳实践
- 适配所有构建方式（本地 make、GitHub Actions、CNB）

### 2. 架构设计

**组件结构：**
1. `public/version/version.go` - 版本信息包，定义全局变量
2. `controller/base_controller.go` - 添加 GetVersion 方法
3. `logic/base_logic.go` - 添加 GetVersion 业务逻辑
4. `model/request/base_request.go` - 添加 BaseVersionReq 结构
5. `routes/base_routes.go` - 注册 `/base/version` 路由
6. `Makefile` - 修改构建命令，注入版本信息

**数据流：**
```
编译时 (Makefile)
  → 注入版本信息到 version 包全局变量
  → API 请求 /api/base/version
  → Controller → Logic → 读取 version 包变量
  → 返回 JSON 响应
```

### 3. 代码实现

#### 3.1 版本信息包 (`public/version/version.go`)

```go
package version

var (
    Version   string // 版本号（如v1.0.0）
    GitCommit string // Git提交哈希（短格式）
    BuildTime string // 构建时间
    GoVersion string // Go编译器版本
    Platform  string // 编译平台（如linux/amd64）
)

// GetVersion 获取版本信息
func GetVersion() map[string]string {
    return map[string]string{
        "version":   Version,
        "gitCommit": GitCommit,
        "buildTime": BuildTime,
        "goVersion": GoVersion,
        "platform":  Platform,
    }
}
```

#### 3.2 Controller 层

在 `controller/base_controller.go` 中添加：

```go
// GetVersion 获取版本信息
// @Summary 获取版本信息
// @Description 获取系统版本号、Git提交哈希和构建时间
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/version [get]
func (m *BaseController) GetVersion(c *gin.Context) {
    req := new(request.BaseVersionReq)
    Run(c, req, func() (any, any) {
        return logic.Base.GetVersion(c, req)
    })
}
```

#### 3.3 Logic 层

在 `logic/base_logic.go` 中添加：

```go
import (
    "github.com/eryajf/go-ldap-admin/public/version"
    // ... 其他导入
)

// GetVersion 获取版本信息
func (l BaseLogic) GetVersion(c *gin.Context, req *request.BaseVersionReq) (data interface{}, err error) {
    return version.GetVersion(), nil
}
```

#### 3.4 Request 结构

在 `model/request/base_request.go` 中添加：

```go
// BaseVersionReq 获取版本信息请求
type BaseVersionReq struct{}
```

#### 3.5 路由注册

在 `routes/base_routes.go` 的 `InitBaseRoutes` 函数中添加：

```go
base.GET("version", controller.Base.GetVersion) // 获取版本信息
```

#### 3.6 Makefile 修改

```makefile
# 版本信息
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S')
GO_VERSION := $(shell go version | awk '{print $$3}')

# ldflags
LDFLAGS := -X 'github.com/eryajf/go-ldap-admin/public/version.Version=$(VERSION)' \
           -X 'github.com/eryajf/go-ldap-admin/public/version.GitCommit=$(GIT_COMMIT)' \
           -X 'github.com/eryajf/go-ldap-admin/public/version.BuildTime=$(BUILD_TIME)' \
           -X 'github.com/eryajf/go-ldap-admin/public/version.GoVersion=$(GO_VERSION)'

# 修改 build 命令
build:
	go build -ldflags "$(LDFLAGS) -X 'github.com/eryajf/go-ldap-admin/public/version.Platform=$(shell go env GOOS)/$(shell go env GOARCH)'" -o ${BINARY_NAME} main.go

# 修改 build-linux 命令
build-linux:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "$(LDFLAGS) -X 'github.com/eryajf/go-ldap-admin/public/version.Platform=linux/amd64'" -o ${BINARY_NAME} main.go

# 修改 gox-linux 命令
gox-linux:
	CGO_ENABLED=0 gox -osarch="linux/amd64 linux/arm64" \
		-ldflags "$(LDFLAGS) -X 'github.com/eryajf/go-ldap-admin/public/version.Platform={{.OS}}/{{.Arch}}'" \
		-output="${OUTPUT_DIR}/${BINARY_NAME}_{{.OS}}_{{.Arch}}"
	@for b in $$(ls ${OUTPUT_DIR}); do \
		OUTPUT_FILE="${OUTPUT_DIR}/$$b"; \
		upx -9 -q "$$OUTPUT_FILE"; \
	done

# 修改 gox-all 命令
gox-all:
	CGO_ENABLED=0 gox -osarch="darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 linux/ppc64le windows/amd64" \
		-ldflags "$(LDFLAGS) -X 'github.com/eryajf/go-ldap-admin/public/version.Platform={{.OS}}/{{.Arch}}'" \
		-output="${OUTPUT_DIR}/${BINARY_NAME}_{{.OS}}_{{.Arch}}"
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
```

### 4. API 响应示例

**请求：**
```
GET /api/base/version
```

**响应：**
```json
{
    "code": 200,
    "data": {
        "version": "v1.5.0",
        "gitCommit": "dc03d68",
        "buildTime": "2026-02-06 10:30:45",
        "goVersion": "go1.21.5",
        "platform": "linux/amd64"
    },
    "msg": "success"
}
```

## 构建方式兼容性

### 本地构建
- `make build` - 自动注入版本信息
- `make build-linux` - Linux 平台构建
- `make gox-linux` - 多架构 Linux 构建
- `make gox-all` - 全平台构建

### GitHub Actions
- `build-docker-image.yml` - 使用 `make gox-linux`，自动注入
- `buildAndPush-binary-to-release.yml` - 需要更新以支持 ldflags

### CNB 构建
- `.cnb/workflows/build-docker-images.yml` - 使用 `make gox-linux`，自动注入

## 实施步骤

1. 创建 `public/version/version.go` 文件
2. 修改 `Makefile` 添加版本信息注入
3. 在 `model/request/base_request.go` 添加 BaseVersionReq
4. 在 `logic/base_logic.go` 添加 GetVersion 方法
5. 在 `controller/base_controller.go` 添加 GetVersion 方法
6. 在 `routes/base_routes.go` 注册路由
7. 测试验证

## 注意事项

1. 接口无需鉴权，方便前端和运维查询
2. Git commit 使用短哈希（7-8 位）
3. 构建时间使用 UTC 时间格式
4. 如果不在 Git 仓库中构建，版本信息会显示默认值（dev/unknown）
