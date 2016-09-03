
.PHONY: all clean build build-cli build-server install install-server install-cli fmt simplify check version build-image run
.PHONY: test test-storage test-influx test-stat test-logs test-build test-project test-service

SHELL := /bin/bash
BASEDIR := $(shell echo $${PWD})

# build variables (provided to binaries by linker LDFLAGS below)
VERSION := 1.0.0
BUILD := $(shell git rev-parse HEAD | cut -c1-8)

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# ignore vendor directory for go files
SRC := $(shell find . -type f -name '*.go' -not -path './vendor/*' -not -path './.git/*')

# for walking directory tree (like for proto rule)
EXCLUDE_FILES_FILTER := -not -path './vendor/*' -not -path './.git/*' -not -path './.glide/*'
EXCLUDE_DIRS_FILTER := $(EXCLUDE_FILES_FILTER) -not -path '.' -not -path './vendor' -not -path './.git' -not -path './.glide'

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

TAG := latest
IMAGE := $(OWNER)/amp:$(TAG)

# tools
DOCKER_RUN := docker run -t --rm

GOTOOLS := appcelerator/gotools2
GOOS := $(shell uname | tr [:upper:] [:lower:])
GOARCH := amd64
GO := $(DOCKER_RUN) --name go -v $${HOME}/.ssh:/root/.ssh -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) -e GOOS=$(GOOS) -e GOARCH=$(GOARCH) $(GOTOOLS) go
GOTEST := $(DOCKER_RUN) --name go -v $${HOME}/.ssh:/root/.ssh -v $${GOPATH}/bin:/go/bin -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) $(GOTOOLS) go test -v

GLIDE_DIRS := $${HOME}/.glide $${PWD}/.glide vendor
GLIDE := $(DOCKER_RUN) -v $${HOME}/.ssh:/root/.ssh -v $${HOME}/.glide:/root/.glide -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) $(GOTOOLS) glide $${GLIDE_OPTS}
GLIDE_INSTALL := $(GLIDE) install
GLIDE_UPDATE := $(GLIDE) update

# need UID:GID because files created by containerized tools when mounting
# cwd are set to root:root
UG := $(shell echo "$$(id -u $${USER}):$$(id -g $${USER})")

all: version check build

arch:
	@echo $(GOOS)

version:
	@echo "version: $(VERSION) (build: $(BUILD))"

clean:
	@rm -rf $(GENERATED)
	@rm -f $$(which $(CLI)) ./$(CLI)
	@rm -f $$(which $(SERVER)) ./$(SERVER)

install-deps:
	@$(GLIDE_INSTALL)
	@chown -R $(UG) $(GLIDE_DIRS)

update-deps:
	@$(GLIDE_UPDATE)
	@chown -R $(UG) $(GLIDE_DIRS)

install: install-cli install-server

install-cli: proto
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLI)

install-server: proto
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(SERVER)

build: build-cli build-server

build-cli: proto
	@hack/build $(CLI)
	@chown -R $(UG) $(CLI)

build-server: proto
	@hack/build $(SERVER)
	@chown -R $(UG) $(SERVER)

proto: $(PROTOFILES)
	@for DIR in $(DIRS); do cd $(BASEDIR)/$${DIR}; ls *.proto > /dev/null 2>&1 && docker run --rm --name protoc -t -v $${PWD}:/go/src -v /var/run/docker.sock:/var/run/docker.sock appcelerator/protoc *.proto --go_out=plugins=grpc:. || true; done
	@find . -type f -name '*.pb.go' $(EXCLUDE_FILES_FILTER) | xargs chown $(UG)

# used to install when you're already inside a container
install-host: proto-host
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLI)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(SERVER)

# used to run protoc when you're already inside a container
proto-host: $(PROTOFILES)
	@for DIR in $(DIRS); do cd $(BASEDIR)/$${DIR}; ls *.proto > /dev/null 2>&1 && protoc *.proto --go_out=plugins=grpc:. || true; done

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

test: test-storage test-influx test-stat test-logs test-project test-build test-service

test-storage:
	@go $(REPO)/data/storage/etcd

test-influx:
	@go test -v $(REPO)/data/influx

test-stat:
#	@go test -v $(REPO)/api/rpc/stat

test-logs:
	@go test -v $(REPO)/api/rpc/logs

test-project:
	@go test -v $(REPO)/api/rpc/project

test-service:
	@go test -v $(REPO)/api/rpc/service

test-build:
	@go test -v $(REPO)/api/rpc/build
