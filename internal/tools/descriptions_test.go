package tools_test

import (
	"testing"

	"github.com/atla-digital/mcp-gitlab-lite/internal/tools"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/descriptions"
)

func TestAllToolsHaveDescriptions(t *testing.T) {
	d := descriptions.Load()
	for _, e := range tools.All(d) {
		name := e.Tool.Name

		if d.Tool(name) == "" {
			t.Errorf("tool %q: missing tool description in .md files", name)
		}

		for pname := range e.Tool.InputSchema.Properties {
			if d.Param(name, pname) == "" {
				t.Errorf("tool %q param %q: missing description in .md files", name, pname)
			}
		}
	}
}

func TestCategoryCompleteness(t *testing.T) {
	d := descriptions.Load()
	for _, e := range tools.All(d) {
		if e.Category == "" {
			t.Errorf("tool %q has empty category", e.Tool.Name)
		}
	}
}

func TestNoDuplicateToolNames(t *testing.T) {
	d := descriptions.Load()
	seen := map[string]bool{}
	for _, e := range tools.All(d) {
		if seen[e.Tool.Name] {
			t.Errorf("duplicate tool name: %q", e.Tool.Name)
		}
		seen[e.Tool.Name] = true
	}
}
