SHELL := /bin/bash

build:
	@go get golang.org/x/tools/cmd/stringer
	go generate ./...
	go run interpreter/generate_docs.go -dir=./docs/

precommithook: build
	git add docs/functions.md
	find . -name '*_allop.go' | xargs git add

	@go test -coverprofile cover.prof ./...; if [ $$? == 0 ]; then curl -o "build.svg" "https://img.shields.io/badge/build-passing-brightgreen.svg"; else curl -o "build.svg" "https://img.shields.io/badge/build-failing-red.svg"; fi;
	@sleep 1s
	curl -o "coverage.svg" $(shell go run generate_coverage_badge.go -profile cover.prof)

	git add build.svg coverage.svg

test: build
	go test -race -count=1 -cover ./...

check: test
	gometalinter --vendor --enable-all --disable=lll ./...
