// Package descriptions provides embedded tool/param descriptions loaded from markdown.
package descriptions

import (
	"embed"
	"strings"
)

//go:embed *.md
var fs embed.FS

// ToolDoc holds a tool's description and its parameter descriptions.
type ToolDoc struct {
	Description string
	Params      map[string]string
}

// Catalog maps tool names to their documentation.
type Catalog map[string]ToolDoc

// Load reads all embedded .md files and returns a Catalog.
func Load() Catalog {
	c := make(Catalog)
	entries, err := fs.ReadDir(".")
	if err != nil {
		return c
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		data, err := fs.ReadFile(e.Name())
		if err != nil {
			continue
		}
		for name, doc := range parse(data) {
			c[name] = doc
		}
	}
	return c
}

// Tool returns the description for the named tool, or "".
func (c Catalog) Tool(name string) string {
	if d, ok := c[name]; ok {
		return d.Description
	}
	return ""
}

// Param returns the description for a param of the named tool, or "".
func (c Catalog) Param(tool, param string) string {
	if d, ok := c[tool]; ok {
		return d.Params[param]
	}
	return ""
}

// parse extracts tool docs from markdown.
// ## headings define tools; ### headings define params.
func parse(data []byte) map[string]ToolDoc {
	result := make(map[string]ToolDoc)
	var currentTool string
	var currentParam string
	var buf strings.Builder

	flush := func() {
		text := strings.TrimSpace(buf.String())
		buf.Reset()
		if text == "" {
			return
		}
		if currentTool == "" {
			return
		}
		doc := result[currentTool]
		if currentParam != "" {
			if doc.Params == nil {
				doc.Params = make(map[string]string)
			}
			doc.Params[currentParam] = text
		} else {
			doc.Description = text
		}
		result[currentTool] = doc
	}

	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(trimmed, "### "):
			flush()
			currentParam = strings.TrimSpace(strings.TrimPrefix(trimmed, "### "))
		case strings.HasPrefix(trimmed, "## "):
			flush()
			currentTool = strings.TrimSpace(strings.TrimPrefix(trimmed, "## "))
			currentParam = ""
			if _, ok := result[currentTool]; !ok {
				result[currentTool] = ToolDoc{Params: make(map[string]string)}
			}
		default:
			if buf.Len() > 0 {
				buf.WriteByte('\n')
			}
			buf.WriteString(line)
		}
	}
	flush()
	return result
}
