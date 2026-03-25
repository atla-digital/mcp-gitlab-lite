package descriptions

import "github.com/mark3labs/mcp-go/mcp"

// Builder constructs a *mcp.Tool for one named tool, auto-injecting
// descriptions from the Catalog for the tool and each parameter.
type Builder struct {
	name    string
	catalog Catalog
	opts    []mcp.ToolOption
}

// For starts a builder for the named tool.
func (c Catalog) For(name string) *Builder {
	return &Builder{
		name:    name,
		catalog: c,
		opts:    []mcp.ToolOption{mcp.WithDescription(c.Tool(name))},
	}
}

// Str adds a string parameter.
func (b *Builder) Str(param string, extra ...mcp.PropertyOption) *Builder {
	return b.add(mcp.WithString(param, b.props(param, extra)...))
}

// Num adds a number parameter.
func (b *Builder) Num(param string, extra ...mcp.PropertyOption) *Builder {
	return b.add(mcp.WithNumber(param, b.props(param, extra)...))
}

// Bool adds a boolean parameter.
func (b *Builder) Bool(param string, extra ...mcp.PropertyOption) *Builder {
	return b.add(mcp.WithBoolean(param, b.props(param, extra)...))
}

// Build returns the completed *mcp.Tool.
func (b *Builder) Build() *mcp.Tool {
	t := mcp.NewTool(b.name, b.opts...)
	return &t
}

// props prepends the catalog description to the caller's extra options.
func (b *Builder) props(param string, extra []mcp.PropertyOption) []mcp.PropertyOption {
	desc := b.catalog.Param(b.name, param)
	out := make([]mcp.PropertyOption, 0, 1+len(extra))
	if desc != "" {
		out = append(out, mcp.Description(desc))
	}
	return append(out, extra...)
}

func (b *Builder) add(opt mcp.ToolOption) *Builder {
	b.opts = append(b.opts, opt)
	return b
}
