// Package config provides configuration management for the issue2md application.
package config

import (
	"net/http"
	"time"
)

// Config holds the application configuration.
type Config struct {
	// GitHubToken is the personal access token for GitHub API authentication.
	// Optional but recommended for higher rate limits.
	GitHubToken string

	// UserAgent is the user agent string for HTTP requests.
	UserAgent string

	// APITimeout is the timeout for GitHub API requests.
	APITimeout time.Duration

	// HTTPClient is the HTTP client used for API requests.
	HTTPClient *http.Client

	// Options control conversion behavior.
	Options ConvertOptions
}

// ConvertOptions controls how resources are converted to Markdown.
type ConvertOptions struct {
	// EnableReactions includes reaction counts in the output.
	EnableReactions bool

	// EnableUserLinks converts @username to markdown links.
	EnableUserLinks bool

	// OutputFile is the path to the output file.
	// Empty string means stdout.
	OutputFile string
}

// CLIArgs represents the parsed command line arguments.
type CLIArgs struct {
	// URL is the GitHub resource URL to convert.
	URL string

	// OutputFile is the optional output file path.
	OutputFile string

	// EnableReactions enables reaction output.
	EnableReactions bool

	// EnableUserLinks enables user link rendering.
	EnableUserLinks bool

	// ShowHelp indicates whether to display help information.
	ShowHelp bool

	// ShowVersion indicates whether to display version information.
	ShowVersion bool
}

// DefaultConfig returns a configuration with default values.
func DefaultConfig() *Config {
	return &Config{
		UserAgent:  "issue2md",
		APITimeout: 30 * time.Second,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Options: ConvertOptions{
			EnableReactions:  false,
			EnableUserLinks:  false,
			OutputFile:       "",
		},
	}
}

// GitHubToken returns the configured GitHub token.
func (c *Config) GitHubTokenOrEmpty() string {
	return c.GitHubToken
}

// UserAgent returns the configured user agent.
func (c *Config) UserAgentOrEmpty() string {
	return c.UserAgent
}

// APITimeoutOrDefault returns the configured timeout or the default.
func (c *Config) APITimeoutOrDefault() time.Duration {
	if c.APITimeout > 0 {
		return c.APITimeout
	}
	return 30 * time.Second
}

// HTTPClientOrDefault returns the configured HTTP client or creates a new one.
func (c *Config) HTTPClientOrDefault() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{
		Timeout: c.APITimeoutOrDefault(),
	}
}

// ToConvertOptions converts CLIArgs to ConvertOptions.
func (a *CLIArgs) ToConvertOptions() ConvertOptions {
	return ConvertOptions{
		EnableReactions: a.EnableReactions,
		EnableUserLinks: a.EnableUserLinks,
		OutputFile:      a.OutputFile,
	}
}
