// Package github provides data structures for GitHub API resources.
package github

import "time"

// User represents a GitHub user.
type User struct {
	Login     string // Username
	AvatarURL string // URL to user's avatar image
	URL       string // URL to user's GitHub profile
}

// Comment represents a comment on an Issue, PR, or Discussion.
type Comment struct {
	ID        int64     // Unique comment identifier
	Author    User      // Author of the comment
	Body      string    // Comment content in markdown
	CreatedAt time.Time // Timestamp when comment was created
	UpdatedAt time.Time // Timestamp when comment was last updated
	Reactions Reactions // Reaction counts for this comment

	// Discussion-specific fields
	IsAnswer bool   // True if this comment is marked as the accepted answer
	ReplyTo  int64  // ID of the comment this replies to (for Discussions)
}

// Reactions represents reaction counts for a resource.
type Reactions struct {
	ThumbsUp   int // +1 reactions count
	ThumbsDown int // -1 reactions count
	Laugh      int // laugh reactions count
	Hooray     int // hooray reactions count
	Confused   int // confused reactions count
	Heart      int // heart reactions count
	Rocket     int // rocket reactions count
	Eyes       int // eyes reactions count
	TotalCount int // Total reactions count
}

// Label represents a GitHub issue/PR label.
type Label struct {
	Name        string // Label name
	Description string // Label description
	Color       string // Label color (hex)
}

// Milestone represents a GitHub milestone.
type Milestone struct {
	Title       string    // Milestone title
	Description string    // Milestone description
	State       string    // Milestone state (open, closed)
	DueDate     time.Time // Milestone due date
}

// DiscussionCategory represents a GitHub discussion category.
type DiscussionCategory struct {
	ID          int64  // Category ID
	Slug        string // Category slug (URL-friendly identifier)
	Name        string // Category display name
	Description string // Category description
}

// ResourceType represents the type of GitHub resource.
type ResourceType string

const (
	// ResourceTypeIssue represents a GitHub Issue.
	ResourceTypeIssue ResourceType = "issue"
	// ResourceTypePR represents a GitHub Pull Request.
	ResourceTypePR ResourceType = "pr"
	// ResourceTypeDiscussion represents a GitHub Discussion.
	ResourceTypeDiscussion ResourceType = "discussion"
)

// ResourceStatus represents the status of a GitHub resource.
type ResourceStatus string

const (
	// StatusOpen represents an open resource.
	StatusOpen ResourceStatus = "open"
	// StatusClosed represents a closed resource.
	StatusClosed ResourceStatus = "closed"
	// StatusMerged represents a merged pull request.
	StatusMerged ResourceStatus = "merged"
)

// Resource is the base interface for all GitHub resources.
type Resource interface {
	// GetType returns the type of the resource.
	GetType() ResourceType
	// GetStatus returns the status of the resource.
	GetStatus() ResourceStatus
	// GetTitle returns the title of the resource.
	GetTitle() string
	// GetAuthor returns the author of the resource.
	GetAuthor() User
	// GetCreatedAt returns the creation timestamp.
	GetCreatedAt() time.Time
	// GetUpdatedAt returns the last update timestamp.
	GetUpdatedAt() time.Time
	// GetBody returns the body content.
	GetBody() string
	// GetComments returns all comments.
	GetComments() []Comment
	// GetReactions returns the reaction counts.
	GetReactions() Reactions
	// GetCommentCount returns the total number of comments.
	GetCommentCount() int
}

// Issue represents a GitHub Issue.
type Issue struct {
	Number         int64        // Issue number
	Title          string       // Issue title
	Body           string       // Issue body content in markdown
	Author         User         // Issue author
	State          ResourceStatus // Issue state (open/closed)
	CreatedAt      time.Time    // Creation timestamp
	UpdatedAt      time.Time    // Last update timestamp
	ClosedAt       *time.Time   // Closure timestamp (nil if open)
	Comments       []Comment    // All issue comments
	Reactions      Reactions    // Reaction counts for the issue
	Labels         []Label      // Issue labels
	Milestone      *Milestone   // Associated milestone (nil if none)
	TotalComments  int          // Total comment count
	URL            string       // GitHub URL for this issue
	RepositoryURL  string       // Repository URL
}

// GetType returns the resource type.
func (i *Issue) GetType() ResourceType {
	return ResourceTypeIssue
}

// GetStatus returns the issue status.
func (i *Issue) GetStatus() ResourceStatus {
	return i.State
}

// GetTitle returns the issue title.
func (i *Issue) GetTitle() string {
	return i.Title
}

// GetAuthor returns the issue author.
func (i *Issue) GetAuthor() User {
	return i.Author
}

// GetCreatedAt returns the creation timestamp.
func (i *Issue) GetCreatedAt() time.Time {
	return i.CreatedAt
}

// GetUpdatedAt returns the last update timestamp.
func (i *Issue) GetUpdatedAt() time.Time {
	return i.UpdatedAt
}

// GetBody returns the issue body.
func (i *Issue) GetBody() string {
	return i.Body
}

// GetComments returns all comments.
func (i *Issue) GetComments() []Comment {
	return i.Comments
}

// GetReactions returns the reaction counts.
func (i *Issue) GetReactions() Reactions {
	return i.Reactions
}

// GetCommentCount returns the total comment count.
func (i *Issue) GetCommentCount() int {
	return i.TotalComments
}

// PullRequest represents a GitHub Pull Request.
type PullRequest struct {
	Number         int64        // PR number
	Title          string       // PR title
	Body           string       // PR body content in markdown
	Author         User         // PR author
	State          ResourceStatus // PR state (open/closed/merged)
	CreatedAt      time.Time    // Creation timestamp
	UpdatedAt      time.Time    // Last update timestamp
	ClosedAt       *time.Time   // Closure timestamp (nil if open)
	MergedAt       *time.Time   // Merge timestamp (nil if not merged)
	Comments       []Comment    // All PR comments (including review comments)
	Reactions      Reactions    // Reaction counts for the PR
	Labels         []Label      // PR labels
	Milestone      *Milestone   // Associated milestone (nil if none)
	TotalComments  int          // Total comment count
	URL            string       // GitHub URL for this PR
	RepositoryURL  string       // Repository URL
	IsDraft        bool         // True if this is a draft PR
	Mergeable      bool         // True if the PR can be merged
	Additions      int          // Number of lines added
	Deletions      int          // Number of lines deleted
}

// GetType returns the resource type.
func (p *PullRequest) GetType() ResourceType {
	return ResourceTypePR
}

// GetStatus returns the PR status.
func (p *PullRequest) GetStatus() ResourceStatus {
	return p.State
}

// GetTitle returns the PR title.
func (p *PullRequest) GetTitle() string {
	return p.Title
}

// GetAuthor returns the PR author.
func (p *PullRequest) GetAuthor() User {
	return p.Author
}

// GetCreatedAt returns the creation timestamp.
func (p *PullRequest) GetCreatedAt() time.Time {
	return p.CreatedAt
}

// GetUpdatedAt returns the last update timestamp.
func (p *PullRequest) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

// GetBody returns the PR body.
func (p *PullRequest) GetBody() string {
	return p.Body
}

// GetComments returns all comments.
func (p *PullRequest) GetComments() []Comment {
	return p.Comments
}

// GetReactions returns the reaction counts.
func (p *PullRequest) GetReactions() Reactions {
	return p.Reactions
}

// GetCommentCount returns the total comment count.
func (p *PullRequest) GetCommentCount() int {
	return p.TotalComments
}

// Discussion represents a GitHub Discussion.
type Discussion struct {
	Number         int64               // Discussion number
	Title          string              // Discussion title
	Body           string              // Discussion body content in markdown
	Author         User                // Discussion author
	State          ResourceStatus      // Discussion state (open/closed)
	CreatedAt      time.Time           // Creation timestamp
	UpdatedAt      time.Time           // Last update timestamp
	ClosedAt       *time.Time          // Closure timestamp (nil if open)
	Comments       []Comment           // All discussion comments
	Reactions      Reactions           // Reaction counts for the discussion
	TotalComments  int                 // Total comment count
	URL            string              // GitHub URL for this discussion
	RepositoryURL  string              // Repository URL
	Category       DiscussionCategory  // Discussion category
	IsAnswered     bool                // True if an answer has been selected
	AnswerComment  *Comment            // The comment marked as answer (nil if none)
}

// GetType returns the resource type.
func (d *Discussion) GetType() ResourceType {
	return ResourceTypeDiscussion
}

// GetStatus returns the discussion status.
func (d *Discussion) GetStatus() ResourceStatus {
	return d.State
}

// GetTitle returns the discussion title.
func (d *Discussion) GetTitle() string {
	return d.Title
}

// GetAuthor returns the discussion author.
func (d *Discussion) GetAuthor() User {
	return d.Author
}

// GetCreatedAt returns the creation timestamp.
func (d *Discussion) GetCreatedAt() time.Time {
	return d.CreatedAt
}

// GetUpdatedAt returns the last update timestamp.
func (d *Discussion) GetUpdatedAt() time.Time {
	return d.UpdatedAt
}

// GetBody returns the discussion body.
func (d *Discussion) GetBody() string {
	return d.Body
}

// GetComments returns all comments.
func (d *Discussion) GetComments() []Comment {
	return d.Comments
}

// GetReactions returns the reaction counts.
func (d *Discussion) GetReactions() Reactions {
	return d.Reactions
}

// GetCommentCount returns the total comment count.
func (d *Discussion) GetCommentCount() int {
	return d.TotalComments
}
