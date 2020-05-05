export GO111MODULE=on
-include .env


VERSION?=1.1.0
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
SDK_INFO?= $(shell go version)

WIN_OS?=windows
LINUX_OS?=linux
DOCKER_OS?=docker
MAC_OS?=darwin
APP?=NFC CLI
APP_ALIAS?=nfc-cli

export

all: deps lint test

build:
	APP_ID=$(APP_ID) CERT=$(CERT) SECRET=$(SECRET) goreleaser --snapshot --skip-publish --rm-dist

publish:
	APP_ID=$(APP_ID) CERT=$(CERT) SECRET=$(SECRET) goreleaser --rm-dist


build-windows:
	GOOS=$(WIN_OS) go build -mod=vendor -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME} -X main.Platform=${WIN_OS} -X 'main.SDKInfo=$(SDK_INFO)' -X 'main.AppID=${APP_ID}' -X 'main.AppCert=${CERT}' -X 'main.AppSecret=${SECRET}'" -o ${APP_ALIAS}.exe ./main.go

build-mac:
	GOOS=$(MAC_OS) go build -mod=vendor -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME} -X main.Platform=${MAC_OS} -X 'main.SDKInfo=$(SDK_INFO)'" -o ${APP_ALIAS} ./main.go

build-linux:
	 GOOS=$(LINUX_OS) go build -mod=vendor -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME} -X main.Platform=${LINUX_OS} -X 'main.SDKInfo=$(SDK_INFO)'" -o ${APP_ALIAS} ./main.go

lint:
	golangci-lint run ./...

test:
	go test -count=1 -cover -v `go list ./...`

deps:
	go mod download
	go mod vendor
