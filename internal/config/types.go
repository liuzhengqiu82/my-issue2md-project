package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config 应用程序配置
// 包含GitHub令牌、输出配置和解析器配置
type Config struct {
	GitHubToken string      `json:"github_token"`
	Output      OutputConfig `json:"output"`
	Parser      ParserConfig `json:"parser"`
}

// OutputConfig 输出配置
type OutputConfig struct {
	Format      string `json:"format"`       // markdown, html, json
	Filename    string `json:"filename"`
	Destination string `json:"destination"`
	Overwrite   bool   `json:"overwrite"`
}

// ParserConfig 解析器配置
type ParserConfig struct {
	IncludeComments    bool `json:"include_comments"`
	IncludeMetadata    bool `json:"include_metadata"`
	IncludeTimestamps  bool `json:"include_timestamps"`
	IncludeUserLinks   bool `json:"include_user_links"`
	EmojisEnabled      bool `json:"emojis_enabled"`
	PreserveLineBreaks bool `json:"preserve_line_breaks"`
}

// Environment 环境变量配置
type Environment struct {
	GitHubToken string
	Debug       bool
	NoColor     bool
}

// DefaultConfig 返回默认配置
// 返回一个包含默认值的配置对象
// 返回值: *Config - 默认配置实例
func DefaultConfig() *Config {
	return &Config{
		Output: OutputConfig{
			Format:      "markdown",
			Destination: "output",
			Overwrite:   false,
		},
		Parser: ParserConfig{
			IncludeComments:    true,
			IncludeMetadata:    true,
			IncludeTimestamps:  true,
			IncludeUserLinks:   true,
			EmojisEnabled:      true,
			PreserveLineBreaks: true,
		},
	}
}

// LoadFromEnv 从环境变量加载配置
// 将环境变量中的配置值加载到当前配置对象中
// 参数: 无
// 返回值: 无
func (c *Config) LoadFromEnv() {
	env := GetEnvironment()

	if token := env.GitHubToken; token != "" {
		c.GitHubToken = token
	}

	if debug := env.Debug; debug {
		c.Parser.EmojisEnabled = false // Example debug setting
	}
}

// GetEnvironment 获取环境变量
// 从系统环境变量中读取配置并返回Environment结构体
// 参数: 无
// 返回值: *Environment - 包含环境变量值的结构体指针
// 异常: 无，如果环境变量解析失败会使用默认值
func GetEnvironment() *Environment {
	env := &Environment{
		GitHubToken: os.Getenv("GITHUB_TOKEN"),
		Debug:       getBoolEnv("DEBUG", false),
		NoColor:     getBoolEnv("NO_COLOR", false),
	}
	return env
}

// getBoolEnv 获取布尔环境变量
// 从环境变量中读取布尔值，如果解析失败则记录警告并返回默认值
// 参数:
//   - key: 环境变量名
//   - defaultValue: 默认值，在解析失败时返回
// 返回值: bool - 解析后的布尔值或默认值
// 异常: 无，会通过fmt.Printf输出警告信息
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		} else {
			// 解析失败时记录错误并返回默认值
			fmt.Printf("warning: failed to parse boolean environment variable %q: %v, using default value %v\n", key, err, defaultValue)
		}
	}
	return defaultValue
}

// Validate 验证配置
// 验证配置对象的必填字段是否有效
// 参数: 无
// 返回值: error - 如果验证失败返回ValidationError，否则返回nil
// 异常: 可能返回ValidationError，包含错误字段和描述信息
func (c *Config) Validate() error {
	if c.GitHubToken == "" {
		return &ValidationError{
			Field:   "github_token",
			Message: "GitHub token is required",
		}
	}

	if c.Output.Format == "" {
		return &ValidationError{
			Field:   "output.format",
			Message: "Output format is required",
		}
	}

	return nil
}

// ValidationError 配置验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
