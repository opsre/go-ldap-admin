# 版本信息接口实施计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 添加版本信息返回接口，支持查询系统版本号、Git 提交哈希、构建时间、Go 版本和编译平台

**Architecture:** 采用编译时注入方式，通过 go build -ldflags 将版本信息注入到 public/version 包的全局变量中。API 层通过 Controller → Logic 模式读取版本信息并返回 JSON 响应。

**Tech Stack:** Go 1.21+, Gin Web Framework, Go build ldflags, Makefile

---

## Task 1: 创建版本信息包

**Files:**
- Create: `.worktrees/version-api/public/version/version.go`

**Step 1: 创建 version 包文件**

创建 `public/version/version.go` 文件，定义版本信息全局变量和获取函数：

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

**Step 2: 验证包可以编译**

Run: `go build ./public/version`
Expected: 编译成功，无错误输出

**Step 3: 提交代码**

```bash
git add public/version/version.go
git commit -m "feat(version): 添加版本信息包"
```

---

## Task 2: 添加 Request 结构体

**Files:**
- Modify: `.worktrees/version-api/model/request/base_req.go:31`

**Step 1: 添加 BaseVersionReq 结构体**

在 `model/request/base_req.go` 文件末尾添加：

```go
// BaseVersionReq 获取版本信息结构体
type BaseVersionReq struct {
}
```

**Step 2: 验证编译**

Run: `go build ./model/request`
Expected: 编译成功

**Step 3: 提交代码**

```bash
git add model/request/base_req.go
git commit -m "feat(request): 添加版本信息请求结构体"
```

---

## Task 3: 添加 Logic 层实现

**Files:**
- Modify: `.worktrees/version-api/logic/base_logic.go:3-12` (import section)
- Modify: `.worktrees/version-api/logic/base_logic.go:219` (end of file)

**Step 1: 添加 version 包导入**

在 `logic/base_logic.go` 的 import 部分添加：

```go
import (
	"fmt"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/model/request"
	"github.com/eryajf/go-ldap-admin/model/response"
	"github.com/eryajf/go-ldap-admin/public/tools"
	"github.com/eryajf/go-ldap-admin/public/version"
	"github.com/eryajf/go-ldap-admin/service/ildap"
	"github.com/eryajf/go-ldap-admin/service/isql"

	"github.com/gin-gonic/gin"
)
```

**Step 2: 添加 GetVersion 方法**

在 `logic/base_logic.go` 文件末尾添加：

```go
// GetVersion 获取版本信息
func (l BaseLogic) GetVersion(c *gin.Context, req any) (data any, rspError any) {
	_, ok := req.(*request.BaseVersionReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	return version.GetVersion(), nil
}
```

**Step 3: 验证编译**

Run: `go build ./logic`
Expected: 编译成功

**Step 4: 提交代码**

```bash
git add logic/base_logic.go
git commit -m "feat(logic): 添加版本信息获取逻辑"
```

---

## Task 4: 添加 Controller 层实现

**Files:**
- Modify: `.worktrees/version-api/controller/base_controller.go:105` (end of file)

**Step 1: 添加 GetVersion 方法**

在 `controller/base_controller.go` 文件末尾添加：

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

**Step 2: 验证编译**

Run: `go build ./controller`
Expected: 编译成功

**Step 3: 提交代码**

```bash
git add controller/base_controller.go
git commit -m "feat(controller): 添加版本信息接口处理器"
```

---

## Task 5: 注册路由

**Files:**
- Modify: `.worktrees/version-api/routes/base_routes.go:51`

**Step 1: 添加版本信息路由**

在 `routes/base_routes.go` 的 `InitBaseRoutes` 函数中，在 `base.GET("config", controller.Base.GetConfig)` 这行之后添加：

```go
		base.GET("config", controller.Base.GetConfig)         // 获取系统配置
		base.GET("version", controller.Base.GetVersion)       // 获取版本信息
		// 登录登出刷新token无需鉴权
```

**Step 2: 验证编译**

Run: `go build ./routes`
Expected: 编译成功

**Step 3: 提交代码**

```bash
git add routes/base_routes.go
git commit -m "feat(routes): 注册版本信息接口路由"
```

---

## Task 6: 修改 Makefile 添加版本信息注入

**Files:**
- Modify: `.worktrees/version-api/Makefile:1-17`

**Step 1: 在 Makefile 顶部添加版本信息变量**

在 `Makefile` 的 `OUTPUT_DIR` 定义之后添加版本信息变量和 ldflags：

```makefile
# 定义项目名称
BINARY_NAME=go-ldap-admin

# 定义输出目录
OUTPUT_DIR=bin

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
```

**Step 2: 修改 build 命令**

将原来的 `build` 目标替换为：

```makefile
.PHONY: build
build:
	go build -ldflags "$(LDFLAGS) -X 'github.com/eryajf/go-ldap-admin/public/version.Platform=$(shell go env GOOS)/$(shell go env GOARCH)'" -o ${BINARY_NAME} main.go
```

**Step 3: 修改 build-linux 命令**

将原来的 `build-linux` 目标替换为：

```makefile
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "$(LDFLAGS) -X 'github.com/eryajf/go-ldap-admin/public/version.Platform=linux/amd64'" -o ${BINARY_NAME} main.go
```

**Step 4: 修改 gox-linux 命令**

将原来的 `gox-linux` 目标替换为：

```makefile
.PHONY: gox-linux
gox-linux:
	CGO_ENABLED=0 gox -osarch="linux/amd64 linux/arm64" \
		-ldflags "$(LDFLAGS) -X 'github.com/eryajf/go-ldap-admin/public/version.Platform={{.OS}}/{{.Arch}}'" \
		-output="${OUTPUT_DIR}/${BINARY_NAME}_{{.OS}}_{{.Arch}}"
	@for b in $$(ls ${OUTPUT_DIR}); do \
		OUTPUT_FILE="${OUTPUT_DIR}/$$b"; \
		upx -9 -q "$$OUTPUT_FILE"; \
	done
```

**Step 5: 修改 gox-all 命令**

将原来的 `gox-all` 目标替换为：

```makefile
.PHONY: gox-all
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

**Step 6: 验证 Makefile 语法**

Run: `make -n build`
Expected: 显示将要执行的命令，无语法错误

**Step 7: 提交代码**

```bash
git add Makefile
git commit -m "feat(build): 添加版本信息编译时注入"
```

---

## Task 7: 测试版本信息接口

**Step 1: 构建项目**

Run: `make build`
Expected: 编译成功，生成 go-ldap-admin 二进制文件

**Step 2: 启动服务（如果配置允许）**

Run: `./go-ldap-admin` (或 `make run`)
Expected: 服务启动成功

**Step 3: 测试版本信息接口**

Run: `curl http://localhost:8888/api/base/version`
Expected: 返回 JSON 响应，包含版本信息：
```json
{
    "code": 200,
    "data": {
        "version": "v1.5.0-xxx",
        "gitCommit": "e478855",
        "buildTime": "2026-02-06 09:02:15",
        "goVersion": "go1.25.4",
        "platform": "darwin/arm64"
    },
    "msg": "success"
}
```

**Step 4: 验证版本信息正确性**

- 检查 gitCommit 是否匹配当前 commit: `git rev-parse --short HEAD`
- 检查 goVersion 是否匹配: `go version`
- 检查 platform 是否匹配: `go env GOOS` 和 `go env GOARCH`

**Step 5: 最终提交**

如果测试通过，创建最终提交：

```bash
git add -A
git commit -m "test: 验证版本信息接口功能正常"
```

---

## Task 8: 更新文档

**Files:**
- Modify: `.worktrees/version-api/docs/plans/2026-02-06-version-info-design.md`

**Step 1: 在设计文档中添加实施完成标记**

在设计文档顶部添加：

```markdown
# 版本信息接口设计文档

**状态:** ✅ 已实施完成

**实施日期:** 2026-02-06

**实施分支:** feature/version-api
```

**Step 2: 提交文档更新**

```bash
git add docs/plans/2026-02-06-version-info-design.md
git commit -m "docs: 更新版本信息接口设计文档状态"
```

---

## 验收标准

1. ✅ `public/version/version.go` 包创建成功，包含 5 个版本信息变量
2. ✅ Request、Logic、Controller 层代码添加完成
3. ✅ 路由注册成功，接口路径为 `/api/base/version`
4. ✅ Makefile 修改完成，所有构建命令支持版本信息注入
5. ✅ 本地构建测试通过，版本信息正确返回
6. ✅ API 响应格式符合项目规范
7. ✅ 所有代码已提交到 feature/version-api 分支

## 注意事项

1. 接口无需鉴权，任何人都可以访问
2. 如果不在 Git 仓库中构建，版本信息会显示默认值（dev/unknown）
3. 构建时间使用 UTC 时间格式
4. Git commit 使用短哈希（7-8 位）
5. 所有构建方式（本地、GitHub Actions、CNB）都会自动注入版本信息

## 后续工作

1. 可以考虑在前端页面底部显示版本信息
2. 可以在系统日志中记录启动时的版本信息
3. 可以添加版本信息到 Swagger 文档中
