# issue2md 开发任务列表

**Version**: 1.0
**Status**: Ready for Implementation
**Generated**: 2024-01-19

---

## 任务执行说明

### 符号说明
- `[P]` - 可与其他标记为`[P]`的任务并行执行（无依赖关系）
- `[TDD]` - 测试先行任务：必须先写测试（Red），再写实现（Green）
- `>>` - 任务依赖关系：左侧任务依赖右侧任务完成

### 执行原则
1. **严格遵守TDD循环**：每个功能模块必须先写测试文件，再写实现文件
2. **原子化执行**：每个任务只涉及一个文件的创建或修改
3. **依赖管理**：严格按照依赖关系顺序执行，并行任务可同时执行

---

## Phase 1: Foundation (数据结构定义)

> 目标：建立项目基础结构，定义核心数据类型，实现URL解析和配置管理。

### 1.1 项目初始化 [P]

```bash
# 任务列表
```

#### Task 1.1.1 [P]: 创建 go.mod
**文件**: `go.mod`
**内容**:
```go
module github.com/bigwhite/issue2md

go 1.18

require (
    github.com/google/go-github v65.0.0
    golang.org/x/oauth2 v0.15.0
)
```

#### Task 1.1.2 [P]: 创建 Makefile
**文件**: `Makefile`
**内容**: 定义 `test`, `build`, `web`, `clean` 等标准目标

#### Task 1.1.3 [P]: 创建目录结构
**创建目录**:
- `cmd/issue2md/`
- `internal/parser/`
- `internal/github/`
- `internal/converter/`
- `internal/config/`
- `internal/cli/`
- `internal/errors/`
- `test/fixtures/`

---

### 1.2 核心错误定义 [P]

#### Task 1.2.1 [TDD]: 创建 internal/errors/errors_test.go
**文件**: `internal/errors/errors_test.go`
**测试内容**:
- 验证预定义错误变量不为nil
- 测试 Wrap() 函数的正确性

#### Task 1.2.2: 创建 internal/errors/errors.go
**文件**: `internal/errors/errors.go`
**内容**:
```go
package errors

import (
    "errors"
    "fmt"
)

var (
    ErrInvalidURL         = errors.New("invalid GitHub URL format")
    ErrUnsupportedType    = errors.New("unsupported resource type")
    ErrResourceNotFound   = errors.New("resource not found")
    ErrAuthenticationFailed = errors.New("authentication failed")
    ErrRateLimitExceeded  = errors.New("rate limit exceeded")
    ErrNetworkError       = errors.New("network error")
    ErrEmptyData          = errors.New("empty issue data")
)

func Wrap(err error, message string) error {
    if err == nil {
        return nil
    }
    return fmt.Errorf("%s: %w", message, err)
}
```

---

### 1.3 Parser 包 - URL 解析

#### Task 1.3.1 [TDD]: 创建 internal/parser/types_test.go
**文件**: `internal/parser/types_test.go`
**测试内容**:
- 验证 ResourceType.String() 方法输出正确
- 测试 TypeIssue, TypePullRequest, TypeDiscussion, TypeUnknown

#### Task 1.3.2: 创建 internal/parser/types.go
**文件**: `internal/parser/types.go`
**内容**: ResourceType 类型定义、常量、String() 方法

#### Task 1.3.3 [TDD]: 创建 internal/parser/parser_test.go
**文件**: `internal/parser/parser_test.go`
**测试内容** (表格驱动):
```go
tests := []struct {
    name    string
    url     string
    want    *ParseResult
    wantErr error
}{
    { name: "valid issue url", ... },
    { name: "valid pr url", ... },
    { name: "valid discussion url", ... },
    { name: "invalid url format", ... },
    { name: "unsupported url type", ... },
    { name: "missing owner", ... },
    { name: "missing repo", ... },
    { name: "invalid number", ... },
}
```

#### Task 1.3.4: 创建 internal/parser/parser.go
**文件**: `internal/parser/parser.go`
**内容**:
- `ParseURL(urlStr string) (*ParseResult, error)` 函数
- URL 解析逻辑（使用 net/url 包）
- 类型识别逻辑

#### Task 1.3.5 [TDD]: 扩展 internal/parser/parser_test.go - Validate方法
**测试内容**: ParseResult.Validate() 方法的边界情况

#### Task 1.3.6: 修改 internal/parser/types.go
**添加**: `Validate() error` 方法到 ParseResult

---

### 1.4 Config 包 - 配置管理

#### Task 1.4.1 [TDD]: 创建 internal/config/config_test.go
**文件**: `internal/config/config_test.go`
**测试内容** (表格驱动):
```go
tests := []struct {
    name          string
    envToken      string
    wantToken     string
}{
    { name: "token from env", ... },
    { name: "no token set", ... },
}
```

#### Task 1.4.2: 创建 internal/config/config.go
**文件**: `internal/config/config.go`
**内容**:
- Config 结构体定义
- `LoadConfig() *Config` 函数（从环境变量读取 GITHUB_TOKEN）
- `MergeFlags(enableReactions, enableUserLinks bool, outputFile string)` 方法

---

## Phase 2: GitHub Fetcher (API交互逻辑)

> 目标：实现与 GitHub GraphQL API 的交互，获取 Issue/PR/Discussion 数据。

### 2.1 GitHub 包 - 数据结构定义

#### Task 2.1.1 [TDD]: 创建 internal/github/types_test.go
**文件**: `internal/github/types_test.go`
**测试内容**:
- 验证 ReactionCounts.IsEmpty() 方法
- 测试各字段类型正确性

#### Task 2.1.2: 创建 internal/github/types.go
**文件**: `internal/github/types.go`
**内容**:
```go
type User struct { ... }
type ReactionCounts struct { ...; IsEmpty() bool }
type Comment struct { ... }
type IssueData struct { ... }
```

---

### 2.2 GraphQL 查询定义

#### Task 2.2.1: 创建 internal/github/graphql.go
**文件**: `internal/github/graphql.go`
**内容**:
- issueQuery 常量（GraphQL 查询字符串）
- pullRequestQuery 常量
- discussionQuery 常量
- 变量结构体定义

---

### 2.3 GitHub 客户端实现

#### Task 2.3.1 [TDD]: 创建 internal/github/client_test.go - NewClient
**文件**: `internal/github/client_test.go`
**测试内容** (表格驱动):
```go
tests := []struct {
    name  string
    token string
    check func(*Client) bool
}{
    { name: "client with token", ... },
    { name: "client without token", ... },
}
```

#### Task 2.3.2: 创建 internal/github/client.go - NewClient
**文件**: `internal/github/client.go`
**内容**:
- Client 结构体定义
- `NewClient(token string) *Client` 函数

#### Task 2.3.3 [TDD]: 扩展 internal/github/client_test.go - buildQuery
**测试内容**: 验证不同资源类型返回正确的 GraphQL 查询

#### Task 2.3.4: 修改 internal/github/client.go - buildQuery
**添加**: `buildQuery(rt ResourceType) string` 私有方法

#### Task 2.3.5 [TDD]: 扩展 internal/github/client_test.go - parseResponse
**测试内容**: 验证 GraphQL 响应解析逻辑（使用 test fixtures）

#### Task 2.3.6 >> Task 2.3.7: 创建 test/fixtures/issue_response.json
**文件**: `test/fixtures/issue_response.json`
**内容**: 真实的 GitHub GraphQL API 响应示例

#### Task 2.3.7: 修改 internal/github/client.go - parseResponse
**添加**: `parseResponse(data []byte) (*IssueData, error)` 私有方法

#### Task 2.3.8 [TDD]: 扩展 internal/github/client_test.go - FetchIssueData
**测试内容** (集成测试，需要真实 GitHub token):
```go
tests := []struct {
    name    string
    url     string
    wantErr bool
}{
    { name: "fetch public issue", ... },
    { name: "fetch public pr", ... },
    { name: "not found resource", ... },
}
```

#### Task 2.3.9: 修改 internal/github/client.go - FetchIssueData
**添加**: `FetchIssueData(ctx context.Context, result *ParseResult) (*IssueData, error)` 方法

---

## Phase 3: Markdown Converter (转换逻辑)

> 目标：将获取的 GitHub 数据转换为格式化的 Markdown 输出。

### 3.1 Converter 包 - 数据结构定义

#### Task 3.1.1 [TDD]: 创建 internal/converter/types_test.go
**文件**: `internal/converter/types_test.go`
**测试内容**:
- 验证 Options 结构体默认值
- 测试选项组合

#### Task 3.1.2: 创建 internal/converter/types.go
**文件**: `internal/converter/types.go`
**内容**:
```go
type Options struct {
    EnableReactions bool
    EnableUserLinks bool
}
```

---

### 3.2 YAML Frontmatter 生成

#### Task 3.2.1 [TDD]: 创建 internal/converter/frontmatter_test.go
**文件**: `internal/converter/frontmatter_test.go`
**测试内容** (表格驱动):
```go
tests := []struct {
    name           string
    data           *IssueData
    enableReactions bool
    wantContains   []string // YAML 中必须包含的字段
}{
    { name: "basic frontmatter", ... },
    { name: "with reactions", ... },
    { name: "discussion type", ... },
}
```

#### Task 3.2.2: 创建 internal/converter/frontmatter.go
**文件**: `internal/converter/frontmatter.go`
**内容**:
- `generateFrontmatter(data *IssueData, enableReactions bool) string` 函数
- YAML 格式化逻辑

---

### 3.3 Markdown 正文生成

#### Task 3.3.1 [TDD]: 创建 internal/converter/body_test.go
**文件**: `internal/converter/body_test.go`
**测试内容** (表格驱动):
```go
tests := []struct {
    name          string
    data          *IssueData
    opts          *Options
    wantContains  []string
}{
    { name: "issue with comments", ... },
    { name: "pr with merged status", ... },
    { name: "discussion with answer", ... },
    { name: "with user links", ... },
}
```

#### Task 3.3.2: 创建 internal/converter/body.go
**文件**: `internal/converter/body.go`
**内容**:
- `generateBody(data *IssueData, opts *Options) string` 函数
- 标题行生成
- 元信息表格生成
- Description 格式化
- Comments 格式化（包括 Answer 标记）
- `@username` 到链接的转换

#### Task 3.3.3 [TDD]: 创建 internal/converter/converter_test.go
**文件**: `internal/converter/converter_test.go`
**测试内容**:
- 完整的 Convert() 函数集成测试
- 使用 test/fixtures/expected_output.md 验证输出

#### Task 3.3.4 >> Task 3.3.5: 创建 test/fixtures/expected_output.md
**文件**: `test/fixtures/expected_output.md`
**内容**: 完整的预期 Markdown 输出示例

#### Task 3.3.5: 创建 internal/converter/converter.go
**文件**: `internal/converter/converter.go`
**内容**:
- `Convert(data *github.IssueData, opts *Options) (string, error)` 函数
- 组合 frontmatter 和 body

---

### 3.4 文件写入

#### Task 3.4.1 [TDD]: 扩展 internal/converter/converter_test.go - WriteFile
**测试内容** (表格驱动):
```go
tests := []struct {
    name     string
    markdown string
    path     string
    wantErr  bool
}{
    { name: "write to stdout", ... },
    { name: "write to file", ... },
    { name: "invalid path", ... },
}
```

#### Task 3.4.2: 修改 internal/converter/converter.go - WriteFile
**添加**: `WriteFile(markdown, path string) error` 函数

---

## Phase 4: CLI Assembly (命令行入口集成)

> 目标：组装所有模块，实现完整的命令行工具。

### 4.1 CLI 包 - 命令行接口

#### Task 4.1.1 [TDD]: 创建 internal/cli/flags_test.go
**文件**: `internal/cli/flags_test.go`
**测试内容** (表格驱动):
```go
tests := []struct {
    name     string
    args     []string
    want     *Flags
    wantErr  bool
}{
    { name: "basic url only", ... },
    { name: "with output file", ... },
    { name: "with enable-reactions", ... },
    { name: "with enable-user-links", ... },
    { name: "help flag", ... },
    { name: "version flag", ... },
}
```

#### Task 4.1.2: 创建 internal/cli/flags.go
**文件**: `internal/cli/flags.go`
**内容**:
- Flags 结构体定义
- `ParseFlags(args []string) (*Flags, error)` 函数
- 使用标准库 flag 包

#### Task 4.1.3 [TDD]: 创建 internal/cli/cli_test.go
**文件**: `internal/cli/cli_test.go`
**测试内容**:
- 测试 Execute() 函数的完整流程
- 模拟不同场景（成功、失败）

#### Task 4.1.4: 创建 internal/cli/cli.go
**文件**: `internal/cli/cli.go`
**内容**:
- `Run(args []string) int` 函数
- `Execute(cfg *config.Config, args []string) error` 函数
- 错误处理和退出码管理

---

### 4.2 入口点实现

#### Task 4.2.1: 创建 cmd/issue2md/main.go
**文件**: `cmd/issue2md/main.go`
**内容**:
```go
package main

import (
    "os"
    "github.com/bigwhite/issue2md/internal/cli"
)

func main() {
    os.Exit(cli.Run(os.Args[1:]))
}
```

---

## Phase 5: 完善与验收

### 5.1 文档

#### Task 5.1.1 [P]: 创建 README.md
**文件**: `README.md`
**内容**:
- 项目简介
- 安装说明
- 使用示例
- 环境变量说明
- Badge（build status, coverage等）

#### Task 5.1.2 [P]: 创建 CONTRIBUTING.md
**文件**: `CONTRIBUTING.md`
**内容**: 贡献指南

#### Task 5.1.3 [P]: 更新 go.mod 元数据
**添加**: 包文档注释

---

### 5.2 版本信息

#### Task 5.2.1: 创建 internal/version/version.go
**文件**: `internal/version/version.go`
**内容**:
```go
package version

var (
    Version = "dev"
    Commit  = "none"
    Date    = "unknown"
)
```

#### Task 5.2.2: 修改 internal/cli/cli.go - 添加版本支持
**添加**: `-version` flag 的处理逻辑

---

### 5.3 构建配置

#### Task 5.3.1: 创建 .goreleaser.yml
**文件**: `.goreleaser.yml`
**内容**: GoReleaser 配置（用于发布）

#### Task 5.3.2 [P]: 创建 .gitignore
**文件**: `.gitignore`
**内容**: 标准 Go 项目忽略文件

#### Task 5.3.3 [P]: 创建 GitHub Actions workflow
**文件**: `.github/workflows/test.yml`
**内容**: CI/CD 配置

---

### 5.4 最终验收

#### Task 5.4.1: 运行完整测试套件
**命令**: `make test`
**验证**: 所有测试通过

#### Task 5.4.2: 检查测试覆盖率
**命令**: `go test -cover ./...`
**目标**: >= 80%

#### Task 5.4.3: 运行 go vet
**命令**: `go vet ./...`
**验证**: 无警告

#### Task 5.4.4: 手动 E2E 测试
**命令**:
```bash
go run cmd/issue2md/main.go https://github.com/golang/go/issues/1
go run cmd/issue2md/main.go -enable-reactions https://github.com/golang/go/issues/1 output.md
```
**验证**: 输出格式正确

---

## 依赖关系图

```
Phase 1 (Foundation)
├── 1.1 项目初始化 [P]
├── 1.2 错误定义 [P]
├── 1.3 Parser 包 (1.3.1→1.3.2→1.3.3→1.3.4→1.3.5→1.3.6)
└── 1.4 Config 包 (1.4.1→1.4.2)

Phase 2 (GitHub Fetcher)
├── 2.1 数据结构 (2.1.1→2.1.2) >> 依赖 Phase 1
├── 2.2 GraphQL 查询 (2.2.1) >> 依赖 2.1
└── 2.3 客户端实现 (2.3.1→2.3.2→2.3.3→2.3.4→2.3.5→2.3.6→2.3.7→2.3.8→2.3.9)
    └── >> 依赖 2.2, 1.3

Phase 3 (Markdown Converter)
├── 3.1 数据结构 (3.1.1→3.1.2) [P]
├── 3.2 Frontmatter (3.2.1→3.2.2) >> 依赖 2.1
├── 3.3 Body生成 (3.3.1→3.3.2→3.3.3→3.3.4→3.3.5) >> 依赖 3.2
└── 3.4 文件写入 (3.4.1→3.4.2) >> 依赖 3.3

Phase 4 (CLI Assembly)
├── 4.1 CLI包 (4.1.1→4.1.2→4.1.3→4.1.4) >> 依赖 Phase 1,2,3
└── 4.2 入口点 (4.2.1) >> 依赖 4.1

Phase 5 (完善与验收)
├── 5.1 文档 [P]
├── 5.2 版本 (5.2.1→5.2.2)
├── 5.3 构建 [P]
└── 5.4 验收 >> 依赖所有之前阶段
```

---

## 任务统计

| Phase | 任务数 | 测试任务 | 实现任务 | 并行任务 |
|-------|--------|----------|----------|----------|
| Phase 1 | 14 | 6 | 8 | 3 |
| Phase 2 | 9 | 5 | 4 | 0 |
| Phase 3 | 11 | 6 | 5 | 1 |
| Phase 4 | 4 | 2 | 2 | 0 |
| Phase 5 | 10 | 0 | 10 | 3 |
| **总计** | **48** | **19** | **29** | **7** |

---

**生成时间**: 2024-01-19
**遵循规范**: constitution.md v1.0
