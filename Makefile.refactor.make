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
# Everything that should be excluded when walking directory tree
EXCLUDE_FILES_FILTER := -not -path './vendor/*' -not -path './.test/*' -not -path './.git/*' -not -path './.glide/*'
EXCLUDE_DIRS_FILTER := $(EXCLUDE_FILES_FILTER) -not -path '.' -not -path './.test' -not -path './vendor' -not -path './.git' -not -path './.glide'

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
# VENDOR MANAGEMENT (GLIDE)
# =============================================================================
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


