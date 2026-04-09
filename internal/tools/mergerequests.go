package tools

import (
	"context"
	"fmt"

	"github.com/atla-digital/mcp-gitlab-lite/internal/model"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/args"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/descriptions"
	"github.com/mark3labs/mcp-go/mcp"
	gl "gitlab.com/gitlab-org/api/client-go"
)

func MREntries(d descriptions.Catalog) []Entry {
	return tag(CatMRs,
		Entry{
			Tool: d.For("list_merge_requests").
				Str("project_id", mcp.Required()).
				Str("state").
				Str("target_branch").
				Str("source_branch").
				Str("search").
				Num("page").
				Build(),
			Handler: listMRs,
		},
		Entry{
			Tool:    d.For("get_merge_request").Str("project_id", mcp.Required()).Num("merge_request_iid", mcp.Required()).Build(),
			Handler: getMR,
		},
		Entry{
			Tool: d.For("create_merge_request").
				Str("project_id", mcp.Required()).
				Str("title", mcp.Required()).
				Str("source_branch", mcp.Required()).
				Str("target_branch", mcp.Required()).
				Str("description").
				Bool("remove_source_branch").
				Bool("squash").
				Build(),
			Handler: createMR,
		},
		Entry{
			Tool: d.For("update_merge_request").
				Str("project_id", mcp.Required()).
				Num("merge_request_iid", mcp.Required()).
				Str("title").Str("description").Str("target_branch").
				Str("state_event").Str("labels").
				Build(),
			Handler: updateMR,
		},
		Entry{
			Tool: d.For("merge_merge_request").
				Str("project_id", mcp.Required()).
				Num("merge_request_iid", mcp.Required()).
				Str("merge_commit_message").
				Bool("squash").
				Bool("should_remove_source_branch").
				Build(),
			Handler: mergeMR,
		},
		Entry{
			Tool:    d.For("approve_merge_request").Str("project_id", mcp.Required()).Num("merge_request_iid", mcp.Required()).Build(),
			Handler: approveMR,
		},
		Entry{
			Tool:    d.For("list_mr_notes").Str("project_id", mcp.Required()).Num("merge_request_iid", mcp.Required()).Build(),
			Handler: listMRNotes,
		},
		Entry{
			Tool:    d.For("create_mr_note").Str("project_id", mcp.Required()).Num("merge_request_iid", mcp.Required()).Str("body", mcp.Required()).Build(),
			Handler: createMRNote,
		},
		Entry{
			Tool: d.For("create_mr_discussion").
				Str("project_id", mcp.Required()).
				Num("merge_request_iid", mcp.Required()).
				Str("body", mcp.Required()).
				Str("base_sha", mcp.Required()).
				Str("head_sha", mcp.Required()).
				Str("start_sha", mcp.Required()).
				Str("new_path").Str("old_path").
				Num("new_line").Num("old_line").
				Build(),
			Handler: createMRDiscussion,
		},
		Entry{
			Tool:    d.For("list_mr_diffs").Str("project_id", mcp.Required()).Num("merge_request_iid", mcp.Required()).Build(),
			Handler: listMRDiffs,
		},
		Entry{
			Tool:    d.For("list_mr_commits").Str("project_id", mcp.Required()).Num("merge_request_iid", mcp.Required()).Build(),
			Handler: listMRCommits,
		},
		Entry{
			Tool:    d.For("list_mr_pipelines").Str("project_id", mcp.Required()).Num("merge_request_iid", mcp.Required()).Build(),
			Handler: listMRPipelines,
		},
		Entry{
			Tool:    d.For("search_merge_requests").Str("search", mcp.Required()).Str("state").Num("page").Build(),
			Handler: searchMRs,
		},
	)
}

func listMRs(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.ListProjectMergeRequestsOptions{
		State:       gl.Ptr(a.State("opened")),
		ListOptions: gl.ListOptions{Page: a.Page(), PerPage: model.DefaultPerPage},
	}
	if b := a.Str("target_branch"); b != "" {
		opts.TargetBranch = gl.Ptr(b)
	}
	if b := a.Str("source_branch"); b != "" {
		opts.SourceBranch = gl.Ptr(b)
	}
	if q := a.Search(); q != "" {
		opts.Search = gl.Ptr(q)
	}
	mrs, resp, err := client.MergeRequests.ListProjectMergeRequests(a.ProjectID(), opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.NewPaged(model.ToMRRefsFromBasic(mrs), resp, a.Page()))
}

func getMR(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	mr, _, err := client.MergeRequests.GetMergeRequest(a.ProjectID(), a.MrIID(), nil)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToMR(mr))
}

func createMR(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.CreateMergeRequestOptions{
		Title:        gl.Ptr(a.Title()),
		SourceBranch: gl.Ptr(a.Str("source_branch")),
		TargetBranch: gl.Ptr(a.Str("target_branch")),
	}
	if d := a.Str("description"); d != "" {
		opts.Description = gl.Ptr(d)
	}
	if v, ok := a.Bool("remove_source_branch"); ok {
		opts.RemoveSourceBranch = gl.Ptr(v)
	}
	if v, ok := a.Bool("squash"); ok {
		opts.Squash = gl.Ptr(v)
	}
	mr, _, err := client.MergeRequests.CreateMergeRequest(a.ProjectID(), opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToMR(mr))
}

func updateMR(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.UpdateMergeRequestOptions{}
	if t := a.Title(); t != "" {
		opts.Title = gl.Ptr(t)
	}
	if d := a.Str("description"); d != "" {
		opts.Description = gl.Ptr(d)
	}
	if b := a.Str("target_branch"); b != "" {
		opts.TargetBranch = gl.Ptr(b)
	}
	if s := a.Str("state_event"); s != "" {
		opts.StateEvent = gl.Ptr(s)
	}
	if l := a.Labels(); l != "" {
		ll := gl.LabelOptions{l}
		opts.Labels = &ll
	}
	mr, _, err := client.MergeRequests.UpdateMergeRequest(a.ProjectID(), a.MrIID(), opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToMR(mr))
}

func mergeMR(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.AcceptMergeRequestOptions{}
	if m := a.Str("merge_commit_message"); m != "" {
		opts.MergeCommitMessage = gl.Ptr(m)
	}
	if v, ok := a.Bool("squash"); ok {
		opts.Squash = gl.Ptr(v)
	}
	if v, ok := a.Bool("should_remove_source_branch"); ok {
		opts.ShouldRemoveSourceBranch = gl.Ptr(v)
	}
	mr, _, err := client.MergeRequests.AcceptMergeRequest(a.ProjectID(), a.MrIID(), opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToMR(mr))
}

func approveMR(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	approval, _, err := client.MergeRequestApprovals.ApproveMergeRequest(a.ProjectID(), a.MrIID(), nil)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(approval)
}

func listMRNotes(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	notes, _, err := client.Notes.ListMergeRequestNotes(a.ProjectID(), a.MrIID(),
		&gl.ListMergeRequestNotesOptions{ListOptions: gl.ListOptions{PerPage: model.MaxPerPage}})
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToNotes(notes))
}

func createMRNote(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	note, _, err := client.Notes.CreateMergeRequestNote(a.ProjectID(), a.MrIID(),
		&gl.CreateMergeRequestNoteOptions{Body: gl.Ptr(a.Body())})
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToNote(note))
}

func createMRDiscussion(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	pos := &gl.PositionOptions{
		BaseSHA:      gl.Ptr(a.Str("base_sha")),
		HeadSHA:      gl.Ptr(a.Str("head_sha")),
		StartSHA:     gl.Ptr(a.Str("start_sha")),
		PositionType: gl.Ptr("text"),
	}
	if p := a.Str("new_path"); p != "" {
		pos.NewPath = gl.Ptr(p)
	}
	if p := a.Str("old_path"); p != "" {
		pos.OldPath = gl.Ptr(p)
	}
	if v := a.Int64("new_line"); v != 0 {
		pos.NewLine = gl.Ptr(v)
	}
	if v := a.Int64("old_line"); v != 0 {
		pos.OldLine = gl.Ptr(v)
	}
	disc, _, err := client.Discussions.CreateMergeRequestDiscussion(
		a.ProjectID(), a.MrIID(),
		&gl.CreateMergeRequestDiscussionOptions{
			Body:     gl.Ptr(a.Body()),
			Position: pos,
		},
	)
	if err != nil {
		return errResult(fmt.Errorf("create discussion: %w", err))
	}
	return jsonResult(disc)
}

func listMRDiffs(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	diffs, _, err := client.MergeRequests.ListMergeRequestDiffs(a.ProjectID(), a.MrIID(), nil)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(diffs)
}

func listMRCommits(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	commits, _, err := client.MergeRequests.GetMergeRequestCommits(a.ProjectID(), a.MrIID(), nil)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToCommitRefs(commits))
}

func listMRPipelines(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	pipelines, _, err := client.MergeRequests.ListMergeRequestPipelines(a.ProjectID(), a.MrIID())
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToPipelineRefs(pipelines))
}

func searchMRs(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.ListMergeRequestsOptions{
		Search:      gl.Ptr(a.Search()),
		State:       gl.Ptr(a.State("opened")),
		Scope:       gl.Ptr("all"),
		ListOptions: gl.ListOptions{Page: a.Page(), PerPage: model.DefaultPerPage},
	}
	mrs, resp, err := client.MergeRequests.ListMergeRequests(opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.NewPaged(model.ToMRRefsFromBasic(mrs), resp, a.Page()))
}
