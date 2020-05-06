-include .env

SDK_INFO:=$(shell go version)
export SDK_INFO

all: deps lint test

build:
	APP_ID=$(APP_ID) CERT=$(CERT) SECRET=$(SECRET) goreleaser --snapshot --skip-publish --rm-dist

lint:
	golangci-lint run ./...

test:
	go test -count=1 -cover -v ./...

deps:
	go mod download
	go mod vendor
