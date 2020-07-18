newversion := $(shell go run . version)
newversion ?= 0.1.0

.DEFAULT_GOAL := build

.PHONY: clean
clean:
	@rm -rf bin

.PHONY: alpha
alpha: clean
	@mkdir -p bin
	@go build -o bin -ldflags="-X main.applicationVersion=${newversion}-alpha"

.PHONY: build
build: clean
	@mkdir -p bin
	@go build -o bin -ldflags="-X main.applicationVersion=${newversion}"
