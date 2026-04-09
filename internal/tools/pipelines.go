package tools

import (
	"context"
	"encoding/json"
	"time"

	glclient "github.com/atla-digital/mcp-gitlab-lite/internal/gitlab"
	"github.com/atla-digital/mcp-gitlab-lite/internal/model"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/args"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/descriptions"
	"github.com/mark3labs/mcp-go/mcp"
	gl "gitlab.com/gitlab-org/api/client-go"
)

var terminalStates = map[string]bool{
	"success": true, "failed": true, "canceled": true,
	"skipped": true, "manual": true,
}

func PipelineEntries(d descriptions.Catalog) []Entry {
	return tag(CatPipelines,
		Entry{
			Tool: d.For("list_pipelines").
				Str("project_id", mcp.Required()).
				Str("status").Str("ref").Num("page").
				Build(),
			Handler: listPipelines,
		},
		Entry{
			Tool:    d.For("get_pipeline").Str("project_id", mcp.Required()).Num("pipeline_id", mcp.Required()).Build(),
			Handler: getPipeline,
		},
		Entry{
			Tool:    d.For("create_pipeline").Str("project_id", mcp.Required()).Str("ref", mcp.Required()).Str("variables").Build(),
			Handler: createPipeline,
		},
		Entry{
			Tool:    d.For("cancel_pipeline").Str("project_id", mcp.Required()).Num("pipeline_id", mcp.Required()).Build(),
			Handler: cancelPipeline,
		},
		Entry{
			Tool:    d.For("retry_pipeline").Str("project_id", mcp.Required()).Num("pipeline_id", mcp.Required()).Build(),
			Handler: retryPipeline,
		},
		Entry{
			Tool: d.For("wait_pipeline_change").
				Str("project_id", mcp.Required()).
				Num("pipeline_id", mcp.Required()).
				Num("timeout_seconds").
				Num("poll_interval_seconds").
				Build(),
			Handler: waitPipelineChange,
		},
		Entry{
			Tool: d.For("wait_pipeline_finish").
				Str("project_id", mcp.Required()).
				Num("pipeline_id", mcp.Required()).
				Num("timeout_seconds").
				Build(),
			Handler: waitPipelineFinish,
		},
	)
}

func listPipelines(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.ListProjectPipelinesOptions{
		ListOptions: gl.ListOptions{Page: a.Page(), PerPage: model.DefaultPerPage},
	}
	if s := a.Str("status"); s != "" {
		v := gl.BuildStateValue(s)
		opts.Status = &v
	}
	if r := a.Ref(); r != "" {
		opts.Ref = gl.Ptr(r)
	}
	ps, resp, err := client.Pipelines.ListProjectPipelines(a.ProjectID(), opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.NewPaged(model.ToPipelineRefs(ps), resp, a.Page()))
}

func getPipeline(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	p, _, err := client.Pipelines.GetPipeline(a.ProjectID(), a.PipelineID())
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToPipeline(p))
}

func createPipeline(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	opts := &gl.CreatePipelineOptions{Ref: gl.Ptr(a.Ref())}
	if varsJSON := a.Str("variables"); varsJSON != "" {
		var m map[string]string
		if err := json.Unmarshal([]byte(varsJSON), &m); err == nil {
			vars := make([]*gl.PipelineVariableOptions, 0, len(m))
			for k, v := range m {
				vars = append(vars, &gl.PipelineVariableOptions{Key: gl.Ptr(k), Value: gl.Ptr(v)})
			}
			opts.Variables = &vars
		}
	}
	p, _, err := client.Pipelines.CreatePipeline(a.ProjectID(), opts)
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToPipeline(p))
}

func cancelPipeline(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	p, _, err := client.Pipelines.CancelPipelineBuild(a.ProjectID(), a.PipelineID())
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToPipeline(p))
}

func retryPipeline(_ context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	p, _, err := client.Pipelines.RetryPipelineBuild(a.ProjectID(), a.PipelineID())
	if err != nil {
		return errResult(err)
	}
	return jsonResult(model.ToPipeline(p))
}

func waitPipelineChange(ctx context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	pid, plid := a.ProjectID(), a.PipelineID()

	initial, _, err := glclient.RetryOnRateLimit(func() (*gl.Pipeline, *gl.Response, error) {
		return client.Pipelines.GetPipeline(pid, plid)
	})
	if err != nil {
		return errResult(err)
	}
	initialStatus := initial.Status

	timeout := time.After(time.Duration(a.Timeout(300)) * time.Second)
	ticker := time.NewTicker(time.Duration(a.PollInterval(3, 5)) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return errResult(ctx.Err())
		case <-timeout:
			current, _, err := glclient.RetryOnRateLimit(func() (*gl.Pipeline, *gl.Response, error) {
				return client.Pipelines.GetPipeline(pid, plid)
			})
			if err != nil {
				return errResult(err)
			}
			return jsonResult(map[string]any{
				"changed":         false,
				"timed_out":       true,
				"previous_status": initialStatus,
				"current_status":  current.Status,
				"pipeline":        model.ToPipeline(current),
			})
		case <-ticker.C:
			current, _, err := glclient.RetryOnRateLimit(func() (*gl.Pipeline, *gl.Response, error) {
				return client.Pipelines.GetPipeline(pid, plid)
			})
			if err != nil {
				continue
			}
			if current.Status != initialStatus {
				return jsonResult(map[string]any{
					"changed":         true,
					"timed_out":       false,
					"previous_status": initialStatus,
					"current_status":  current.Status,
					"pipeline":        model.ToPipeline(current),
				})
			}
		}
	}
}

func waitPipelineFinish(ctx context.Context, client *gl.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a := args.From(req)
	pid, plid := a.ProjectID(), a.PipelineID()

	poll := func() (*mcp.CallToolResult, bool) {
		p, _, err := glclient.RetryOnRateLimit(func() (*gl.Pipeline, *gl.Response, error) {
			return client.Pipelines.GetPipeline(pid, plid)
		})
		if err != nil {
			return nil, false
		}
		if terminalStates[p.Status] {
			res, _ := jsonResult(map[string]any{
				"finished":  true,
				"timed_out": false,
				"status":    p.Status,
				"pipeline":  model.ToPipeline(p),
			})
			return res, true
		}
		return nil, false
	}

	if res, done := poll(); done {
		return res, nil
	}

	timeout := time.After(time.Duration(a.Timeout(1800)) * time.Second)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return errResult(ctx.Err())
		case <-timeout:
			p, _, err := glclient.RetryOnRateLimit(func() (*gl.Pipeline, *gl.Response, error) {
				return client.Pipelines.GetPipeline(pid, plid)
			})
			if err != nil {
				return errResult(err)
			}
			return jsonResult(map[string]any{
				"finished":  false,
				"timed_out": true,
				"status":    p.Status,
				"pipeline":  model.ToPipeline(p),
			})
		case <-ticker.C:
			if res, done := poll(); done {
				return res, nil
			}
		}
	}
}
