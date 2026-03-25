package tools

import (
	"context"

	"github.com/atla-digital/mcp-gitlab-lite/internal/model"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/args"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/descriptions"
	"github.com/mark3labs/mcp-go/mcp"
	gl "gitlab.com/gitlab-org/api/client-go"
)

func IssueEntries(d descriptions.Catalog) []Entry {
	return tag(CatIssues,
		Entry{
			Tool: d.For("list_issues").
				Str("project_id", mcp.Required()).
				Str("state").
				Str("labels").
				Str("assignee_username").
				Str("search").
				Num("page").
				Build(),
			Handler: listIssues,
		},
		Entry{
			Tool: d.For("get_issue").
				Str("project_id", mcp.Required()).
				Num("issue_iid", mcp.Required()).
				Build(),
			Handler: getIssue,
		},
		Entry{
			Tool: d.For("create_issue").
				Str("project_id", mcp.Required()).
				Str("title", mcp.Required()).
				Str("description").
				Str("labels").
				Build(),
			Handler: createIssue,
		},
		Entry{
			Tool: d.For("update_issue").
				Str("project_id", mcp.Required()).
				Num("issue_iid", mcp.Required()).
				Str("title").
				Str("description").
				Str("state_event").
				Str("labels").
				Build(),
			Handler: updateIssue,
		},
		Entry{
			Tool: d.For("list_issue_notes").
				Str("project_id", mcp.Required()).
				Num("issue_iid", mcp.Required()).
				Build(),
			Handler: listIssueNotes,
		},
		Entry{
			Tool: d.For("create_issue_note").
				Str("project_id", mcp.Required()).
				Num("issue_iid", mcp.Required()).
				Str("body", mcp.Required()).
				Build(),
			Handler: createIssueNote,
		},
		Entry{
			Tool: d.For("search_issues").
				Str("search", mcp.Required()).
				Str("state").
				Num("page").
				Build(),
			Handler: searchIssues,
		},
	)
}

func listIssues(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.ListProjectIssuesOptions{
		State:       gl.Ptr(a.State("opened")),
		ListOptions: gl.ListOptions{Page: a.Page(), PerPage: model.DefaultPerPage},
	}
	if l := a.Labels(); l != "" {
		ll := gl.LabelOptions{l}
		opts.Labels = &ll
	}
	if u := a.Str("assignee_username"); u != "" {
		opts.AssigneeUsername = gl.Ptr(u)
	}
	if q := a.Search(); q != "" {
		opts.Search = gl.Ptr(q)
	}
	issues, resp, err := client.Issues.ListProjectIssues(a.ProjectID(), opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.NewPaged(model.ToIssueRefs(issues), resp, a.Page()))
}

func getIssue(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	issue, _, err := client.Issues.GetIssue(a.ProjectID(), a.IssueIID())
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToIssue(issue))
}

func createIssue(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.CreateIssueOptions{Title: gl.Ptr(a.Title())}
	if d := a.Str("description"); d != "" {
		opts.Description = gl.Ptr(d)
	}
	if l := a.Labels(); l != "" {
		ll := gl.LabelOptions{l}
		opts.Labels = &ll
	}
	issue, _, err := client.Issues.CreateIssue(a.ProjectID(), opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToIssue(issue))
}

func updateIssue(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.UpdateIssueOptions{}
	if t := a.Title(); t != "" {
		opts.Title = gl.Ptr(t)
	}
	if d := a.Str("description"); d != "" {
		opts.Description = gl.Ptr(d)
	}
	if s := a.Str("state_event"); s != "" {
		opts.StateEvent = gl.Ptr(s)
	}
	if l := a.Labels(); l != "" {
		ll := gl.LabelOptions{l}
		opts.Labels = &ll
	}
	issue, _, err := client.Issues.UpdateIssue(a.ProjectID(), a.IssueIID(), opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToIssue(issue))
}

func listIssueNotes(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	notes, _, err := client.Notes.ListIssueNotes(a.ProjectID(), a.IssueIID(),
		&gl.ListIssueNotesOptions{ListOptions: gl.ListOptions{PerPage: model.MaxPerPage}})
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToNotes(notes))
}

func createIssueNote(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	note, _, err := client.Notes.CreateIssueNote(a.ProjectID(), a.IssueIID(),
		&gl.CreateIssueNoteOptions{Body: gl.Ptr(a.Body())})
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToNote(note))
}

func searchIssues(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.ListIssuesOptions{
		Search:      gl.Ptr(a.Search()),
		State:       gl.Ptr(a.State("opened")),
		Scope:       gl.Ptr("all"),
		ListOptions: gl.ListOptions{Page: a.Page(), PerPage: model.DefaultPerPage},
	}
	issues, resp, err := client.Issues.ListIssues(opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.NewPaged(model.ToIssueRefs(issues), resp, a.Page()))
}
