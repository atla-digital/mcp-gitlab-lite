package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	glclient "github.com/atla-digital/mcp-gitlab-lite/internal/gitlab"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools"
	"github.com/atla-digital/mcp-gitlab-lite/internal/tools/descriptions"
	"github.com/mark3labs/mcp-go/server"
)

const (
	version = "1.0.0"
	port    = "3000"
)

func main() {
	// Built-in health check for Docker HEALTHCHECK on distroless.
	if len(os.Args) > 1 && os.Args[1] == "--healthcheck" {
		runHealthCheck()
		return
	}

	s := server.NewMCPServer("gitlab-mcp", version,
		server.WithToolCapabilities(true),
	)

	d := descriptions.Load()
	tools.Register(s, d)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"status":"ok","version":"%s"}`, version)
	})

	httpServer := server.NewStreamableHTTPServer(s,
		server.WithStateLess(true),
		server.WithHTTPContextFunc(glclient.InjectHeaders),
		server.WithStreamableHTTPServer(&http.Server{
			Addr:              ":" + port,
			Handler:           mux,
			ReadHeaderTimeout: 10 * time.Second,
		}),
	)

	mux.Handle("/mcp", httpServer)

	fmt.Printf("gitlab-mcp v%s listening on :%s\n", version, port)
	if err := httpServer.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}

// runHealthCheck hits the /health endpoint and exits 0 on success, 1 on failure.
// Used by Docker HEALTHCHECK on distroless images (no curl/wget available).
func runHealthCheck() {
	client := &http.Client{Timeout: 3 * time.Second}

	resp, err := client.Get("http://localhost:" + port + "/health")
	if err != nil {
		os.Exit(1)
	}

	status := resp.StatusCode
	_ = resp.Body.Close()

	if status != http.StatusOK {
		os.Exit(1)
	}
}
