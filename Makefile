
.PHONY: all clean install fmt simplify check version build-image run
.PHONY: build build-cli-linux build-cli-darwin build-cli-windows build-server build-server-linux build-server-darwin build-server-windows
.PHONY: dist dist-linux dist-darwin dist-windows
.PHONY: test test-unit test-cli test-integration test-integration-host

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
UNIT_TEST_PACKAGES        := $(shell find .                   -type f -name '*_test.go' -not -path './tests/*' $(EXCLUDE_DIRS_FILTER) -exec dirname {} \; | sort -u)
CLI_TEST_PACKAGES         := $(shell find ./tests/cli         -type f -name '*_test.go'                        $(EXCLUDE_DIRS_FILTER) -exec dirname {} \; | sort -u)
INTEGRATION_TEST_PACKAGES := $(shell find ./tests/integration -type f -name '*_test.go'                        $(EXCLUDE_DIRS_FILTER) -exec dirname {} \; | sort -u)

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

# Binaries
CLI := amp
SERVER := amplifier
AGENT := amp-agent
LOGWORKER := amp-log-worker
GATEWAY := amplifier-gateway
CLUSTERSERVER := cluster-server
CLUSTERAGENT := cluster-agent
AMPCLUSTER := amp-cluster

TAG := latest
IMAGE := $(OWNER)/amp:$(TAG)

# tools
# need UID:GID because files created by containerized tools when mounting
# cwd are set to root:root
UG := $(shell echo "$$(id -u $${USER}):$$(id -g $${USER})")

DOCKER_RUN := docker run -t --rm -u $(UG)

GOTOOLS := appcelerator/gotools2:1.0.0
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

install-deps:
	@$(GLIDE_INSTALL)

update-deps:
	@$(GLIDE_UPDATE)

proto: $(PROTOFILES)
	@go run hack/proto.go

# used to run protoc when you're already inside a container
proto-host: $(PROTOFILES)
	@go run hack/proto.go -protoc

clean:
	@rm -rf $(GENERATED)
	@rm -f $$(which $(CLI)) ./$(CLI)
	@rm -f $$(which $(SERVER)) ./$(SERVER)
	@rm -f coverage.out coverage-all.out
	@rm -f $$(which $(AGENT)) ./$(AGENT)
	@rm -f $$(which $(LOGWORKER)) ./$(LOGWORKER)
	@rm -f $$(which $(GATEWAY)) ./$(GATEWAY)
	@rm -f $$(which $(CLUSTERSERVER)) ./$(CLUSTERSERVER)
	@rm -f $$(which $(CLUSTERAGENT)) ./$(CLUSTERAGENT)
	@rm -f $$(which $(AMPCLUSTER)) ./$(AMPCLUSTER)

install:
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLI)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(SERVER)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(AGENT)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(LOGWORKER)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(GATEWAY)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLUSTERSERVER)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLUSTERAGENT)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(AMPCLUSTER)

build:
	@hack/build $(CLI)
	@hack/build $(SERVER)
	@hack/build $(AGENT)
	@hack/build $(LOGWORKER)
	@hack/build $(GATEWAY)
	@hack/build $(CLUSTERAGENT)
	@hack/build $(CLUSTERSERVER)
	@hack/build $(AMPCLUSTER)

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

build-clustercli-linux:
	@rm -f $(CLI)
	@env GOOS=linux GOARCH=amd64 VERSION=$(VERSION) hack/build $(AMPCLUSTER)

build-clustercli-darwin:
	@rm -f $(CLI)
	@env GOOS=darwin GOARCH=amd64 VERSION=$(VERSION) hack/build $(AMPCLUSTER)

build-clustercli-windows:
	@rm -f $(CLI).exe
	@env GOOS=windows GOARCH=amd64 VERSION=$(VERSION) hack/build $(AMPCLUSTER)
	
dist-linux: build-cli-linux build-server-linux build-clustercli-linux
	@rm -f dist/Linux/x86_64/amp-$(VERSION).tgz
	@mkdir -p dist/Linux/x86_64
	@tar czf dist/Linux/x86_64/amp-$(VERSION).tgz $(CLI) $(SERVER) $(AMPCLUSTER)

dist-darwin: build-cli-darwin build-server-darwin build-clustercli-derwing
	@rm -f dist/Darwin/x86_64/amp-$(VERSION).tgz
	@mkdir -p dist/Darwin/x86_64
	@tar czf dist/Darwin/x86_64/amp-$(VERSION).tgz $(CLI) $(SERVER) $(AMPCLUSTER)

dist-windows: build-cli-windows build-server-windows build-clustercli-windows
	@rm -f dist/Windows/x86_64/amp-$(VERSION).zip
	@mkdir -p dist/Windows/x86_64
	@zip -q dist/Windows/x86_64/amp-$(VERSION).zip $(CLI).exe $(SERVER).exe $(AMPCLUSTER).exe

dist: dist-linux dist-darwin dist-windows

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

test-cli:
	@for pkg in $(CLI_TEST_PACKAGES) ; do \
		go test $$pkg ; \
	done

test-unit:
	@for pkg in $(UNIT_TEST_PACKAGES) ; do \
		go test $$pkg ; \
	done

test-integration:
	@docker service rm amp-integration-test > /dev/null 2>&1 || true
	@docker build -f Dockerfile.test -t appcelerator/amp-integration-test .
	@docker service create --network amp-infra --name amp-integration-test --restart-condition none appcelerator/amp-integration-test
	@containerid=""; \
	while [[ $${containerid} == "" ]] ; do \
		containerid=`docker ps -qf 'name=amp-integration'`; \
		sleep 1 ; \
	done; \
	docker logs -f $$containerid; \
	docker service rm amp-integration-test > /dev/null 2>&1 || true \
	exit `docker inspect --format='{{.State.ExitCode}}' $$containerid`

test-integration-host:
	@for pkg in $(INTEGRATION_TEST_PACKAGES) ; do \
		go test $$pkg ; \
	done

test: test-unit test-integration test-cli

cover:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(TEST_PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out
