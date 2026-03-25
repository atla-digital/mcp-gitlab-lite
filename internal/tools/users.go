package tools

import (
	"context"

	"github.com/atla-digital/mcp-gitlab-lite/internal/model"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/args"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/descriptions"
	"github.com/mark3labs/mcp-go/mcp"
	gl "gitlab.com/gitlab-org/api/client-go"
)

func UserEntries(d descriptions.Catalog) []Entry {
	return tag(CatUsers,
		Entry{
			Tool:    d.For("get_user").Num("user_id", mcp.Required()).Build(),
			Handler: getUser,
		},
		Entry{
			Tool:    d.For("get_current_user").Build(),
			Handler: getCurrentUser,
		},
	)
}

func getUser(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	u, _, err := client.Users.GetUser(a.UserID(), gl.GetUsersOptions{})
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToUser(u))
}

func getCurrentUser(_ context.Context, client *gl.Client, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	u, _, err := client.Users.CurrentUser()
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToUser(u))
}
