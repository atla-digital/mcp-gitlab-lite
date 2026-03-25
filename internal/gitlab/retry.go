package glclient

import (
	"net/http"
	"strconv"
	"time"

	gl "gitlab.com/gitlab-org/api/client-go"
)

// RetryOnRateLimit executes fn; on HTTP 429 it sleeps Retry-After seconds
// (default 10s) and retries once. Subsequent 429s propagate.
func RetryOnRateLimit[T any](fn func() (T, *gl.Response, error)) (T, *gl.Response, error) {
	result, resp, err := fn()
	if err == nil || resp == nil || resp.StatusCode != http.StatusTooManyRequests {
		return result, resp, err
	}
	wait := 10 * time.Second
	if ra := resp.Header.Get("Retry-After"); ra != "" {
		if secs, perr := strconv.Atoi(ra); perr == nil {
			wait = time.Duration(secs) * time.Second
		}
	}
	time.Sleep(wait)
	return fn()
}
