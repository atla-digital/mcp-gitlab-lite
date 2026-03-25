package tools

import (
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/descriptions"
	"github.com/mark3labs/mcp-go/mcp"
)

// Category is a typed string so mis-spellings are compile errors.
type Category string

const (
	CatProjects  Category = "Projects"
	CatIssues    Category = "Issues"
	CatMRs       Category = "Merge Requests"
	CatPipelines Category = "Pipelines"
	CatJobs      Category = "Jobs"
	CatRepos     Category = "Repositories"
	CatUsers     Category = "Users"
)

// Entry is one registered tool: its category, MCP schema, and handler.
type Entry struct {
	Category Category
	Tool     *mcp.Tool
	Handler  GLHandler
}

// tag stamps a single category onto every entry in the slice.
func tag(cat Category, entries ...Entry) []Entry {
	for i := range entries {
		entries[i].Category = cat
	}
	return entries
}

// All returns every tool in registration order.
func All(d descriptions.Catalog) []Entry {
	groups := [][]Entry{
		ProjectEntries(d),
		IssueEntries(d),
		MREntries(d),
		PipelineEntries(d),
		JobEntries(d),
		RepoEntries(d),
		UserEntries(d),
	}

	total := 0
	for _, g := range groups {
		total += len(g)
	}

	e := make([]Entry, 0, total)
	for _, g := range groups {
		e = append(e, g...)
	}

	return e
}
