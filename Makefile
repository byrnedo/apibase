.PHONY: all build buildstatic buildstatic-in-docker build-dockerfile install-deps proto build-docker-image test help

default: help

build:					## Builds a dynamic linked binary
	go version
	go build

deps:					## Installs dependencies using glide
	@glide install

test:							## Run all tests
	@go test $(shell go list ./... | grep -v /vendor/)

vet:
	@go vet $(shell go list ./...|grep -v /vendor/)

fmt:
	@go fmt $(shell go list ./...|grep -v /vendor/)

help:							## Show this help.
	@grep -e "^[a-zA-Z_-]*:" Makefile|awk -F'##' '{gsub(/[ \t]+$$/, "", $$1);printf "%-30s\t%s\n", $$1, $$2}'

