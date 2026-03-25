package model

import "time"

// Paged wraps a slice with pagination metadata.
type Paged[T any] struct {
	Items     []T   `json:"items"`
	Page      int64 `json:"page"`
	PerPage   int64 `json:"per_page"`
	Total     int64 `json:"total"`
	Truncated bool  `json:"truncated"`
}

// PersonRef is a compact author/assignee reference.
type PersonRef struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// ProjectRef is a compact project reference for list responses.
type ProjectRef struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
	WebURL            string `json:"web_url"`
}

// Project is the full project representation.
type Project struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
	Description       string `json:"description,omitempty"`
	Visibility        string `json:"visibility"`
	DefaultBranch     string `json:"default_branch"`
	WebURL            string `json:"web_url"`
}

// IssueRef is a compact issue reference for list responses.
type IssueRef struct {
	IID    int64     `json:"iid"`
	Title  string    `json:"title"`
	State  string    `json:"state"`
	Author PersonRef `json:"author"`
	Labels []string  `json:"labels"`
	WebURL string    `json:"web_url"`
}

// Issue is the full issue representation.
type Issue struct {
	IID         int64       `json:"iid"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	State       string      `json:"state"`
	Author      PersonRef   `json:"author"`
	Assignees   []PersonRef `json:"assignees"`
	Labels      []string    `json:"labels"`
	CreatedAt   *time.Time  `json:"created_at,omitempty"`
	UpdatedAt   *time.Time  `json:"updated_at,omitempty"`
	WebURL      string      `json:"web_url"`
}

// MRRef is a compact merge request reference for list responses.
type MRRef struct {
	IID          int64     `json:"iid"`
	Title        string    `json:"title"`
	State        string    `json:"state"`
	SourceBranch string    `json:"source_branch"`
	TargetBranch string    `json:"target_branch"`
	Author       PersonRef `json:"author"`
	WebURL       string    `json:"web_url"`
}

// MR is the full merge request representation.
type MR struct {
	IID          int64       `json:"iid"`
	Title        string      `json:"title"`
	Description  string      `json:"description,omitempty"`
	State        string      `json:"state"`
	SourceBranch string      `json:"source_branch"`
	TargetBranch string      `json:"target_branch"`
	Author       PersonRef   `json:"author"`
	Assignees    []PersonRef `json:"assignees"`
	Labels       []string    `json:"labels"`
	MergeStatus  string      `json:"merge_status"`
	SHA          string      `json:"sha,omitempty"`
	CreatedAt    *time.Time  `json:"created_at,omitempty"`
	UpdatedAt    *time.Time  `json:"updated_at,omitempty"`
	WebURL       string      `json:"web_url"`
}

// PipelineRef is a compact pipeline reference for list responses.
type PipelineRef struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
	Ref    string `json:"ref"`
	SHA    string `json:"sha"`
	WebURL string `json:"web_url"`
}

// Pipeline is the full pipeline representation.
type Pipeline struct {
	ID        int64      `json:"id"`
	Status    string     `json:"status"`
	Ref       string     `json:"ref"`
	SHA       string     `json:"sha"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Duration  int        `json:"duration"`
	WebURL    string     `json:"web_url"`
}

// Job represents a CI/CD job.
type Job struct {
	ID       int64       `json:"id"`
	Name     string      `json:"name"`
	Stage    string      `json:"stage"`
	Status   string      `json:"status"`
	Pipeline PipelineRef `json:"pipeline"`
	WebURL   string      `json:"web_url"`
	Duration float64     `json:"duration"`
}

// Note represents a comment on an issue or MR.
type Note struct {
	ID        int64      `json:"id"`
	Body      string     `json:"body"`
	Author    PersonRef  `json:"author"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	System    bool       `json:"system"`
}

// Branch represents a git branch.
type Branch struct {
	Name      string    `json:"name"`
	Commit    CommitRef `json:"commit"`
	Protected bool      `json:"protected"`
	Default   bool      `json:"default"`
}

// CommitRef is a compact commit reference.
type CommitRef struct {
	ID         string     `json:"id"`
	ShortID    string     `json:"short_id"`
	Title      string     `json:"title"`
	AuthorName string     `json:"author_name"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
}

// User represents a GitLab user.
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	State    string `json:"state"`
	WebURL   string `json:"web_url"`
}
