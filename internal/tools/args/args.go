// Package args provides typed, named accessors for MCP tool arguments.
// All JSON numbers arrive as float64; domain methods handle conversion.
package args

import (
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
)

// R wraps the raw arguments map with typed accessors.
// The zero value is safe; missing keys return zero values.
type R struct {
	m map[string]any
}

// From constructs an R from a CallToolRequest.
func From(req mcp.CallToolRequest) R { return R{m: req.GetArguments()} }

// ── Generic accessors ─────────────────────────────────────────────────────────

func (r R) Str(key string) string {
	v, _ := r.m[key].(string)
	return v
}

func (r R) Int(key string) int {
	switch v := r.m[key].(type) {
	case float64:
		return int(v)
	case string:
		n, _ := strconv.Atoi(v)
		return n
	}
	return 0
}

func (r R) Int64(key string) int64 {
	switch v := r.m[key].(type) {
	case float64:
		return int64(v)
	case string:
		n, _ := strconv.ParseInt(v, 10, 64)
		return n
	}
	return 0
}

func (r R) Bool(key string) (val, ok bool) {
	v, ok := r.m[key].(bool)
	return v, ok
}

// ── IDs ───────────────────────────────────────────────────────────────────────

func (r R) ProjectID() string { return r.Str("project_id") }
func (r R) GroupID() string   { return r.Str("group_id") }
func (r R) IssueIID() int64   { return r.Int64("issue_iid") }
func (r R) MrIID() int64      { return r.Int64("mr_iid") }
func (r R) PipelineID() int64 { return r.Int64("pipeline_id") }
func (r R) JobID() int64      { return r.Int64("job_id") }
func (r R) UserID() int64     { return r.Int64("user_id") }

// ── Common string params ──────────────────────────────────────────────────────

func (r R) Body() string     { return r.Str("body") }
func (r R) Title() string    { return r.Str("title") }
func (r R) Search() string   { return r.Str("search") }
func (r R) Ref() string      { return r.Str("ref") }
func (r R) Branch() string   { return r.Str("branch") }
func (r R) SHA() string      { return r.Str("sha") }
func (r R) FilePath() string { return r.Str("file_path") }
func (r R) Scope() string    { return r.Str("scope") }
func (r R) Labels() string   { return r.Str("labels") }

// State returns the state argument, substituting def if empty.
func (r R) State(def string) string {
	if s := r.Str("state"); s != "" {
		return s
	}
	return def
}

// ── Pagination ─────────────────────────────────────────────────────────────────

// Page returns the page argument, defaulting to 1.
func (r R) Page() int64 {
	if p := r.Int64("page"); p > 0 {
		return p
	}
	return 1
}

// ── Wait-tool helpers ─────────────────────────────────────────────────────────

// Timeout returns timeout_seconds, substituting def if unset.
func (r R) Timeout(def int) int {
	if t := r.Int("timeout_seconds"); t > 0 {
		return t
	}
	return def
}

// PollInterval returns poll_interval_seconds, clamped to min, with def as fallback.
func (r R) PollInterval(minVal, def int) int {
	if p := r.Int("poll_interval_seconds"); p >= minVal {
		return p
	}
	return def
}
