package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	gl "gitlab.com/gitlab-org/api/client-go"
)

// fakeMergeAPI behaves like a GitLab project with require_sha_for_merge on: the merge is
// refused without a sha, and a sha other than the MR head is refused too.
func fakeMergeAPI(t *testing.T, headSHA string, gotSHA *string) *gl.Client {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v4/projects/web%2Fservices%2Fapp/merge_requests/7", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"iid": 7, "sha": headSHA, "state": "opened"})
	})
	mux.HandleFunc("PUT /api/v4/projects/web%2Fservices%2Fapp/merge_requests/7/merge", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			SHA string `json:"sha"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		*gotSHA = body.SHA
		switch body.SHA {
		case "":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"SHA must be provided when merging"}`))
		case headSHA:
			_ = json.NewEncoder(w).Encode(map[string]any{"iid": 7, "sha": headSHA, "state": "merged"})
		default:
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"message":"SHA does not match HEAD of source branch"}`))
		}
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	c, err := gl.NewClient("token", gl.WithBaseURL(srv.URL+"/api/v4"))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func mergeReq(args map[string]any) mcp.CallToolRequest {
	var req mcp.CallToolRequest
	req.Params.Arguments = args
	return req
}

func TestMergeMRSendsHeadSHAWhenNoneGiven(t *testing.T) {
	var got string
	c := fakeMergeAPI(t, "abc123", &got)
	res, err := mergeMR(context.Background(), c, mergeReq(map[string]any{
		"project_id": "web/services/app", "merge_request_iid": float64(7),
	}))
	if err != nil || res.IsError {
		t.Fatalf("merge failed: err=%v result=%+v", err, res)
	}
	if got != "abc123" {
		t.Errorf("sha sent = %q, want the MR head abc123", got)
	}
}

func TestMergeMRPassesGivenSHA(t *testing.T) {
	var got string
	c := fakeMergeAPI(t, "abc123", &got)
	res, err := mergeMR(context.Background(), c, mergeReq(map[string]any{
		"project_id": "web/services/app", "merge_request_iid": float64(7), "sha": "old999",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if got != "old999" {
		t.Errorf("sha sent = %q, want the given old999", got)
	}
	if !res.IsError {
		t.Error("a sha that is not the head must fail the merge, not be replaced")
	}
}
