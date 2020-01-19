#!/usr/bin/make
 
.DEFAULT_GOAL := all

GOFILES	= $(shell find . -type f -name '*.go' -not -path "./.git/*")

setup:
	@go get golang.org/x/lint/golint
	@go get golang.org/x/tools/cmd/goimports	

build:
	@export GO111MODULE=on;go build ./...

fmt:
	$(eval FMT_LOG := $(shell mktemp -t gofmt.XXXXX))
	@gofmt -d -s -e $(GOFILES) > $(FMT_LOG) || true
	@[ ! -s "$(FMT_LOG)" ] || (echo "gofmt failed:" | cat - $(FMT_LOG) && false)
	
imports:
	$(eval IMP_LOG := $(shell mktemp -t goimp.XXXXX))
	@goimports -d -e -l $(GOFILES) > $(IMP_LOG) || true
	@[ ! -s "$(IMP_LOG)" ] || (echo "goimports failed:" | cat - $(IMP_LOG) && false)

lint:
	@golint -set_exit_status $(shell go list ./...)

verify:
	@make -s fmt
	@make -s imports
	@make -s lint

test:
	@go test ./... -race -coverprofile=coverage.txt -covermode=atomic

all:
	@make -s build test verify
