SHELL := /bin/bash
BASEDIR := $(shell echo $${PWD})

# =============================================================================
# BUILD MANAGEMENT
# Variables declared here are used by this Makefile *and* are exported to
# override default values used by supporting scripts in the hack directory
# =============================================================================
export UG := $(shell echo "$$(id -u):$$(id -g)")

export VERSION := $(shell cat VERSION)
export BUILD := $(shell git rev-parse HEAD | cut -c1-8)
export LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

export OWNER := appcelerator
export REPO := github.com/$(OWNER)/amp

# COMMON DIRECTORIES
# =============================================================================
CMDDIR := cmd

# =============================================================================
# DEFAULT TARGET
# =============================================================================
all: build

# =============================================================================
# VENDOR MANAGEMENT (GLIDE)
# =============================================================================
GLIDETARGETS := vendor

$(GLIDETARGETS): glide.yaml
	@glide install
# TODO: temporary fix for trace conflict, remove when resolved
	@rm -rf vendor/github.com/docker/docker/vendor/golang.org/x/net/trace

install-deps: $(GLIDETARGETS)

.PHONY: update-deps
update-deps:
	@glide update
# TODO: temporary fix for trace conflict, remove when resolved
	@rm -rf vendor/github.com/docker/docker/vendor/golang.org/x/net/trace

.PHONY: clean-glide
clean-glide:
	@rm -rf vendor

.PHONY: cleanall-glide
cleanall-glide: clean-glide
	@rm -rf .glide

# =============================================================================
# PROTOC (PROTOCOL BUFFER COMPILER)
# Generate *.pb.go, *.pb.gw.go files in any non-excluded directory
# with *.proto files.
# =============================================================================
PROTODIRS := api cmd data tests
PROTOFILES := $(shell find $(PROTODIRS) -type f -name '*.proto')
PROTOGWFILES := $(shell find $(PROTODIRS) -type f -name '*.proto' -exec grep -l 'google.api.http' {} \;)
# Generate swagger.json files for protobuf types even if only exposed over gRPC, not REST API
PROTOTARGETS := $(PROTOFILES:.proto=.pb.go) $(PROTOGWFILES:.proto=.pb.gw.go) $(PROTOFILES:.proto=.swagger.json)

PROTOOPTS := \
	-I/go/src/ \
	-I/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:/go/src/ \
	--grpc-gateway_out=logtostderr=true:/go/src \
	--swagger_out=logtostderr=true:/go/src/

%.pb.go %.pb.gw.go %.swagger.json: %.proto
	@echo $<
	@protoc $(PROTOOPTS) /go/src/$(REPO)/$<

protoc: $(PROTOTARGETS)

.PHONY: clean-protoc
clean-protoc:
	@find . \( -name "*.pb.go" -o -name "*.pb.gw.go" -o -name "*.swagger.json" \) \
			$(EXCLUDE_DIRS_FILTER) -type f -delete

# =============================================================================
# CLEAN
# =============================================================================
.PHONY: clean cleanall
clean: clean-glide clean-protoc clean-cli
cleanall: clean cleanall-glide

# =============================================================================
# BUILD CLI (`amp`)
# Saves binary to `cmd/amp/amp.alpine`, then builds `appcelerator/amp` image
# =============================================================================
CLI := amp
CLIBINARY=$(CLI).alpine
CLIIMG := appcelerator/amp
CLITARGET := $(CMDDIR)/$(CLI)/$(CLIBINARY)
CLISRC := $(shell find ./cmd/amp -type f -name '*.go')

$(CLITARGET): $(GLIDETARGETS) $(PROTOTARGETS) $(CLISRC)
	@go build $(LDFLAGS) -o $(CLITARGET) $(REPO)/$(CMDDIR)/$(CLI)
	@docker build -t $(CLIIMG) $(CMDDIR)/$(CLI)

build-cli: $(CLITARGET)

.PHONY: clean-cli
clean-cli:
	@rm -f $(CLITARGET)

# =============================================================================
# BUILD
# =============================================================================
build: $(PROTOTARGETS)
	@echo To be implemented...

dump:
	@echo $(CLISRC)
