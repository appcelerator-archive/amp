.PHONY: all clean install install-server install-cli fmt simplify check version build run

SHELL := /bin/bash
BASEDIR := $(shell echo $${PWD})

# build variables (provided to binaries by linker LDFLAGS below)
VERSION := 1.0.0
BUILD := $(shell git rev-parse HEAD | cut -c1-8)

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# ignore vendor directory for go files
SRC := $(shell find . -type f -name '*.go' -not -path './vendor/*')

# generated file dependencies for rpc rule
PROTODIR := api/rpc
PROTOFILES = $(shell find $(PROTODIR) -type f -name '*.proto' -not -path './vendor/*')

# generated files that can be cleaned
GENERATED := $(shell find $(PROTODIR) -type f -name '*.pb.go' -not -path './vendor/*')

# ignore generated files when formatting/linting/vetting
CHECKSRC := $(shell find . -type f -name '*.go' -not -name '*.pb.go' -not -path './vendor/*')

OWNER := appcelerator
REPO := github.com/$(OWNER)/amp

CMDDIR := cmd
CLI := amp
SERVER := amplifier

TAG := latest
IMAGE := $(OWNER)/amp/:$(TAG)

all: version check install

clean:
	@rm -rf $(GENERATED)
	@rm -f $$(which amp)

install: install-cli install-server

install-cli: rpc
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLI)

install-server: rpc
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(SERVER)

rpc: $(PROTOFILES)
	@for PKG in $$(ls -d $(BASEDIR)/$(PROTODIR)/*/); do cd $${PKG}; docker run -v $${PWD}:/go/src -v /var/run/docker.sock:/var/run/docker.sock appcelerator/protoc *.proto --go_out=plugins=grpc:.; done

# used to build under Docker
install-host: rpc-host
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLI)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(SERVER)

# used to build under Docker
rpc-host: $(PROTOFILES)
	@for PKG in $$(ls -d $(BASEDIR)/$(PROTODIR)/*/); do cd $${PKG}; protoc *.proto --go_out=plugins=grpc:.; done

fmt:
	@gofmt -l -w $(CHECKSRC)

simplify:
	@gofmt -s -l -w $(CHECKSRC)

check:
	@test -z $(shell gofmt -l ${CHECKSRC} | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d} | sed '/pb\.go/d'; done
	@go tool vet ${CHECKSRC}

version:
	@echo "version: $(VERSION) (build: $(BUILD))"

build:
	@docker build -t $(IMAGE) .

run: build
	@CID=$(shell docker run -d -p 50051:50051 --name $(SERVER) $(IMAGE)) && echo $${CID}
