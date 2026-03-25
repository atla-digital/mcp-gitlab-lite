.PHONY: build run test generate lint fmt vuln check clean

BINARY := gitlab-mcp

build:
	go build -o $(BINARY) ./cmd/server

run: build
	./$(BINARY)

test:
	go test ./...

generate:
	go run ./cmd/gendocs

# Run golangci-lint (includes gosec, gocritic, revive, exhaustive, etc.)
lint:
	golangci-lint run ./...

# Format code with gofumpt (stricter superset of gofmt)
fmt:
	gofumpt -w .

# Scan for known vulnerabilities in dependencies and stdlib
vuln:
	govulncheck ./...

# Run all quality checks: format, lint, test, vuln scan
check: fmt lint test vuln

clean:
	rm -f $(BINARY)
