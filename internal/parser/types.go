// Package parser provides URL parsing and type identification for GitHub URLs.
package parser

import (
	"regexp"

	"github.com/bigwhite/issue2md/internal/github"
)

// ResourceURL represents a parsed GitHub resource URL.
type ResourceURL struct {
	OriginalURL string            // The original URL string
	Owner       string            // Repository owner
	Repo        string            // Repository name
	Number      int64             // Issue/PR/Discussion number
	Type        github.ResourceType // Resource type (issue/pr/discussion)
	IsValid     bool              // Whether the URL is valid
}

// URLPattern represents a URL pattern for matching GitHub URLs.
type URLPattern struct {
	Pattern     *regexp.Regexp // Compiled regex pattern
	Type        github.ResourceType // Resource type this pattern matches
	Description string         // Human-readable description
}

// SupportedPatterns contains all supported URL patterns.
var SupportedPatterns = []URLPattern{
	{
		// Issue pattern: https://github.com/{owner}/{repo}/issues/{number}
		Pattern:     regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+)/issues/(\d+)$`),
		Type:        github.ResourceTypeIssue,
		Description: "GitHub Issue URL",
	},
	{
		// Pull Request pattern: https://github.com/{owner}/{repo}/pull/{number}
		Pattern:     regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+)/pull/(\d+)$`),
		Type:        github.ResourceTypePR,
		Description: "GitHub Pull Request URL",
	},
	{
		// Discussion pattern: https://github.com/{owner}/{repo}/discussions/{number}
		Pattern:     regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+)/discussions/(\d+)$`),
		Type:        github.ResourceTypeDiscussion,
		Description: "GitHub Discussion URL",
	},
}

// SupportedTypes returns a list of all supported resource types.
func SupportedTypes() []github.ResourceType {
	types := make([]github.ResourceType, len(SupportedPatterns))
	for i, p := range SupportedPatterns {
		types[i] = p.Type
	}
	return types
}

// Match returns the first pattern that matches the given URL.
// Returns nil if no pattern matches.
func Match(url string) *URLPattern {
	for _, p := range SupportedPatterns {
		if p.Pattern.MatchString(url) {
			return &p
		}
	}
	return nil
}

// IsSupported returns true if the URL matches a supported pattern.
func IsSupported(url string) bool {
	return Match(url) != nil
}
