GOBIN ?= ${PWD}
NEW_VERSION ?= 0.1.0
TARGET_OS ?= $(shell go env GOOS)
TARGET_ARCH ?= $(shell go env GOARCH)

.DEFAULT_GOAL := build

.PHONY: bootstrap
bootstrap:
	@GOBIN=${GOBIN} go get github.com/mitchellh/gox
	@go mod download

.PHONY: clean
clean:
	@rm -rf bin

.PHONY: alpha
alpha: clean
	@mkdir -p bin
	@go build -o bin -ldflags="-X main.applicationVersion=$(shell date +%s)-alpha"

.PHONY: build
build: clean
	@mkdir -p bin
	@${GOBIN}/gox -os=${TARGET_OS} -arch=${TARGET_ARCH} -output bin/ccp -ldflags="-X main.applicationVersion=${NEW_VERSION}"

.PHONY: local-version
local-version:
	@go run . version
