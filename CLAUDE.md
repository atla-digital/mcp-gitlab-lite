# gitlab-mcp-lite

> **Keep this file up-to-date** as the project evolves. When adding tools, changing architecture, or modifying conventions, update the relevant sections here.

## Overview

A GitLab MCP (Model Context Protocol) server written in Go. Provides ~40 tools for AI assistants to interact with GitLab (projects, issues, merge requests, pipelines, jobs, repositories, users). Runs as a Streamable HTTP server on port 3000.

## Tech stack

- **Language**: Go 1.24+
- **MCP SDK**: `github.com/mark3labs/mcp-go`
- **GitLab client**: `gitlab.com/gitlab-org/api/client-go` (aliased as `gl`)
- **Transport**: Streamable HTTP (`server.NewStreamableHTTPServer`)

## Project layout

```
├── cmd/
│   ├── server/main.go          # HTTP entry point
│   └── gendocs/main.go         # Generates TOOLS.md
├── internal/
│   ├── gitlab/
│   │   ├── client.go           # FromRequest — builds *gl.Client
│   │   └── retry.go            # RetryOnRateLimit[T]
│   ├── model/
│   │   ├── types.go            # Domain types (Issue, MR, Pipeline, …)
│   │   └── convert.go          # gl.* → domain converters, Paged[T]
│   └── tools/
│       ├── handler.go          # GLHandler, Wrap, errResult, jsonResult
│       ├── registry.go         # Category type, tag(), All()
│       ├── register.go         # Register(s, d) — wires tools to MCP server
│       ├── args/args.go        # Typed argument accessors
│       ├── descriptions/
│       │   ├── loader.go       # embed.FS markdown parser → Catalog
│       │   ├── builder.go      # Fluent tool builder: d.For("name").Str(…).Build()
│       │   └── *.md            # Tool/param descriptions (one file per category)
│       ├── issues.go
│       ├── mergerequests.go
│       ├── pipelines.go
│       ├── jobs.go
│       ├── repos.go
│       ├── projects.go
│       └── users.go
├── .golangci.yml               # golangci-lint v2 config
├── Makefile
├── TOOLS.md                    # Auto-generated — do not edit manually
└── go.mod
```

## Build & run

```sh
make build          # compile binary
make run            # build + start server on :3000
make test           # run all tests
make generate       # regenerate TOOLS.md
make lint           # golangci-lint (gosec, gocritic, revive, staticcheck, …)
make fmt            # gofumpt formatting
make vuln           # govulncheck vulnerability scan
make check          # fmt + lint + test + vuln (full quality gate)
```

## Authentication

Credentials are passed per-request via HTTP headers (the server holds no credentials):

| Header | Required | Default | Description |
|--------|----------|---------|-------------|
| `X-GitLab-Token` | yes | — | GitLab personal access token |
| `X-GitLab-URL` | no | `https://gitlab.com` | GitLab instance base URL |

This allows a single server to connect to any number of GitLab instances with different tokens.

The server listens on port **3000** (hardcoded).

## Adding a new tool

1. Add the handler function in the appropriate `internal/tools/<category>.go` file.
2. Add the `Entry{}` to the category's `XxxEntries(d)` function using the builder pattern:
   ```go
   Entry{
       Tool: d.For("tool_name").Str("param", mcp.Required()).Build(),
       Handler: myHandler,
   },
   ```
3. Add descriptions in `internal/tools/descriptions/<category>.md`:
   ```markdown
   ## tool_name
   What this tool does.
   ### param
   What this parameter means.
   ```
4. Run `make test` — the description test will fail if any tool/param lacks a description.
5. Run `make generate` to update TOOLS.md.

## Conventions

- **Tool names**: `snake_case` (e.g. `list_issues`, `get_merge_request`)
- **Handler signature**: `func(ctx, *gl.Client, mcp.CallToolRequest) (*mcp.CallToolResult, error)`
- **Error returns**: always `errResult(err)` — never return a Go error from handlers
- **JSON responses**: always `jsonResult(v)` — marshals any value to MCP text result
- **Args access**: use `args.From(req)` then named accessors (`a.ProjectID()`, `a.Title()`, etc.)
- **Descriptions**: stored in embedded markdown, not in Go code
- **Categories**: typed `Category` constants — typos are compile errors
- **Imports**: alias GitLab client as `gl`, internal gitlab package as `glclient`
- **No backward compat hacks**: breaking changes are the default

## Versioning

Managed by [commitizen](https://commitizen-tools.github.io/commitizen/). Version is tracked in:
- `.cz.toml` (source of truth)
- `cmd/server/main.go` (compiled into binary)

Commit messages must follow [Conventional Commits](https://www.conventionalcommits.org/):
- `feat:` → MINOR bump
- `fix:` → PATCH bump
- `feat!:` or `BREAKING CHANGE:` → MAJOR bump

Bump version: `cz bump` (updates both files, creates tag + commit).

## Code quality tooling

All configured in `.golangci.yml` (v2 format) and runnable via `make`.

### golangci-lint v2 (`.golangci.yml`)
Linters enabled:
- **Essential**: govet, errcheck, staticcheck, unused, ineffassign
- **Security**: gosec
- **Style**: gocritic, revive, misspell, unconvert, unparam, nolintlint, copyloopvar
- **Bug prevention**: bodyclose, exhaustive
- **Performance**: prealloc

Formatters enabled: gofumpt, goimports

### gofumpt
Stricter superset of gofmt. Run `make fmt` or configure your editor to format on save.

### govulncheck
Scans for known vulnerabilities in dependencies and Go stdlib. Run `make vuln`.

### Workflow
Run `make check` before committing — it runs fmt, lint, test, and vuln in sequence.
Run `make lint` for a quick lint-only check.
