SHELL := /bin/bash

.DEFAULT_GOAL := build

# version and build for the target
VERSION := 1.0.0
BUILD := `git rev-parse HEAD`

# use linker flags to provide version and build settings to the target
LDFLAGS = -ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go')

# documentation directory
DOCSDIR := ${CURDIR}/docs

# create build directory
OUTDIR := ${CURDIR}/out
MAKE_OUTDIR = $(shell mkdir -p $(OUTDIR))

# create bin directory
BINDIR := $(GOPATH)/bin
MAKE_BINDIR = $(shell mkdir -p $(BINDIR))

# default deploy target
DEPLOY_TARGET := remote

.PHONY: all get-tools init build clean test install uninstall docs fmt simplify check

all: check test install

get-tools: $(MAKE_BINDIR)
	@go get -u golang.org/x/tools/cmd/goimports

init:
ifeq (, $(wildcard ./go.mod))
	@go mod init github.com/skeeterhealth/espressopp
endif

build: init $(MAKE_OUTDIR)
	@go build $(LDFLAGS) -o $(OUTDIR) ./cmd/espressopp

clean:
	@rm -rf $(OUTDIR)

test:
	@go test ./...

install: init $(MAKE_BINDIR)
	@go install $(LDFLAGS) ./cmd/espressopp

uninstall:
	@rm -f $(BINDIR)/espressopp*

docs:
	@$(foreach file, $(wildcard $(DOCSDIR)/*.adoc), asciidoctor -D $(OUTDIR)/docs $(file);)

fmt: get-tools
	@gofmt -l -w $(SRC)
	@goimports -l -w $(SRC)

simplify:
	@gofmt -s -l -w $(SRC)

check:
	@test -z $(shell gofmt -l $(SRC) | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@go vet ${SRC}
