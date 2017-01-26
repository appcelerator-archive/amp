SHELL := /bin/bash
BASEDIR := $(shell echo $${PWD})

# =============================================================================
# BUILD MANAGEMENT
# Variables declared here are used by this Makefile *and* are exported to
# override default values used by supporting scripts in the hack directory
# =============================================================================
export UG := $(shell echo "$$(id -u $${USER}):$$(id -g $${USER})")

export VERSION := $(shell cat VERSION)
export BUILD := $(shell git rev-parse HEAD | cut -c1-8)
export LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

export OWNER := appcelerator
export REPO := github.com/$(OWNER)/amp

# =============================================================================
# COMMON FILE AND DIRECTORY FILTERS AND GLOB VARS
# =============================================================================
EXCLUDE_DIRS_FILTER := -not -path './.*' -not -path './.*/*' \
	-not -path './dist' -not -path './dist/*' \
	-not -path './docs' -not -path './docs/*' \
	-not -path './hack' -not -path './hack/*' \
	-not -path './images' -not -path './images/*' \
	-not -path './project' -not -path './project/*' \
	-not -path './vendor' -not -path './vendor/*'

INCLUDE_DIRS_FILTER := -path './*' -path './*/*' $(EXCLUDE_DIRS_FILTER)

SRCDIRS := $(shell find . -type d $(EXCLUDE_DIRS_FILTER))
GOSRC := $(shell find . -type f -name '*.go' $(EXCLUDE_DIRS_FILTER))

# COMMON DIRECTORIES
# =============================================================================
CMDDIR := cmd

# =============================================================================
# COMMON CONTAINER TOOLS
# =============================================================================
# Used by: glide, protoc, go
BUILDTOOL := appcelerator/amptools:latest

# =============================================================================
# DEFAULT TARGET
# =============================================================================
all: build

# =============================================================================
# VENDOR MANAGEMENT (GLIDE)
# =============================================================================
GLIDETARGETS := glide.lock vendor

$(GLIDETARGETS): glide.yaml
	@hack/amptools glide install
# TODO: temporary fix for trace conflict, remove when resolved
	@rm -rf vendor/github.com/docker/docker/vendor/golang.org/x/net/trace

install-deps: $(GLIDETARGETS)

.PHONY: update-deps
update-deps:
	@hack/amptools glide update
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
PROTOFILES := $(shell find . -type f -name '*.proto' $(EXCLUDE_DIRS_FILTER))
PROTOGWFILES := $(shell find . -type f -name '*.proto' $(EXCLUDE_DIRS_FILTER) -exec grep -l 'google.api.http' {} \;)
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
	@hack/amptools protoc $(PROTOOPTS) /go/src/$(REPO)/$<

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
CLIIMG := appcelerator/aamp
CLITARGET := $(CMDDIR)/$(CLI)/$(CLIBINARY)
CLISRC := $(shell find ./cmd/amp -type f -name '*.go' $(EXCLUDE_DIRS_FILTER))

$(CLITARGET): $(GLIDETARGETS) $(PROTOTARGETS) $(CLISRC)
	@hack/amptools go build $(LDFLAGS) -o $(CLITARGET) $(REPO)/$(CMDDIR)/$(CLI)
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

