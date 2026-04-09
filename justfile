set shell := ["bash", "-euo", "pipefail", "-c"]
set dotenv-load

default:
    @just --list

setup:
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4
    go install github.com/goreleaser/goreleaser/v2@v2.8.2
    go install github.com/zricethezav/gitleaks/v8@v8.24.0
    go install github.com/evilmartians/lefthook@v1.11.13
    go install golang.org/x/vuln/cmd/govulncheck@latest
    CGO_ENABLED=1 go install -tags extended github.com/gohugoio/hugo@v0.160.1
    lefthook install
    @echo "Done. All tools installed and hooks configured."

check: fmt lint test

fmt:
    gofmt -w .

lint:
    golangci-lint run

test:
    go test -race ./...

cover:
    go test -race -coverprofile=coverage.out ./...
    go tool cover -func=coverage.out

build:
    go build -o bin/nux .

install:
    go install .

schemas:
    go run ./internal/config/gen

docs:
    cd docs && hugo server -D

docs-build:
    cd docs && hugo --minify

cover-html:
    go test -race -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out

clean:
    rm -rf bin/ coverage.out docs/public/
