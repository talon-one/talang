SHELL := /bin/bash

build:
	go generate ./...
	go run interpreter/generate_docs.go

precommithook: build
	git add functions.md
	find . -name '*_allop.go' | xargs git add

	@go test -coverprofile cover.prof ./...; if [ $$? == 0 ]; then curl -o "build.svg" "https://img.shields.io/badge/build-passing-green.svg"; else curl -o "build.svg" "https://img.shields.io/badge/build-failing-red.svg"; fi;
	
	curl -o "coverage.svg" $(shell go run generate_badges.go -profile cover.prof)

	git add build.svg coverage.svg

test: build
	go test -cover ./...

check: test
	gometalinter --vendor --enable-all --disable=lll ./...
