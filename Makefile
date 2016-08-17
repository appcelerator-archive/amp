.PHONY: all clean install install-server install-cli fmt simplify check version build run test test-storage

SHELL := /bin/bash
BASEDIR := $(shell echo $${PWD})

# build variables (provided to binaries by linker LDFLAGS below)
VERSION := 1.0.0
BUILD := $(shell git rev-parse HEAD | cut -c1-8)

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# ignore vendor directory for go files
SRC := $(shell find . -type f -name '*.go' -not -path './vendor/*' -not -path './.git/*')

# for walking directory tree (like for proto rule)
DIRS = $(shell find . -type d -not -path '.' -not -path './vendor' -not -path './vendor/*' -not -path './.git' -not -path './.git/*')

# generated file dependencies for proto rule
PROTOFILES = $(shell find . -type f -name '*.proto' -not -path './vendor/*' -not -path './.git/*')

# generated files that can be cleaned
GENERATED := $(shell find . -type f -name '*.pb.go' -not -path './vendor/*' -not -path './.git/*')

# ignore generated files when formatting/linting/vetting
CHECKSRC := $(shell find . -type f -name '*.go' -not -name '*.pb.go' -not -path './vendor/*' -not -path './.git/*')

OWNER := appcelerator
REPO := github.com/$(OWNER)/amp

CMDDIR := cmd
CLI := amp
SERVER := amplifier

TAG := latest
IMAGE := $(OWNER)/amp:$(TAG)

all: version check install

version:
	@echo "version: $(VERSION) (build: $(BUILD))"

clean:
	@rm -rf $(GENERATED)
	@rm -f $$(which amp)

install: install-cli install-server

install-cli: proto
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLI)

install-server: proto
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(SERVER)

proto: $(PROTOFILES)
	@for DIR in $(DIRS); do cd $(BASEDIR)/$${DIR}; ls *.proto > /dev/null 2>&1 && docker run -v $${PWD}:/go/src -v /var/run/docker.sock:/var/run/docker.sock appcelerator/protoc *.proto --go_out=plugins=grpc:. || true; done

# used to build under Docker
install-host: proto-host
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLI)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(SERVER)

# used to build under Docker
proto-host: $(PROTOFILES)
	@for DIR in $(DIRS); do cd $(BASEDIR)/$${DIR}; ls *.proto > /dev/null 2>&1 && protoc *.proto --go_out=plugins=grpc:. || true; done

# format and simplify if possible (https://golang.org/cmd/gofmt/#hdr-The_simplify_command)
fmt:
	@gofmt -s -l -w $(CHECKSRC)

check:
	@test -z $(shell gofmt -l ${CHECKSRC} | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d} | sed '/pb\.go/d'; done
	@go tool vet ${CHECKSRC}

build:
	@docker build -t $(IMAGE) .

run: build
	@CID=$(shell docker run --net=host -d --name $(SERVER) $(IMAGE)) && echo $${CID}

install-deps:
	@glide install --strip-vcs --strip-vendor --update-vendored

update-deps:
	@glide update --strip-vcs --strip-vendor --update-vendored

test: test-storage
	@go test -v $(REPO)/data/influx
	@go test -v $(REPO)/api/rpc/logs
	@go test -v $(REPO)/api/rpc/project
	@go test -v $(REPO)/api/rpc/service
	@go test -v $(REPO)/api/rpc/stat

test-storage:
	@go test -v $(REPO)/data/storage/etcd
	
