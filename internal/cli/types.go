package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// CLI 命令行接口
// 提供命令行参数解析和命令执行功能
type CLI struct {
	name   string
	args   []string
	output *Output
}

// NewCLI 创建新的CLI实例
// 创建一个CLI实例并初始化默认的输出配置
// 参数:
//   - name: 应用程序名称
//   - args: 命令行参数数组
// 返回值: *CLI - 新创建的CLI实例
func NewCLI(name string, args []string) *CLI {
	return &CLI{
		name: name,
		args: args,
		output: &Output{
			Writer: os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}
}

// Output 输出配置
// 定义输出和错误输出的目标，使用io.Writer接口确保类型安全
type Output struct {
	Writer      io.Writer
	ErrorWriter io.Writer
}

// Command 命令定义
type Command struct {
	Name        string
	Description string
	Flags       *flag.FlagSet
	Run         func(*Context) error
}

// Context 命令执行上下文
type Context struct {
	Args   []string
	Flags  map[string]string
	Output *Output
}

// CommandRegistry 命令注册表
type CommandRegistry struct {
	commands map[string]*Command
}

// NewCommandRegistry 创建命令注册表
// 创建一个新的命令注册表实例
// 参数: 无
// 返回值: *CommandRegistry - 新创建的命令注册表实例
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]*Command),
	}
}

// Register 注册命令
// 将一个命令注册到命令注册表中
// 参数:
//   - cmd: 要注册的命令对象
// 返回值: error - 如果命令已存在则返回错误，否则返回nil
// 异常: 可能返回fmt.Errorf错误，表示命令已存在
func (cr *CommandRegistry) Register(cmd *Command) error {
	if _, exists := cr.commands[cmd.Name]; exists {
		return fmt.Errorf("command %s already registered", cmd.Name)
	}
	cr.commands[cmd.Name] = cmd
	return nil
}

// Get 获取命令
// 根据命令名称获取已注册的命令
// 参数:
//   - name: 命令名称
// 返回值: (*Command, bool) - 如果找到命令则返回命令对象和true，否则返回nil和false
func (cr *CommandRegistry) Get(name string) (*Command, bool) {
	cmd, exists := cr.commands[name]
	return cmd, exists
}

// List 列出所有命令
// 返回注册表中所有已注册的命令列表
// 参数: 无
// 返回值: []*Command - 命令对象的切片
func (cr *CommandRegistry) List() []*Command {
	commands := make([]*Command, 0, len(cr.commands))
	for _, cmd := range cr.commands {
		commands = append(commands, cmd)
	}
	return commands
}

// Error CLI错误类型
type Error struct {
	Message string
	Code    int
}

func (e *Error) Error() string {
	return e.Message
}

// NewError 创建CLI错误
// 创建一个CLI错误实例
// 参数:
//   - message: 错误描述信息
//   - code: 错误代码
// 返回值: *Error - 新创建的CLI错误实例
func NewError(message string, code int) *Error {
	return &Error{
		Message: message,
		Code:    code,
	}
}
