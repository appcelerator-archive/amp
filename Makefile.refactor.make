SHELL := /bin/bash
BASEDIR := $(shell echo $${PWD})

# =============================================================================
# BUILD MANAGEMENT
# =============================================================================
# VERSION and BUILD are build variables supplied to binaries by go linker LDFLAGS option
VERSION_FILE := VERSION
VERSION := $(shell cat $(VERSION_FILE))
BUILD ?= $(shell git rev-parse HEAD | cut -c1-8)
LDFLAGS := -ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

OWNER := appcelerator
REPO := github.com/$(OWNER)/amp

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

# =============================================================================
# DOCKER SUPPORT
# =============================================================================
# Used so that files created in containers using mounted volumes are
# owned by current UID:GID instead of root:root
UG := $(shell echo "$$(id -u $${USER}):$$(id -g $${USER})")

# Base docker command
DOCKER_RUN_CMD := docker run -t --rm -u $(UG)

# =============================================================================
# COMMON CONTAINER TOOLS
# =============================================================================
# Used by: glide, protoc, go
GOTOOLS := appcelerator/gotools:latest

# =============================================================================
# DEFAULT TARGET
# =============================================================================
all: build

# =============================================================================
# VENDOR MANAGEMENT (GLIDE)
# =============================================================================
.PHONY: install-deps update-deps

# Mount ~/.ssh (for access to private git repos), glide cache, and working directory (for ~/vendor)
GLIDE_BASE_CMD := $(DOCKER_RUN_CMD) \
                  -e HOME=$${HOME} \
                  -v $${HOME}/.ssh:$${HOME}/.ssh:ro \
                  -v $${HOME}/.gitconfig:$${HOME}/.gitconfig:ro \
                  -e GLIDE_HOME=/tmp/glide \
                  -v $${PWD}:/go/src/$(REPO) \
                  -v glide:/tmp/glide \
                  -w /go/src/$(REPO) \
                  $(GOTOOLS) glide $${GLIDE_OPTS}
GLIDE_INSTALL_CMD := $(GLIDE_BASE_CMD) install
GLIDE_UPDATE_CMD := $(GLIDE_BASE_CMD) update

install-deps:
	@echo $(GLIDE_INSTALL_CMD)
	@$(GLIDE_INSTALL_CMD)
# TODO: temporary fix for trace conflict, remove when resolved
	@rm -rf vendor/github.com/docker/docker/vendor/golang.org/x/net/trace

update-deps:
	@$(GLIDE_UPDATE_CMD)
# TODO: temporary fix for trace conflict, remove when resolved
	@rm -rf vendor/github.com/docker/docker/vendor/golang.org/x/net/trace

# =============================================================================
# PROTOC (PROTOCOL BUFFER COMPILER)
# Generate *.pb.go, *.pb.gw.go files in any non-excluded directory
# with *.proto files.
# =============================================================================
.PHONY: protoc protoc-clean

PROTOFILES := $(shell find . -type f -name '*.proto' $(EXCLUDE_DIRS_FILTER))
PROTOGWFILES := $(shell find . -type f -name '*.proto' $(EXCLUDE_DIRS_FILTER) -exec grep -l 'google.api.http' {} \;)
# Generate swagger.json files for protobuf types even if only exposed over gRPC, not REST API
PROTOTARGETS := $(PROTOFILES:.proto=.pb.go) $(PROTOGWFILES:.proto=.pb.gw.go) $(PROTOFILES:.proto=.swagger.json)

%.pb.go %.pb.gw.go %.swagger.json: %.proto
	@echo $@
	@go run hack/proto.go "/go/src/$(REPO)/$<"

protoc: $(PROTOTARGETS)

protoc-clean:
	@find . \( -name "*.pb.go" -o -name "*.pb.gw.go" -o -name "*.swagger.json" \) \
			$(EXCLUDE_DIRS_FILTER) -type f -delete

# =============================================================================
# CLEAN
# =============================================================================
.PHONY: clean
clean: protoc-clean

# =============================================================================
# BUILD
# =============================================================================
build: $(PROTOTARGETS)
	@echo To be implemented...

# =============================================================================
# MISC
# =============================================================================
# TODO: used for debugging makefile, will ultimately this remove when all finished
.PHONY: dump
dump:
	@echo $(SRCDIRS)


