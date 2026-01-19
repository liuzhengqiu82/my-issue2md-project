# issue2md 技术实现方案

**Version**: 1.0
**Status**: Ready for Implementation
**Author**: AI Architect Agent
**Date**: 2024-01-19

---

## 1. 技术上下文总结

### 1.1 技术栈选型

| 组件 | 技术选型 | 理由 |
|------|---------|------|
| **语言** | Go 1.18+ | 项目宪法明确要求，高性能、并发友好 |
| **Web框架** | `net/http` (标准库) | 遵循"简单性原则"，不引入Gin/Echo等外部框架 |
| **GitHub API** | `google/go-github` v65+ | 官方维护，支持GraphQL v4 API |
| **GraphQL** | 内嵌查询字符串 | 避免引入额外GraphQL库，保持简单 |
| **Markdown输出** | 标准库 `fmt` + `strings` | 不使用第三方模板引擎 |
| **数据存储** | 无（实时API获取） | MVP阶段无需持久化 |
| **测试框架** | `testing` (标准库) | 表格驱动测试，拒绝Mock |

### 1.2 外部依赖

```go
// go.mod 最终依赖列表
module github.com/bigwhite/issue2md

go 1.18

require (
    github.com/google/go-github v65.0.0
    golang.org/x/oauth2 v0.15.0 // 仅用于GitHub认证
)
```

---

## 2. "合宪性"审查

### 2.1 对照宪法条款逐条审查

| 宪法条款 | 审查结果 | 说明 |
|---------|---------|------|
| **1.1 YAGNI** | ✅ 合规 | 仅实现spec.md明确要求的功能，无超前设计 |
| **1.2 标准库优先** | ✅ 合规 | Web层使用`net/http`，Markdown输出使用`fmt`/`strings` |
| **1.3 反过度工程** | ✅ 合规 | 无接口抽象，直接使用结构体，无设计模式 |
| **2.1 TDD循环** | ✅ 合规 | 每个模块开发前先写失败的表格驱动测试 |
| **2.2 表格驱动测试** | ✅ 合规 | 所有单元测试采用`[]struct{...}{...}`模式 |
| **2.3 拒绝Mocks** | ✅ 合规 | 优先集成测试，使用真实GitHub API（或testcontainers） |
| **3.1 错误处理** | ✅ 合规 | 所有错误显式处理，使用`fmt.Errorf("...: %w", err)`包装 |
| **3.2 无全局变量** | ✅ 合规 | 所有依赖通过构造函数/参数显式注入 |

### 2.2 违规风险评估

| 风险点 | 评估 | 缓解措施 |
|-------|------|---------|
| GitHub API限流 | 低 | 支持`GITHUB_TOKEN`环境变量，实现重试逻辑 |
| 私有仓库访问 | 低 | 通过Token认证，错误信息明确提示 |
| GraphQL查询复杂度 | 中 | 预定义查询字符串，静态验证 |

---

## 3. 项目结构细化

### 3.1 完整目录结构

```
issue2md/
├── cmd/
│   └── issue2md/
│       └── main.go              # 入口点，调用cli.Run()
│
├── internal/
│   ├── parser/                  # URL解析与类型识别
│   │   ├── parser.go            # ParseURL(), ResourceType定义
│   │   └── parser_test.go       # 表格驱动测试
│   │
│   ├── github/                  # GitHub API交互
│   │   ├── client.go            # Client结构体，FetchIssueData()
│   │   ├── graphql.go           # GraphQL查询字符串
│   │   ├── types.go             # IssueData, Comment, User等
│   │   └── client_test.go       # 集成测试
│   │
│   ├── converter/               # 数据转换为Markdown
│   │   ├── converter.go         # Convert(), WriteFile()
│   │   ├── frontmatter.go       # YAML frontmatter生成
│   │   └── converter_test.go    # 表格驱动测试
│   │
│   ├── config/                  # 配置管理
│   │   ├── config.go            # Config结构体，LoadConfig()
│   │   └── config_test.go       # 表格驱动测试
│   │
│   └── cli/                     # 命令行接口
│       ├── cli.go               # Run(), Execute()
│       ├── flags.go             # Flag解析
│       └── cli_test.go          # 表格驱动测试
│
├── test/
│   └── fixtures/                # 测试固件
│       ├── issue_response.json  # GitHub API响应示例
│       └── expected_output.md   # 预期Markdown输出
│
├── specs/
│   ├── spec.md                  # 功能规格（已存在）
│   ├── plan.md                  # 本文档
│   └── api-sketch.md            # API草案（已存在）
│
├── constitution.md              # 项目宪法（已存在）
├── CLAUDE.md                    # AI协作指南（已存在）
├── go.mod
├── go.sum
├── Makefile                     # 标准化操作
└── README.md
```

### 3.2 包依赖关系

```
┌─────────────────────────────────────────────────────────┐
│                      cmd/issue2md                       │
│                        main.go                          │
└───────────────────────────┬─────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│                      internal/cli                       │
│  ┌─────────────────────────────────────────────────┐   │
│  │  Run() ──► config.LoadConfig()                  │   │
│  │          ──► parser.ParseURL()                   │   │
│  │          ──► github.NewClient()                  │   │
│  │          ──► client.FetchIssueData()             │   │
│  │          ──► converter.Convert()                 │   │
│  │          ──► converter.WriteFile()               │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        ▼                   ▼                   ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│internal/     │   │internal/     │   │internal/     │
│  parser      │   │  github      │   │  converter   │
│              │   │              │   │              │
│ ParseURL()   │   │ NewClient()  │   │ Convert()    │
│ ParseResult  │   │ FetchIssue() │   │ WriteFile()  │
└──────────────┘   └──────────────┘   └──────────────┘
                            ▲
                            │
                    ┌──────────────┐
                    │internal/     │
                    │  config      │
                    │              │
                    │ LoadConfig() │
                    └──────────────┘
```

### 3.3 包职责详解

| 包名 | 职责 | 依赖 |
|------|------|------|
| `parser` | URL解析、类型识别、验证 | 无 |
| `github` | GraphQL API调用、数据获取 | `parser`, `config` |
| `converter` | Markdown格式化、YAML生成 | `github` (IssueData) |
| `config` | 环境变量读取、配置结构 | 无 |
| `cli` | 命令行解析、流程编排 | 以上所有包 |

---

## 4. 核心数据结构

### 4.1 跨包数据流

```
ParseResult ──► IssueData ──► Markdown
   (parser)      (github)      (converter)
```

### 4.2 完整数据结构定义

```go
// ============================================================================
// internal/parser/types.go
// ============================================================================

// ResourceType 表示 GitHub 资源类型
type ResourceType int

const (
    TypeUnknown ResourceType = iota
    TypeIssue
    TypePullRequest
    TypeDiscussion
)

func (t ResourceType) String() string {
    switch t {
    case TypeIssue:
        return "issue"
    case TypePullRequest:
        return "pr"
    case TypeDiscussion:
        return "discussion"
    default:
        return "unknown"
    }
}

// ParseResult 表示 URL 解析结果
type ParseResult struct {
    Type   ResourceType
    Owner  string
    Repo   string
    Number int
    URL    string // 原始URL
}

// ============================================================================
// internal/github/types.go
// ============================================================================

// User 表示 GitHub 用户信息
type User struct {
    Login     string
    AvatarURL string
    URL       string // GitHub 主页链接
}

// ReactionCounts 表示 reactions 统计
type ReactionCounts struct {
    ThumbsUp   int `json:"THUMBS_UP"`
    ThumbsDown int `json:"THUMBS_DOWN"`
    Laugh      int `json:"LAUGH"`
    Hooray     int `json:"HOORAY"`
    Confused   int `json:"CONFUSED"`
    Heart      int `json:"HEART"`
    Rocket     int `json:"ROCKET"`
    Eyes       int `json:"EYES"`
}

// IsEmpty 判断是否有任何reaction
func (r *ReactionCounts) IsEmpty() bool {
    return r.ThumbsUp == 0 && r.ThumbsDown == 0 &&
        r.Laugh == 0 && r.Hooray == 0 &&
        r.Confused == 0 && r.Heart == 0 &&
        r.Rocket == 0 && r.Eyes == 0
}

// Comment 表示一条评论
type Comment struct {
    Author    *User
    Body      string
    CreatedAt time.Time
    UpdatedAt time.Time
    Reactions *ReactionCounts
    IsAnswer  bool // 仅 Discussion 使用：标记为答案
}

// IssueData 表示 Issue/PR/Discussion 的完整数据
// 这是跨包传递的核心数据结构
type IssueData struct {
    Type         ResourceType
    Title        string
    URL          string
    Author       *User
    CreatedAt    time.Time
    UpdatedAt    time.Time
    Status       string // "open", "closed", "merged", "answered"
    Body         string
    Comments     []*Comment
    Reactions    *ReactionCounts
    TotalCommentCount int // 包括评论和review评论总数
}

// ============================================================================
// internal/converter/types.go
// ============================================================================

// Options 控制输出选项
type Options struct {
    EnableReactions bool // 包含 reactions 统计
    EnableUserLinks bool // @username 转换为链接
}

// ============================================================================
// internal/config/types.go
// ============================================================================

// Config 表示应用配置
type Config struct {
    GitHubToken     string
    EnableReactions bool
    EnableUserLinks bool
    OutputFile      string // 空字符串表示 stdout
}
```

---

## 5. 接口设计

### 5.1 `internal/parser` 接口

```go
// ParseURL 解析 GitHub URL 并返回资源信息
//
// 参数:
//   - urlStr: 待解析的URL字符串
//
// 返回:
//   - *ParseResult: 解析结果，包含类型、仓库信息、编号
//   - error: 当URL格式无效或不支持的类型时返回错误
//
// 支持的URL格式:
//   - https://github.com/{owner}/{repo}/issues/{number}
//   - https://github.com/{owner}/{repo}/pull/{number}
//   - https://github.com/{owner}/{repo}/discussions/{number}
//
// 示例:
//   result, err := parser.ParseURL("https://github.com/golang/go/issues/12345")
func ParseURL(urlStr string) (*ParseResult, error)

// Validate 验证 ParseResult 是否有效
func (r *ParseResult) Validate() error
```

### 5.2 `internal/github` 接口

```go
// NewClient 创建新的 GitHub API 客户端
//
// 参数:
//   - token: GitHub Personal Access Token，空字符串时使用匿名访问
//
// 注意:
//   - 匿名访问有严格的API限流
//   - 推荐通过 GITHUB_TOKEN 环境变量设置
func NewClient(token string) *Client

// FetchIssueData 获取指定 Issue/PR/Discussion 的完整数据
//
// 参数:
//   - ctx: 上下文，用于超时控制
//   - result: ParseResult，包含资源类型和位置信息
//
// 返回:
//   - *IssueData: 完整的数据，包括标题、作者、评论等
//   - error: API错误、网络错误、数据解析错误
//
// 根据 ParseResult.Type 自动选择对应的 GraphQL 查询:
//   - TypeIssue: 使用 issueQuery
//   - TypePullRequest: 使用 pullRequestQuery
//   - TypeDiscussion: 使用 discussionQuery
func (c *Client) FetchIssueData(ctx context.Context, result *ParseResult) (*IssueData, error)
```

### 5.3 `internal/converter` 接口

```go
// Convert 将 IssueData 转换为 Markdown 格式
//
// 参数:
//   - data: 从 GitHub API 获取的完整数据
//   - opts: 输出选项控制
//
// 返回:
//   - string: 完整的 Markdown 字符串，包含 YAML frontmatter 和正文
//   - error: 模板渲染错误（理论上不应发生）
//
// 输出格式:
//   1. YAML frontmatter (元数据)
//   2. 标题行（包含状态）
//   3. 元信息表格（作者、时间、状态等）
//   4. Description (原始body)
//   5. Comments (按时间排序)
func Convert(data *github.IssueData, opts *Options) (string, error)

// WriteFile 将 Markdown 写入文件或 stdout
//
// 参数:
//   - markdown: Markdown 内容
//   - path: 文件路径，空字符串表示写入 stdout
//
// 返回:
//   - error: 文件写入失败时返回错误
func WriteFile(markdown, path string) error
```

### 5.4 `internal/config` 接口

```go
// LoadConfig 从环境变量和命令行参数加载配置
//
// 环境变量:
//   - GITHUB_TOKEN: GitHub Personal Access Token (可选)
//
// 命令行参数由 cli 包解析后传入
func LoadConfig() *Config

// MergeFlags 从命令行参数更新配置
func (c *Config) MergeFlags(enableReactions, enableUserLinks bool, outputFile string)
```

### 5.5 `internal/cli` 接口

```go
// Run 执行 CLI 主逻辑
//
// 参数:
//   - args: 命令行参数（不含程序名）
//
// 返回:
//   - int: 进程退出码 (0=成功, 1=错误)
//
// 命令格式:
//   issue2md [flags] <url> [output_file]
//
// Flags:
//   -enable-reactions    # 包含 reactions 统计
//   -enable-user-links   # 用户名渲染为链接
//   -h, -help            # 显示帮助
//   -v, -version         # 显示版本
func Run(args []string) int

// Execute 是实际的执行函数，便于测试
func Execute(cfg *config.Config, args []string) error
```

---

## 6. GraphQL查询设计

### 6.1 Issue 查询

```graphql
query GetIssue($owner: String!, $repo: String!, $number: Int!) {
    repository(owner: $owner, name: $repo) {
        issue(number: $number) {
            title
            url
            state
            createdAt
            updatedAt
            bodyText
            author {
                login
                avatarUrl
                url
            }
            reactions(first: 100) {
                thumbsUp: reactionCount(content: THUMBS_UP)
                thumbsDown: reactionCount(content: THUMBS_DOWN)
                laugh: reactionCount(content: LAUGH)
                hooray: reactionCount(content: HOORAY)
                confused: reactionCount(content: CONFUSED)
                heart: reactionCount(content: HEART)
                rocket: reactionCount(content: ROCKET)
                eyes: reactionCount(content: EYES)
            }
            comments(first: 100) {
                totalCount
                nodes {
                    author { login avatarUrl url }
                    bodyText
                    createdAt
                    updatedAt
                    reactions(first: 100) {
                        thumbsUp: reactionCount(content: THUMBS_UP)
                        thumbsDown: reactionCount(content: THUMBS_DOWN)
                        laugh: reactionCount(content: LAUGH)
                        hooray: reactionCount(content: HOORAY)
                        confused: reactionCount(content: CONFUSED)
                        heart: reactionCount(content: HEART)
                        rocket: reactionCount(content: ROCKET)
                        eyes: reactionCount(content: EYES)
                    }
                }
            }
        }
    }
}
```

### 6.2 PullRequest 查询

```graphql
query GetPullRequest($owner: String!, $repo: String!, $number: Int!) {
    repository(owner: $owner, name: $repo) {
        pullRequest(number: $number) {
            title
            url
            state
            merged
            createdAt
            updatedAt
            bodyText
            author {
                login
                avatarUrl
                url
            }
            reactions(first: 100) {
                thumbsUp: reactionCount(content: THUMBS_UP)
                thumbsDown: reactionCount(content: THUMBS_DOWN)
                laugh: reactionCount(content: LAUGH)
                hooray: reactionCount(content: HOORAY)
                confused: reactionCount(content: CONFUSED)
                heart: reactionCount(content: HEART)
                rocket: reactionCount(content: ROCKET)
                eyes: reactionCount(content: EYES)
            }
            comments(first: 100) {
                totalCount
                nodes {
                    author { login avatarUrl url }
                    bodyText
                    createdAt
                    updatedAt
                    reactions(first: 100) {
                        thumbsUp: reactionCount(content: THUMBS_UP)
                        thumbsDown: reactionCount(content: THUMBS_DOWN)
                        laugh: reactionCount(content: LAUGH)
                        hooray: reactionCount(content: HOORAY)
                        confused: reactionCount(content: CONFUSED)
                        heart: reactionCount(content: HEART)
                        rocket: reactionCount(content: ROCKET)
                        eyes: reactionCount(content: EYES)
                    }
                }
            }
            reviews(first: 100) {
                totalCount
                nodes {
                    author { login avatarUrl url }
                    bodyText
                    createdAt
                    updatedAt
                    state
                    comments(first: 100) {
                        nodes {
                            author { login avatarUrl url }
                            bodyText
                            createdAt
                            updatedAt
                        }
                    }
                }
            }
        }
    }
}
```

### 6.3 Discussion 查询

```graphql
query GetDiscussion($owner: String!, $repo: String!, $number: Int!) {
    repository(owner: $owner, name: $repo) {
        discussion(number: $number) {
            title
            url
            state
            isAnswered
            createdAt
            updatedAt
            bodyText
            author {
                login
                avatarUrl
                url
            }
            reactions(first: 100) {
                thumbsUp: reactionCount(content: THUMBS_UP)
                thumbsDown: reactionCount(content: THUMBS_DOWN)
                laugh: reactionCount(content: LAUGH)
                hooray: reactionCount(content: HOORAY)
                confused: reactionCount(content: CONFUSED)
                heart: reactionCount(content: HEART)
                rocket: reactionCount(content: ROCKET)
                eyes: reactionCount(content: EYES)
            }
            comments(first: 100) {
                totalCount
                nodes {
                    author { login avatarUrl url }
                    bodyText
                    createdAt
                    updatedAt
                    isAnswer
                    reactions(first: 100) {
                        thumbsUp: reactionCount(content: THUMBS_UP)
                        thumbsDown: reactionCount(content: THUMBS_DOWN)
                        laugh: reactionCount(content: LAUGH)
                        hooray: reactionCount(content: HOORAY)
                        confused: reactionCount(content: CONFUSED)
                        heart: reactionCount(content: HEART)
                        rocket: reactionCount(content: ROCKET)
                        eyes: reactionCount(content: EYES)
                    }
                }
            }
        }
    }
}
```

---

## 7. 错误处理策略

### 7.1 错误类型定义

```go
// internal/errors/errors.go
package errors

import (
    "errors"
    "fmt"
)

var (
    // Parser 错误
    ErrInvalidURL    = errors.New("invalid GitHub URL format")
    ErrUnsupportedType = errors.New("unsupported resource type")

    // GitHub 错误
    ErrResourceNotFound = errors.New("resource not found")
    ErrAuthenticationFailed = errors.New("authentication failed")
    ErrRateLimitExceeded = errors.New("rate limit exceeded")
    ErrNetworkError = errors.New("network error")

    // Converter 错误
    ErrEmptyData = errors.New("empty issue data")

    // Config 错误
    ErrInvalidOutputPath = errors.New("invalid output file path")
)

// Wrap 包装错误并添加上下文
func Wrap(err error, message string) error {
    if err == nil {
        return nil
    }
    return fmt.Errorf("%s: %w", message, err)
}
```

### 7.2 退出码约定

| 退出码 | 含义 | 示例场景 |
|-------|------|---------|
| 0 | 成功 | 正常完成转换 |
| 1 | 一般错误 | 无效URL、网络错误 |
| 2 | API错误 | 资源不存在、认证失败 |
| 3 | 文件错误 | 无法写入输出文件 |

---

## 8. 测试策略

### 8.1 测试金字塔

```
        ┌─────────────┐
        │  E2E Tests  │  10% - 完整流程测试
        ├─────────────┤
        │ Integration │  30% - API交互测试
        │    Tests    │
        ├─────────────┤
        │  Unit Tests │  60% - 表格驱动测试
        └─────────────┘
```

### 8.2 表格驱动测试模板

```go
// internal/parser/parser_test.go

func TestParseURL(t *testing.T) {
    tests := []struct {
        name    string
        url     string
        want    *ParseResult
        wantErr error
    }{
        {
            name: "valid issue url",
            url:  "https://github.com/golang/go/issues/12345",
            want: &ParseResult{
                Type:   TypeIssue,
                Owner:  "golang",
                Repo:   "go",
                Number: 12345,
                URL:    "https://github.com/golang/go/issues/12345",
            },
            wantErr: nil,
        },
        {
            name:    "invalid url",
            url:     "not-a-url",
            want:    nil,
            wantErr: ErrInvalidURL,
        },
        // ... 更多测试用例
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseURL(tt.url)
            if tt.wantErr != nil {
                if !errors.Is(err, tt.wantErr) {
                    t.Errorf("ParseURL() error = %v, wantErr %v", err, tt.wantErr)
                }
                return
            }
            if err != nil {
                t.Fatalf("ParseURL() unexpected error: %v", err)
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ParseURL() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 8.3 测试覆盖目标

| 包 | 目标覆盖率 | 重点测试内容 |
|---|-----------|-------------|
| parser | 90%+ | URL解析边界情况 |
| github | 80%+ | API响应解析、错误处理 |
| converter | 90%+ | Markdown格式、YAML生成 |
| config | 80%+ | 环境变量读取 |
| cli | 70%+ | 命令行解析、流程编排 |

---

## 9. 实现路线图

### Phase 1: 核心基础 (Week 1)
- [ ] 搭建项目结构
- [ ] 实现 `internal/parser` 包
- [ ] 实现 `internal/config` 包
- [ ] 编写 parser 和 config 的表格驱动测试

### Phase 2: GitHub集成 (Week 2)
- [ ] 实现 `internal/github` 包
- [ ] 添加 `google/go-github` 依赖
- [ ] 定义 GraphQL 查询
- [ ] 编写集成测试（使用真实API）

### Phase 3: Markdown转换 (Week 3)
- [ ] 实现 `internal/converter` 包
- [ ] YAML frontmatter 生成
- [ ] 评论格式化和排序
- [ ] Reactions 和用户链接支持

### Phase 4: CLI和集成 (Week 4)
- [ ] 实现 `internal/cli` 包
- [ ] 实现入口 `cmd/issue2md/main.go`
- [ ] 端到端测试
- [ ] 文档编写

---

## 10. 验收检查清单

### 功能验收
- [ ] 支持三种URL类型解析
- [ ] 正确获取Issue/PR/Discussion数据
- [ ] Markdown输出符合spec格式
- [ ] YAML frontmatter完整
- [ ] Reactions可选输出
- [ ] 用户链接可选输出
- [ ] 支持stdout和文件输出
- [ ] GITHUB_TOKEN环境变量支持

### 质量验收
- [ ] 所有测试通过 (`make test`)
- [ ] 测试覆盖率 >= 80%
- [ ] 无`go vet`警告
- [ ] 无`golangci-lint`问题
- [ ] 符合宪法所有条款

### 文档验收
- [ ] README.md 使用说明
- [ ] Godoc注释完整
- [ ] 示例输出文档

---

**文档版本历史:**
- v1.0 (2024-01-19): 初始版本，基于spec.md和constitution.md创建
