# issue2md API 接口草案

> 本文档描述 `internal` 目录下各包对外暴露的主要接口。
> 设计原则：简单、明确、避免过度抽象。

---

## 1. `internal/parser` - URL 解析

### 1.1 数据结构

```go
// ResourceType 表示 GitHub 资源类型
type ResourceType int

const (
    TypeUnknown ResourceType = iota
    TypeIssue
    TypePullRequest
    TypeDiscussion
)

// ParseResult 表示 URL 解析结果
type ParseResult struct {
    Type     ResourceType
    Owner    string
    Repo     string
    Number   int
}
```

### 1.2 主要函数

```go
// ParseURL 解析 GitHub URL 并返回资源信息
// 返回 error 当 URL 格式无效或不支持的类型时
func ParseURL(url string) (*ParseResult, error)
```

---

## 2. `internal/github` - GitHub API 交互

### 2.1 数据结构

```go
// Client 表示 GitHub API 客户端
type Client struct {
    token     string
    client    *http.Client
    endpoint  string
}

// User 表示 GitHub 用户信息
type User struct {
    Login     string
    AvatarURL string
    URL       string
}

// ReactionCounts 表示 reactions 统计
type ReactionCounts struct {
    ThumbsUp    int
    ThumbsDown  int
    Laugh       int
    Hooray      int
    Confused    int
    Heart       int
    Rocket      int
    Eyes        int
}

// Comment 表示一条评论
type Comment struct {
    Author       *User
    Body         string
    CreatedAt    time.Time
    UpdatedAt    time.Time
    Reactions    *ReactionCounts
    IsAnswer     bool // Discussion 特有：标记为答案
}

// IssueData 表示 Issue/PR/Discussion 的完整数据
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
}
```

### 2.2 主要函数

```go
// NewClient 创建新的 GitHub API 客户端
// token 为空字符串时使用匿名访问（仅限公开仓库）
func NewClient(token string) *Client

// FetchIssueData 获取指定 Issue/PR/Discussion 的完整数据
// 根据 ParseResult.Type 调用相应的 GraphQL 查询
func (c *Client) FetchIssueData(ctx context.Context, result *ParseResult) (*IssueData, error)
```

---

## 3. `internal/converter` - 数据转换为 Markdown

### 3.1 数据结构

```go
// Options 控制输出选项
type Options struct {
    EnableReactions    bool
    EnableUserLinks    bool
}
```

### 3.2 主要函数

```go
// Convert 将 IssueData 转换为 Markdown 格式
// 返回完整的 Markdown 字符串，包含 YAML frontmatter 和正文
func Convert(data *github.IssueData, opts *Options) (string, error)

// WriteFile 将 Markdown 写入文件
// 如果 path 为空字符串，则写入 stdout
func WriteFile(markdown, path string) error
```

---

## 4. `internal/config` - 配置管理

### 4.1 数据结构

```go
// Config 表示应用配置
type Config struct {
    GitHubToken   string
    EnableReactions bool
    EnableUserLinks bool
}
```

### 4.2 主要函数

```go
// LoadConfig 从环境变量加载配置
// 支持：GITHUB_TOKEN（可选）
func LoadConfig() *Config
```

---

## 5. `internal/cli` - 命令行接口

### 5.1 主要函数

```go
// Run 执行 CLI 主逻辑
// args 为命令行参数（不含程序名）
// 返回 exit code
func Run(args []string) int

// Execute 是实际的执行函数，便于测试
func Execute(cfg *config.Config, args []string) error
```

---

## 6. 依赖关系图

```
cmd/issue2md/
    └── main.go
        └── cli.Run()
            ├── config.LoadConfig()
            ├── parser.ParseURL()
            ├── github.NewClient()
            ├── client.FetchIssueData()
            ├── converter.Convert()
            └── converter.WriteFile()
```

---

## 7. 错误处理约定

所有包都应定义自己的错误类型，使用 `fmt.Errorf` 包装：

```go
var (
    ErrInvalidURL    = errors.New("invalid GitHub URL")
    ErrUnsupportedType = errors.New("unsupported resource type")
)
```

---

**设计说明：**
- 无接口抽象，直接使用结构体和方法（符合"简单性原则"）
- 依赖通过函数参数显式传递（符合"无全局变量"原则）
- 每个包只做一件事，职责清晰（符合"单一职责原则"）
