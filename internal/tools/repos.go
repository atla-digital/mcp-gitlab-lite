package tools

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/atla-digital/mcp-gitlab-lite/internal/model"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/args"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/descriptions"
	"github.com/mark3labs/mcp-go/mcp"
	gl "gitlab.com/gitlab-org/api/client-go"
)

func RepoEntries(d descriptions.Catalog) []Entry {
	return tag(CatRepos,
		Entry{
			Tool:    d.For("list_branches").Str("project_id", mcp.Required()).Str("search").Build(),
			Handler: listBranches,
		},
		Entry{
			Tool:    d.For("get_branch").Str("project_id", mcp.Required()).Str("branch", mcp.Required()).Build(),
			Handler: getBranch,
		},
		Entry{
			Tool: d.For("create_branch").
				Str("project_id", mcp.Required()).
				Str("branch", mcp.Required()).
				Str("ref", mcp.Required()).
				Build(),
			Handler: createBranch,
		},
		Entry{
			Tool:    d.For("delete_branch").Str("project_id", mcp.Required()).Str("branch", mcp.Required()).Build(),
			Handler: deleteBranch,
		},
		Entry{
			Tool: d.For("list_repository_tree").
				Str("project_id", mcp.Required()).
				Str("path").Str("ref").Bool("recursive").
				Build(),
			Handler: listTree,
		},
		Entry{
			Tool: d.For("get_file_content").
				Str("project_id", mcp.Required()).
				Str("file_path", mcp.Required()).
				Str("ref").
				Build(),
			Handler: getFileContent,
		},
		Entry{
			Tool: d.For("list_commits").
				Str("project_id", mcp.Required()).
				Str("ref_name").Str("path").Num("page").
				Build(),
			Handler: listCommits,
		},
		Entry{
			Tool:    d.For("get_commit").Str("project_id", mcp.Required()).Str("sha", mcp.Required()).Build(),
			Handler: getCommit,
		},
		Entry{
			Tool: d.For("compare_refs").
				Str("project_id", mcp.Required()).
				Str("from", mcp.Required()).
				Str("to", mcp.Required()).
				Build(),
			Handler: compareRefs,
		},
	)
}

func listBranches(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.ListBranchesOptions{ListOptions: gl.ListOptions{PerPage: 50}}
	if s := a.Search(); s != "" {
		opts.Search = gl.Ptr(s)
	}
	branches, _, err := client.Branches.ListBranches(a.ProjectID(), opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToBranches(branches))
}

func getBranch(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	b, _, err := client.Branches.GetBranch(a.ProjectID(), a.Branch())
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToBranch(b))
}

func createBranch(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	b, _, err := client.Branches.CreateBranch(a.ProjectID(), &gl.CreateBranchOptions{
		Branch: gl.Ptr(a.Branch()),
		Ref:    gl.Ptr(a.Ref()),
	})
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToBranch(b))
}

func deleteBranch(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	if _, err := client.Branches.DeleteBranch(a.ProjectID(), a.Branch()); err != nil {
		return errResult(err)
	}
	return mcp.NewToolResultText(`{"deleted":true}`), nil
}

func listTree(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.ListTreeOptions{ListOptions: gl.ListOptions{PerPage: model.MaxPerPage}}
	if p := a.Str("path"); p != "" {
		opts.Path = gl.Ptr(p)
	}
	if r := a.Ref(); r != "" {
		opts.Ref = gl.Ptr(r)
	}
	if v, ok := a.Bool("recursive"); ok {
		opts.Recursive = gl.Ptr(v)
	}
	tree, _, err := client.Repositories.ListTree(a.ProjectID(), opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(tree)
}

func getFileContent(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.GetFileOptions{}
	if r := a.Ref(); r != "" {
		opts.Ref = gl.Ptr(r)
	}
	file, _, err := client.RepositoryFiles.GetFile(a.ProjectID(), a.FilePath(), opts)
	if err != nil {
		return errResult(err)
	}
	content, err := base64.StdEncoding.DecodeString(file.Content)
	if err != nil {
		return errResult(fmt.Errorf("base64 decode: %w", err))
	}
	return mcp.NewToolResultText(string(content)), nil
}

func listCommits(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.ListCommitsOptions{ListOptions: gl.ListOptions{Page: a.Page(), PerPage: model.DefaultPerPage}}
	if r := a.Str("ref_name"); r != "" {
		opts.RefName = gl.Ptr(r)
	}
	if p := a.Str("path"); p != "" {
		opts.Path = gl.Ptr(p)
	}
	commits, resp, err := client.Commits.ListCommits(a.ProjectID(), opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.NewPaged(model.ToCommitRefs(commits), resp, a.Page()))
}

func getCommit(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	c, _, err := client.Commits.GetCommit(a.ProjectID(), a.SHA(), nil)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToCommitRef(c))
}

func compareRefs(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	comparison, _, err := client.Repositories.Compare(a.ProjectID(), &gl.CompareOptions{
		From: gl.Ptr(a.Str("from")),
		To:   gl.Ptr(a.Str("to")),
	})
	if err != nil {
		return errResult(err)
	}
	return jsonResult(comparison)
}
