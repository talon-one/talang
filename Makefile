SHELL := /bin/bash

build: docs
	go build -o ./cmd/talang-cli/talang-cli ./cmd/talang-cli

precommithook: docs
	git add docs/functions.md
	find . -name '*_allop.go' | xargs git add

docs: generate
	go run interpreter/generate_docs.go -dir=./docs/

generate:
	@go get golang.org/x/tools/cmd/stringer
	go generate ./...

test: generate
	go test -race -count=1 -cover ./...

metalint: test
	gometalinter --vendor --enable-all --disable=lll --exclude "_test\.go" --exclude "testhelpers" ./...

metalintall: test
	gometalinter --vendor --enable-all --disable=lll ./...