build:
	go generate ./...
	go run interpreter/generate_docs.go

precommithook: build
	git add functions.md
	find . -name '*_allop.go' | xargs git add

test: build
	go test -cover ./...

check: test
	gometalinter --vendor --enable-all --disable=lll ./...
