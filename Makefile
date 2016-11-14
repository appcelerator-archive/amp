
.PHONY: all clean build build-cli build-cli-linux build-cli-darwin build-cli-windows build-server build-server-linux build-server-darwin build-server-windows dist-linux dist-darwin dist-windows dist build-agent build-log-worker install install-server install-cli install-agent install-log-worker fmt simplify check version build-image run
.PHONY: test

SHELL := /bin/bash
BASEDIR := $(shell echo $${PWD})

VERSION_FILE=VERSION
# build variables (provided to binaries by linker LDFLAGS below)
VERSION := $(shell cat $(VERSION_FILE))
BUILD := $(shell git rev-parse HEAD | cut -c1-8)

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# ignore vendor directory for go files
SRC := $(shell find . -type f -name '*.go' -not -path './vendor/*' -not -path './.git/*')

# for walking directory tree (like for proto rule)
EXCLUDE_FILES_FILTER := -not -path './vendor/*' -not -path './.git/*' -not -path './.glide/*'
EXCLUDE_DIRS_FILTER := $(EXCLUDE_FILES_FILTER) -not -path '.' -not -path './vendor' -not -path './.git' -not -path './.glide'

# for tests
UNIT_TEST_PACKAGES        := $(shell find .       -type f -name '*_test.go' -not -path './tests/*' $(EXCLUDE_DIRS_FILTER) -exec dirname {} \; | sort -u)
INTEGRATION_TEST_PACKAGES := $(shell find ./tests -type f -name '*_test.go'                        $(EXCLUDE_DIRS_FILTER) -exec dirname {} \; | sort -u)

DIRS = $(shell find . -type d $(EXCLUDE_DIRS_FILTER))

# generated file dependencies for proto rule
PROTOFILES = $(shell find . -type f -name '*.proto' $(EXCLUDE_DIRS_FILTER))

# generated files that can be cleaned
GENERATED := $(shell find . -type f -name '*.pb.go' $(EXCLUDE_FILES_FILTER))

# ignore generated files when formatting/linting/vetting
CHECKSRC := $(shell find . -type f -name '*.go' -not -name '*.pb.go' $(EXCLUDE_FILES_FILTER))

OWNER := appcelerator
REPO := github.com/$(OWNER)/amp

CMDDIR := cmd
CLI := amp
SERVER := amplifier
AGENT := amp-agent
LOGWORKER := amp-log-worker

TAG := latest
IMAGE := $(OWNER)/amp:$(TAG)

# tools
# need UID:GID because files created by containerized tools when mounting
# cwd are set to root:root
UG := $(shell echo "$$(id -u $${USER}):$$(id -g $${USER})")

DOCKER_RUN := docker run -t --rm -u $(UG)

GOTOOLS := appcelerator/gotools2
GOOS := $(shell uname | tr [:upper:] [:lower:])
GOARCH := amd64
GO := $(DOCKER_RUN) --name go -v $${HOME}/.ssh:/root/.ssh -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) -e GOOS=$(GOOS) -e GOARCH=$(GOARCH) $(GOTOOLS) go
GOTEST := $(DOCKER_RUN) --name go -v $${HOME}/.ssh:/root/.ssh -v $${GOPATH}/bin:/go/bin -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) $(GOTOOLS) go test -v

GLIDE_DIRS := $${HOME}/.glide $${PWD}/.glide vendor
GLIDE := $(DOCKER_RUN) -u $(UG) -v $${HOME}/.ssh:/root/.ssh -v $${HOME}/.glide:/root/.glide -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) $(GOTOOLS) glide $${GLIDE_OPTS}
GLIDE_INSTALL := $(GLIDE) install
GLIDE_UPDATE := $(GLIDE) update

all: version check build

arch:
	@echo $(GOOS)

version:
	@echo "version: $(VERSION) (build: $(BUILD))"

clean:
	@rm -rf $(GENERATED)
	@rm -f $$(which $(CLI)) ./$(CLI)
	@rm -f $$(which $(SERVER)) ./$(SERVER)
	@rm -f coverage.out coverage-all.out

install-deps:
	@$(GLIDE_INSTALL)

update-deps:
	@$(GLIDE_UPDATE)

install: install-cli install-server install-agent install-log-worker

install-cli: proto
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLI)

install-server: proto
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(SERVER)

install-agent: proto
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(AGENT)

install-log-worker: proto
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(LOGWORKER)

build: build-cli build-server build-agent build-log-worker

build-cli: proto
	@hack/build $(CLI)

build-server: proto
	@hack/build $(SERVER)

build-agent: proto
	@hack/build $(AGENT)

build-log-worker: proto
	@hack/build $(LOGWORKER)

build-server-image:
	@docker build -t appcelerator/$(SERVER):$(TAG) .

build-cli-linux:
	@rm -f $(CLI)
	@env GOOS=linux GOARCH=amd64 VERSION=$(VERSION) hack/build $(CLI)

build-cli-darwin:
	@rm -f $(CLI)
	@env GOOS=darwin GOARCH=amd64 VERSION=$(VERSION) hack/build $(CLI)

build-cli-windows:
	@rm -f $(CLI).exe
	@env GOOS=windows GOARCH=amd64 VERSION=$(VERSION) hack/build $(CLI)

build-server-linux:
	@rm -f $(SERVER)
	@env GOOS=linux GOARCH=amd64 VERSION=$(VERSION) hack/build $(SERVER)

build-server-darwin:
	@rm -f $(SERVER)
	@env GOOS=darwin GOARCH=amd64 VERSION=$(VERSION) hack/build $(SERVER)

build-server-windows:
	@rm -f $(SERVER).exe
	@env GOOS=windows GOARCH=amd64 VERSION=$(VERSION) hack/build $(SERVER)

dist-linux: build-cli-linux build-server-linux
	@rm -f dist/Linux/x86_64/amp-$(VERSION).tgz
	@mkdir -p dist/Linux/x86_64
	@tar czf dist/Linux/x86_64/amp-$(VERSION).tgz $(CLI) $(SERVER)

dist-darwin: build-cli-darwin build-server-darwin
	@rm -f dist/Darwin/x86_64/amp-$(VERSION).tgz
	@mkdir -p dist/Darwin/x86_64
	@tar czf dist/Darwin/x86_64/amp-$(VERSION).tgz $(CLI) $(SERVER)
	
dist-windows: build-cli-windows build-server-windows
	@rm -f dist/Windows/x86_64/amp-$(VERSION).zip
	@mkdir -p dist/Windows/x86_64
	@zip -q dist/Windows/x86_64/amp-$(VERSION).zip $(CLI).exe $(SERVER).exe
	
dist: dist-linux dist-darwin dist-windows
	
proto: $(PROTOFILES)
	@go run hack/proto.go

# used to install when you're already inside a container
install-host: proto-host
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLI)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(SERVER)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(AGENT)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(LOGWORKER)

# used to run protoc when you're already inside a container
proto-host: $(PROTOFILES)
	@go run hack/proto.go -protoc

# format and simplify if possible (https://golang.org/cmd/gofmt/#hdr-The_simplify_command)
fmt:
	@gofmt -s -l -w $(CHECKSRC)

check:
	@test -z $(shell gofmt -l ${CHECKSRC} | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@$(DOCKER_RUN) -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) $(GOTOOLS) bash -c 'for p in $$(go list ./... | grep -v /vendor/); do golint $${p} | sed "/pb\.go/d"; done'
	@go tool vet ${CHECKSRC}

build-image:
	@docker build -t $(IMAGE) .

run: build-image
	@CID=$(shell docker run --net=host -d --name $(SERVER) $(IMAGE)) && echo $${CID}

test-unit:
	@for pkg in $(UNIT_TEST_PACKAGES) ; do \
		go test $$pkg ; \
	done

test-integration:
	@for pkg in $(INTEGRATION_TEST_PACKAGES) ; do \
		go test $$pkg ; \
	done

test: test-unit test-integration

cover:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(TEST_PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out
