package tools

import (
	"bufio"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/atla-digital/mcp-gitlab-lite/internal/model"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/args"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/descriptions"
	"github.com/mark3labs/mcp-go/mcp"
	gl "gitlab.com/gitlab-org/api/client-go"
)

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func JobEntries(d descriptions.Catalog) []Entry {
	return tag(CatJobs,
		Entry{
			Tool: d.For("list_pipeline_jobs").
				Str("project_id", mcp.Required()).
				Num("pipeline_id", mcp.Required()).
				Str("scope").
				Build(),
			Handler: listPipelineJobs,
		},
		Entry{
			Tool:    d.For("get_job").Str("project_id", mcp.Required()).Num("job_id", mcp.Required()).Build(),
			Handler: getJob,
		},
		Entry{
			Tool: d.For("get_job_log").
				Str("project_id", mcp.Required()).
				Num("job_id", mcp.Required()).
				Str("grep").
				Num("head").
				Num("tail").
				Build(),
			Handler: getJobLog,
		},
		Entry{
			Tool:    d.For("retry_job").Str("project_id", mcp.Required()).Num("job_id", mcp.Required()).Build(),
			Handler: retryJob,
		},
		Entry{
			Tool:    d.For("cancel_job").Str("project_id", mcp.Required()).Num("job_id", mcp.Required()).Build(),
			Handler: cancelJob,
		},
	)
}

func listPipelineJobs(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.ListJobsOptions{ListOptions: gl.ListOptions{PerPage: model.MaxPerPage}}
	if s := a.Scope(); s != "" {
		v := gl.BuildStateValue(s)
		opts.Scope = &[]gl.BuildStateValue{v}
	}
	jobs, _, err := client.Jobs.ListPipelineJobs(a.ProjectID(), a.PipelineID(), opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToJobs(jobs))
}

func getJob(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	job, _, err := client.Jobs.GetJob(a.ProjectID(), a.JobID())
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToJob(job))
}

func getJobLog(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	reader, _, err := client.Jobs.GetTraceFile(a.ProjectID(), a.JobID())
	if err != nil {
		return errResult(err)
	}

	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, ansiEscape.ReplaceAllString(scanner.Text(), ""))
	}
	if err := scanner.Err(); err != nil {
		return errResult(fmt.Errorf("reading trace: %w", err))
	}

	if pattern := a.Str("grep"); pattern != "" {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return errResult(fmt.Errorf("invalid grep pattern %q: %w", pattern, err))
		}
		filtered := lines[:0]
		for _, l := range lines {
			if re.MatchString(l) {
				filtered = append(filtered, l)
			}
		}
		lines = filtered
	}

	head, tail := a.Int("head"), a.Int("tail")
	switch {
	case head > 0 && len(lines) > head:
		lines = lines[:head]
	case tail > 0 && head == 0 && len(lines) > tail:
		lines = lines[len(lines)-tail:]
	}
	return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
}

func retryJob(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	job, _, err := client.Jobs.RetryJob(a.ProjectID(), a.JobID())
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToJob(job))
}

func cancelJob(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	job, _, err := client.Jobs.CancelJob(a.ProjectID(), a.JobID())
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToJob(job))
}
