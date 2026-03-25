package glclient

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	gl "gitlab.com/gitlab-org/api/client-go"
)

type contextKey string

const (
	tokenKey contextKey = "gitlab-token"
	urlKey   contextKey = "gitlab-url"
)

// InjectHeaders is an HTTPContextFunc that extracts X-GitLab-Token and
// X-GitLab-URL from request headers and stores them in the context.
func InjectHeaders(ctx context.Context, r *http.Request) context.Context {
	if token := r.Header.Get("X-GitLab-Token"); token != "" {
		ctx = context.WithValue(ctx, tokenKey, token)
	}
	if url := r.Header.Get("X-GitLab-URL"); url != "" {
		ctx = context.WithValue(ctx, urlKey, url)
	}
	return ctx
}

// FromContext builds a *gl.Client from credentials stored in the context
// by InjectHeaders. X-GitLab-Token is required; X-GitLab-URL defaults
// to https://gitlab.com.
func FromContext(ctx context.Context) (*gl.Client, error) {
	token, _ := ctx.Value(tokenKey).(string)
	if token == "" {
		return nil, fmt.Errorf("X-GitLab-Token header is required")
	}

	baseURL, _ := ctx.Value(urlKey).(string)
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}

	baseURL = strings.TrimRight(baseURL, "/")

	return gl.NewClient(token, gl.WithBaseURL(baseURL+"/api/v4"))
}
