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

# create outputdirectory
OUTDIR := ${CURDIR}/out
MAKE_OUTDIR = $(shell mkdir -p $(OUTDIR))

# create test directory
TSTDIR := ${CURDIR}/out/test
MAKE_TSTDIR = $(shell mkdir -p $(TSTDIR))

# create bin directory
BINDIR := $(GOPATH)/bin
MAKE_BINDIR = $(shell mkdir -p $(BINDIR))

# current timestamp
TIMESTAMP := $$(date +%Y%m%d%H%M%S)

.PHONY: all get-tools init build clean test install uninstall docs docs-html docs-pdf fmt simplify check

all: check test install

get-tools: $(MAKE_BINDIR)
	@go get -u golang.org/x/tools/cmd/goimports

init:
ifeq (, $(wildcard ./go.mod))
	@go mod init gitlab.com/skeeterhealth/espressopp
endif

# build: export GOPRIVATE=gitlab.com/skeeterhealth
build: init $(MAKE_OUTDIR)
	@go build $(LDFLAGS) -o $(OUTDIR) ./cmd/espressopp

clean:
	@rm -rf $(OUTDIR)

test: get-tools $(MAKE_TSTDIR)
	$(eval FILENAME := cover_$$(TIMESTAMP))
	@-go test -coverprofile=$(TSTDIR)/$(FILENAME).txt .
ifeq (, $(wildcard $(TSTDIR)/$(FILENAME).txt))
	@go tool cover -html=$(TSTDIR)/$(FILENAME).txt -o $(TSTDIR)/$(FILENAME).html
endif

install: init $(MAKE_BINDIR)
	@go install $(LDFLAGS) ./cmd/espressopp

uninstall:
	@rm -f $(BINDIR)/espressopp*

docs: docs-html docs-pdf

docs-html:
	@$(foreach file, $(wildcard $(DOCSDIR)/*.adoc), asciidoctor \
		-r asciidoctor-diagram \
	   	-D $(OUTDIR)/docs/html $(file);)

docs-pdf:
	@$(foreach file, $(wildcard $(DOCSDIR)/*.adoc), asciidoctor-pdf \
		-r asciidoctor-diagram \
	   	-D $(OUTDIR)/docs/pdf $(file);)

fmt: get-tools
	@gofmt -l -w $(SRC)
	@goimports -l -w $(SRC)

simplify:
	@gofmt -s -l -w $(SRC)

check:
	@test -z $(shell gofmt -l $(SRC) | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@go vet ${SRC}
