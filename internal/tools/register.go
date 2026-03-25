package tools

import (
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/descriptions"
	"github.com/mark3labs/mcp-go/server"
)

// Register adds all tools to the MCP server.
func Register(s *server.MCPServer, d descriptions.Catalog) {
	for _, e := range All(d) {
		s.AddTool(*e.Tool, Wrap(e.Handler))
	}
}
