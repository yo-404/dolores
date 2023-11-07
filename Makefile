VERSION=$(shell git tag --sort=-version:refname | head -1)
SHA=$(shell git rev-parse --short HEAD)
CMD=dolores

LDFLAGS=-X 'main.version=$(VERSION)' -X 'main.commit=$(SHA)'

.PHONY: setup build build_linux test run clean all

.DEFAULT_GOAL: default

default: build test

setup:
	mkdir -p bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.54.2

install:
	go install --ldflags="${LDFLAGS}" ./cmd/${CMD}/

lint: setup
	./bin/golangci-lint run

test: lint
	go test ./...

gomod:
	go mod tidy

build: gomod
	go build --ldflags="${LDFLAGS}" -o ./bin/${CMD} ./cmd/${CMD}/

gorelease_snapshot: build
	goreleaser release --snapshot --rm-dist

lint-fix: go-import-fmt
	@go mod tidy
	@go mod verify
	@golangci-lint run --timeout=10m --fix "./..."

go-import-fmt:
	@./hack/fmt-imports.sh