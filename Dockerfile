FROM golang:1.25.8 AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /gitlab-mcp ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /gitlab-mcp /gitlab-mcp

EXPOSE 3000

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/gitlab-mcp", "--healthcheck"]

USER nonroot:nonroot
ENTRYPOINT ["/gitlab-mcp"]
