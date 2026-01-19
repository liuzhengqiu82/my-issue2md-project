// Package github provides error definitions for GitHub API operations.
package github

import "fmt"

// Predefined error types for GitHub operations.

// ErrInvalidURL is returned when a URL is malformed or doesn't match expected patterns.
type ErrInvalidURL struct {
	URL string
	Err error
}

// Error returns the error message.
func (e *ErrInvalidURL) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("invalid URL %q: %v", e.URL, e.Err)
	}
	return fmt.Sprintf("invalid URL: %s", e.URL)
}

// Unwrap returns the underlying error.
func (e *ErrInvalidURL) Unwrap() error {
	return e.Err
}

// ErrUnsupportedURL is returned when a URL is valid but not a supported GitHub resource type.
type ErrUnsupportedURL struct {
	URL string
}

// Error returns the error message.
func (e *ErrUnsupportedURL) Error() string {
	return fmt.Sprintf("unsupported URL type: %s (supported: issues, pull requests, discussions)", e.URL)
}

// ErrAPIError is returned when the GitHub API returns an error response.
type ErrAPIError struct {
	Message string
	StatusCode int
	Err error
}

// Error returns the error message.
func (e *ErrAPIError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("API error: %s", e.Message)
}

// Unwrap returns the underlying error.
func (e *ErrAPIError) Unwrap() error {
	return e.Err
}

// ErrAuthRequired is returned when authentication is required but not provided.
type ErrAuthRequired struct {
	Resource string
}

// Error returns the error message.
func (e *ErrAuthRequired) Error() string {
	return fmt.Sprintf("authentication required for resource: %s (set GITHUB_TOKEN environment variable)", e.Resource)
}

// ErrResourceNotFound is returned when the requested resource doesn't exist.
type ErrResourceNotFound struct {
	URL string
}

// Error returns the error message.
func (e *ErrResourceNotFound) Error() string {
	return fmt.Sprintf("resource not found: %s", e.URL)
}

// ErrRateLimitExceeded is returned when GitHub API rate limit is exceeded.
type ErrRateLimitExceeded struct {
	ResetTime int64 // Unix timestamp when the limit resets
}

// Error returns the error message.
func (e *ErrRateLimitExceeded) Error() string {
	return fmt.Sprintf("rate limit exceeded, resets at %d", e.ResetTime)
}

// ErrNetworkError is returned when there's a network-related error.
type ErrNetworkError struct {
	Op string // Operation that failed (e.g., "GET", "POST")
	Err error // Underlying error
}

// Error returns the error message.
func (e *ErrNetworkError) Error() string {
	return fmt.Sprintf("network error during %s: %v", e.Op, e.Err)
}

// Unwrap returns the underlying error.
func (e *ErrNetworkError) Unwrap() error {
	return e.Err
}

// ErrInvalidResponse is returned when the GitHub API response cannot be parsed.
type ErrInvalidResponse struct {
	Reason string
	Err error
}

// Error returns the error message.
func (e *ErrInvalidResponse) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("invalid API response: %s: %v", e.Reason, e.Err)
	}
	return fmt.Sprintf("invalid API response: %s", e.Reason)
}

// Unwrap returns the underlying error.
func (e *ErrInvalidResponse) Unwrap() error {
	return e.Err
}
