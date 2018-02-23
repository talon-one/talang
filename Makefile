SHELL := /bin/sh

build:
	go generate ./...
	go run interpreter/generate_docs.go

precommithook: test
	git add functions.md
	find . -name '*_allop.go' | xargs git add

	curl -o "coverage.svg" $(shell go run generate_badges.go -profile cover.prof)
	git add coverage.svg

test: build
	go test -coverprofile cover.prof ./...

check: test
	gometalinter --vendor --enable-all --disable=lll ./...
