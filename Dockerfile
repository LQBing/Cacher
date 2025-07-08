FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o cacher
FROM alpine:latest
WORKDIR /app
VOLUME [ "/app/save" ]
COPY --from=builder /app/cacher .
ENTRYPOINT ["./cacher"]
