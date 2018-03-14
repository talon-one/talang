SHELL := /bin/bash

build:
	@go get golang.org/x/tools/cmd/stringer
	go generate ./...
	go run interpreter/generate_docs.go -dir=./docs/

precommithook: build
	git add docs/functions.md
	find . -name '*_allop.go' | xargs git add

test: build
	go test -race -count=1 -cover ./...

metalint: test
	gometalinter --vendor --enable-all --disable=lll --exclude "_test\.go" --exclude "testhelpers" ./...

metalintall: test
	gometalinter --vendor --enable-all --disable=lll ./...