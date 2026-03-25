package tools

import (
	"context"

	"github.com/atla-digital/mcp-gitlab-lite/internal/model"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/args"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/descriptions"
	"github.com/mark3labs/mcp-go/mcp"
	gl "gitlab.com/gitlab-org/api/client-go"
)

func ProjectEntries(d descriptions.Catalog) []Entry {
	return tag(CatProjects,
		Entry{
			Tool: d.For("list_projects").
				Str("search").
				Str("membership").
				Str("visibility").
				Num("page").
				Build(),
			Handler: listProjects,
		},
		Entry{
			Tool:    d.For("get_project").Str("project_id", mcp.Required()).Build(),
			Handler: getProject,
		},
		Entry{
			Tool: d.For("list_group_projects").
				Str("group_id", mcp.Required()).
				Str("search").
				Bool("include_subgroups").
				Num("page").
				Build(),
			Handler: listGroupProjects,
		},
	)
}

func listProjects(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	membership := a.Str("membership") != "false"
	opts := &gl.ListProjectsOptions{
		Membership:  gl.Ptr(membership),
		ListOptions: gl.ListOptions{Page: a.Page(), PerPage: model.DefaultPerPage},
	}
	if s := a.Search(); s != "" {
		opts.Search = gl.Ptr(s)
	}
	if v := a.Str("visibility"); v != "" {
		vv := gl.VisibilityValue(v)
		opts.Visibility = &vv
	}
	projects, resp, err := client.Projects.ListProjects(opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.NewPaged(model.ToProjectRefs(projects), resp, a.Page()))
}

func getProject(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	p, _, err := client.Projects.GetProject(a.ProjectID(), &gl.GetProjectOptions{})
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToProject(p))
}

func listGroupProjects(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.ListGroupProjectsOptions{
		ListOptions: gl.ListOptions{Page: a.Page(), PerPage: model.DefaultPerPage},
	}
	if s := a.Search(); s != "" {
		opts.Search = gl.Ptr(s)
	}
	if v, ok := a.Bool("include_subgroups"); ok {
		opts.WithShared = gl.Ptr(v)
	}
	projects, resp, err := client.Groups.ListGroupProjects(a.GroupID(), opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.NewPaged(model.ToProjectRefs(projects), resp, a.Page()))
}
