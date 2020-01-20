#!/usr/bin/make
 
.DEFAULT_GOAL := all

# Development
# ---------------------------------------------------------------------------------

# Install the dependencies of build

setup:
	@go get golang.org/x/lint/golint
	@go get golang.org/x/tools/cmd/goimports
	@go get github.com/securego/gosec/cmd/gosec

# Check quality of code

GOFILES	= $(shell find . -type f -name '*.go' -not -path "./.git/*")

fmt:
	$(eval FMT_LOG := $(shell mktemp -t gofmt.XXXXX))
	@gofmt -d -s -e $(GOFILES) > $(FMT_LOG) || true
	@[ ! -s "$(FMT_LOG)" ] || (echo "gofmt failed:" | cat - $(FMT_LOG) && false)
	
imports:
	$(eval IMP_LOG := $(shell mktemp -t goimp.XXXXX))
	@$(GOPATH)/bin/goimports -d -e -l $(GOFILES) > $(IMP_LOG) || true
	@[ ! -s "$(IMP_LOG)" ] || (echo "goimports failed:" | cat - $(IMP_LOG) && false)

lint:
	@$(GOPATH)/bin/golint -set_exit_status $(shell go list ./...)

verify:
	@make -s fmt
	@make -s imports
	@make -s lint

# Inspect source code for security problems
sec:
	@gosec -quiet ./...

# Build the source code for current os and arch
build:
	@export GO111MODULE=on;go build

# Run the tests
test:
	@go test ./... -race -coverprofile=coverage.txt -covermode=atomic

# Run everything

all:
	@make -s build test verify sec

# Releases
# ---------------------------------------------------------------------------------

# Create tag
# TAG_NAME=v0.0.1 make tag
tag:
	git tag -a $(TAG_NAME) -m "kube-wide release - $(TAG_NAME)"
	git push origin $(TAG_NAME)

# Build the source code for many os and arch

DIST_DIR := dist
PLATFORMS := linux/amd64 darwin/amd64 windows/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

builds: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -o 'dist/kw_$(os)-$(arch)'
