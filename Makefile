export GO111MODULE=on

all: deps deps_check lint test build

build:
	go build -mod=vendor ./main.go

lint:
	golangci-lint run pkg/client

test:
	go test -count=1 -cover -v `go list ./...`

deps:
	go mod download
	go mod vendor
