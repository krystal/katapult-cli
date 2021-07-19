ROOT := $(CURDIR)
SOURCES := $(shell find . -name "*.go" -or -name "go.mod" -or -name "go.sum" \
	-or -name "Makefile")

# Verbose output
ifdef VERBOSE
V = -v
endif

#
# Environment
#

BINDIR := bin
TOOLDIR := $(BINDIR)/tools

# Global environment variables for all targets
SHELL ?= /bin/bash
SHELL := env \
	GO111MODULE=on \
	GOBIN=$(CURDIR)/$(TOOLDIR) \
	CGO_ENABLED=0 \
	PATH='$(CURDIR)/$(BINDIR):$(CURDIR)/$(TOOLDIR):$(PATH)' \
	$(SHELL)

#
# Defaults
#

# Default target
.DEFAULT_GOAL := test

#
# Tools
#

# external tool
define tool # 1: binary-name, 2: go-import-path
TOOLS += $(TOOLDIR)/$(1)

$(TOOLDIR)/$(1): Makefile
	GOBIN="$(CURDIR)/$(TOOLDIR)" go install "$(2)"
endef

$(eval $(call tool,gofumpt,mvdan.cc/gofumpt@latest))
$(eval $(call tool,goimports,golang.org/x/tools/cmd/goimports@latest))
$(eval $(call tool,golangci-lint,github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41))
$(eval $(call tool,gomod,github.com/Helcaraxan/gomod@latest))
$(eval $(call tool,goreleaser,github.com/goreleaser/goreleaser@latest))

.PHONY: tools
tools: $(TOOLS)

#
# Build
#

LDFLAGS := -w -s

VERSION ?= $(shell git describe --tags)
DATE ?= $(shell date +%s)
GIT_SHA ?= $(shell git rev-parse --short HEAD)

ifndef VERSION
	VERSION = dev
endif

CMDDIR := cmd
BINS := $(shell test -d "$(CMDDIR)" && cd "$(CMDDIR)" && \
	find * -maxdepth 0 -type d -exec echo $(BINDIR)/{} \;)

.PHONY: build
build: $(BINS)

$(BINS): $(BINDIR)/%: $(SOURCES)
	mkdir -p "$(BINDIR)"
	cd "$(CMDDIR)/$*" && go build -a $(V) \
		-o "$(CURDIR)/$(BINDIR)/$*" \
		-ldflags "$(LDFLAGS) \
			-X main.version=$(VERSION) \
			-X main.commit=$(COMMIT) \
			-X main.date=$(DATE)"

.PHONY: build-snapshot
build-snapshot: $(TOOLDIR)/goreleaser
	goreleaser --snapshot --rm-dist

#
# Development
#

TEST ?= ./...

.PHONY: clean
clean:
	rm -rf $(TOOLDIR) $(BINDIR)
	rm -f ./coverage.out ./go.mod.tidy-check ./go.sum.tidy-check

.PHONY: test
test:
	CGO_ENABLED=1 go test $(V) -count=1 -race  $(TESTARGS) $(TEST)

.PHONY: test-deps
test-deps:
	go test all

.PHONY: lint
lint: $(TOOLDIR)/golangci-lint
	golangci-lint $(V) run

.PHONY: format
format: $(TOOLDIR)/goimports $(TOOLDIR)/gofumpt
	goimports -l -w .
	gofumpt -l -w .

#
# Coverage
#

.PHONY: cov
cov: coverage.out

.PHONY: cov-html
cov-html: coverage.out
	go tool cover -html=./coverage.out

.PHONY: cov-func
cov-func: coverage.out
	go tool cover -func=./coverage.out

coverage.out: $(SOURCES)
	go test $(V) -covermode=count -coverprofile=./coverage.out ./...

#
# Dependencies
#

.PHONY: deps
deps:
	$(info Downloading dependencies)
	go mod download


.PHONY: deps-analyze
deps-analyze: $(TOOLDIR)/gomod
	gomod analyze

.PHONY: tidy
tidy:
	go mod tidy $(V)

.PHONY: verify
verify:
	go mod verify

.SILENT: check-tidy
.PHONY: check-tidy
check-tidy:
	cp go.mod go.mod.tidy-check
	cp go.sum go.sum.tidy-check
	go mod tidy
	( \
		diff go.mod go.mod.tidy-check && \
		diff go.sum go.sum.tidy-check && \
		rm -f go.mod go.sum && \
		mv go.mod.tidy-check go.mod && \
		mv go.sum.tidy-check go.sum \
	) || ( \
		rm -f go.mod go.sum && \
		mv go.mod.tidy-check go.mod && \
		mv go.sum.tidy-check go.sum; \
		exit 1 \
	)

#
# Release
#

.PHONY: new-version
new-version: check-npx
	npx standard-version

.PHONY: next-version
next-version: check-npx
	npx standard-version --dry-run

.PHONY: check-npx
check-npx:
	$(if $(shell which npx),,\
		$(error No npx found in PATH, please install NodeJS))
