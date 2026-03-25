package tools

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	glclient "github.com/atla-digital/mcp-gitlab-lite/internal/gitlab"
	"github.com/mark3labs/mcp-go/mcp"
	gl "gitlab.com/gitlab-org/api/client-go"
)

// GLHandler is the signature every tool handler must implement.
type GLHandler func(ctx context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error)

// ClientFactory abstracts client construction for testability.
type ClientFactory func(ctx context.Context) (*gl.Client, error)

// Wrap is the production adapter.
func Wrap(fn GLHandler) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return WrapWith(glclient.FromContext, fn)
}

// WrapWith injects a custom factory — the test seam.
func WrapWith(factory ClientFactory, fn GLHandler) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		start := time.Now()
		client, err := factory(ctx)
		if err != nil {
			slog.Warn("auth failed", "tool", req.Params.Name, "err", err)
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, callErr := fn(ctx, client, req)
		slog.Info("tool",
			"name", req.Params.Name,
			"ms", time.Since(start).Milliseconds(),
			"is_error", result != nil && result.IsError,
		)
		return result, callErr
	}
}

// errResult wraps a Go error as an MCP tool error result.
func errResult(err error) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError(err.Error()), nil
}

// jsonResult marshals v as JSON and returns it as an MCP text result.
func jsonResult(v any) (*mcp.CallToolResult, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return errResult(err)
	}
	return mcp.NewToolResultText(string(data)), nil
}
