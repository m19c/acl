all: test clean build lint

clean:
	rm -f dist/*

lint:
	golangci-lint -c .golangci.yml run ./...

build:
	go build

test:
	go test ./... -coverprofile dist/cp.out

format:
	gofumpt -w -lang 1.20 ./.

.PHONY: build lint
.DEFAULT_GOAL: all