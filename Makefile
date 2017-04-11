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
export LDFLAGS := "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

export OWNER := appcelerator
export REPO := github.com/$(OWNER)/amp

export GOOS := $(shell go env | grep GOOS | sed 's/"//g' | cut -c6-)
export GOARCH := $(shell go env | grep GOARCH | sed 's/"//g' | cut -c8-)

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
	@glide install || (rm -rf vendor; exit 1)
# TODO: temporary fix for trace conflict, remove when resolved
	@rm -rf vendor/github.com/docker/docker/vendor/golang.org/x/net/trace

install-deps: $(GLIDETARGETS)

.PHONY: update-deps
update-deps:
	@glide update
# TODO: temporary fix for trace conflict, remove when resolved
	@rm -rf vendor/github.com/docker/docker/vendor/golang.org/x/net/trace

.PHONY: clean-deps
clean-deps:
	@rm -rf vendor

.PHONY: cleanall-deps
# cleanall-deps will effectively causes `install-deps` to behave like `update-deps`
cleanall-deps: clean-deps
	@rm -rf .glide glide.lock

# =============================================================================
# PROTOC (PROTOCOL BUFFER COMPILER)
# Generate *.pb.go, *.pb.gw.go files in any non-excluded directory
# with *.proto files.
# =============================================================================
PROTODIRS := api cmd data tests
PROTOFILES := $(shell find $(PROTODIRS) -type f -name '*.proto')
PROTOTARGETS := $(PROTOFILES:.proto=.pb.go)
PROTOOPTS := -I/go/src/ --go_out=plugins=grpc:/go/src/

%.pb.go: %.proto
	@echo $<
	@protoc $(PROTOOPTS) /go/src/$(REPO)/$<

protoc: $(PROTOTARGETS)

.PHONY: clean-protoc
clean-protoc:
	@find $(PROTODIRS) \( -name "*.pb.go" \) -type f -delete

# =============================================================================
# CLEAN
# =============================================================================
.PHONY: clean cleanall
# clean doesn't remove the vendor directory since installing is time-intensive;
# you can do this explicitly: `ampmake clean-deps clean`

clean: clean-protoc cleanall-cli clean-server clean-beat clean-agent
cleanall: clean cleanall-deps

# =============================================================================
# BUILD
# =============================================================================
# When running in the amptools container, set DOCKER_CMD="sudo docker"
DOCKER_CMD ?= "docker"

build: install-deps protoc build-server build-cli build-beat build-agent
buildall: install-deps protoc build-server buildall-cli build-beat build-agent

# =============================================================================
# BUILD CLI (`amp`)
# Saves binary to `cmd/amp/amp.alpine`, then builds `appcelerator/amp` image
# =============================================================================
AMP := amp
AMPTARGET := bin/$(GOOS)/$(GOARCH)/$(AMP)
AMPDIRS := $(CMDDIR)/$(AMP) cli tests
AMPSRC := $(shell find $(AMPDIRS) -type f -name '*.go')
AMPPKG := $(REPO)/$(CMDDIR)/$(AMP)

$(AMPTARGET): $(GLIDETARGETS) $(PROTOTARGETS) $(AMPSRC)
	@echo "Compiling $(AMP) source(s) ($(GOOS)/$(GOARCH))"
	@echo $?
	@GOOS=$(GOOS) GOARCH=$(GOARCH) hack/lbuild $(REPO)/bin $(AMP) $(AMPPKG) $(LDFLAGS)
	@echo "bin/$(GOOS)/$(GOARCH)/$(AMP)"

build-cli: $(AMPTARGET)

.PHONY: rebuild-cli
rebuild-cli: clean-cli build-cli

.PHONY: rebuildall-cli
rebuildall-cli: cleanall-cli buildall-cli

.PHONY: clean-cli
clean-cli:
	@rm -f $(AMPTARGET)

.PHONY: cleanall-cli
cleanall-cli:
# following fails in gogland shell
#	@(shopt -s extglob; rm -f bin/*(darwin|linux|alpine)/amd64/amp)
	@rm -f bin/darwin/amd64/amp bin/linux/amd64/amp bin/alpine/amd64/amp

# Build cross-compiled versions of the cli
buildall-cli: $(AMPTARGET)
	@echo "cross-compiling $(AMP) cli for supported targets"
	@hack/xbuild $(REPO)/bin $(AMP) $(REPO)/$(CMDDIR)/$(AMP) $(LDFLAGS)

# =============================================================================
# BUILD SERVER (`amplifier`)
# Saves binary to `cmd/amplifier/amplifier.alpine`,
# then builds `appcelerator/amplifier` image
# =============================================================================
AMPL := amplifier
AMPLBINARY=$(AMPL).alpine
AMPLTAG := local
AMPLIMG := appcelerator/$(AMPL):$(AMPLTAG)
AMPLTARGET := $(CMDDIR)/$(AMPL)/$(AMPLBINARY)
AMPLDIRS := cmd/$(AMPL) api data tests
AMPLSRC := $(shell find $(AMPLDIRS) -type f -name '*.go')
AMPLPKG := $(REPO)/$(CMDDIR)/$(AMPL)

$(AMPLTARGET): $(GLIDETARGETS) $(PROTOTARGETS) $(AMPLSRC)
	@echo "Compiling $(AMPL) source(s):"
	@echo $?
	@hack/build4alpine $(REPO)/$(AMPLTARGET) $(AMPLPKG) $(LDFLAGS)
	@echo "bin/$(GOOS)/$(GOARCH)/$(AMPL)"

build-server: $(AMPLTARGET)
	@echo "build $(AMPLIMG)"
	@cp -f ~/.config/amp/amplifier.y*ml cmd/amplifier &> /dev/null || (echo "Warning: ~/.config/amp/amplifier.yml not found (sendgrid key needed to send email)" && touch cmd/amplifier/amplifier.yml)
	@$(DOCKER_CMD) build -t $(AMPLIMG) $(CMDDIR)/$(AMPL) || (rm -f $(AMPLTARGET); exit 1)
	@rm -f cmd/amplifier/amplifier.yml

rebuild-server: clean-server build-server

.PHONY: clean-server
clean-server:
	@rm -f $(AMPLTARGET)


# =============================================================================
# BUILD BEAT (`ampbeat`)
# Saves binary to `cmd/ampbeat/ampbeat.alpine`,
# then builds `appcelerator/ampbeat` image
# =============================================================================
BEAT := ampbeat
BEATBINARY=$(BEAT).alpine
BEATTAG := local
BEATIMG := appcelerator/$(BEAT):$(BEATTAG)
BEATTARGET := $(CMDDIR)/$(BEAT)/$(BEATBINARY)
BEATDIRS := cmd/$(BEAT) api data tests
BEATSRC := $(shell find $(BEATDIRS) -type f -name '*.go')
BEATPKG := $(REPO)/$(CMDDIR)/$(BEAT)

$(BEATTARGET): $(GLIDETARGETS) $(PROTOTARGETS) $(BEATSRC)
	@echo "Compiling $(BEAT) source(s):"
	@echo $?
	@hack/build4alpine $(REPO)/$(BEATTARGET) $(BEATPKG) $(LDFLAGS)
	@echo "bin/$(GOOS)/$(GOARCH)/$(BEAT)"

build-beat: $(BEATTARGET)
	@echo "build $(BEATIMG)"
	@$(DOCKER_CMD) build -t $(BEATIMG) $(CMDDIR)/$(BEAT) || (rm -f $(BEATTARGET); exit 1)

rebuild-beat: clean-beat build-beat

.PHONY: clean-beat
clean-beat:
	@rm -f $(BEATTARGET)

# =============================================================================
# BUILD AGENT (`agent`)
# Saves binary to `cmd/agent/agent.alpine`,
# then builds `appcelerator/agent` image
# =============================================================================
AGENT := agent
AGENTBINARY=$(AGENT).alpine
AGENTTAG := local
AGENTIMG := appcelerator/$(AGENT):$(AGENTTAG)
AGENTTARGET := $(CMDDIR)/$(AGENT)/$(AGENTBINARY)
AGENTDIRS := cmd/$(AGENT) api data tests
AGENTSRC := $(shell find $(AGENTDIRS) -type f -name '*.go')
AGENTPKG := $(REPO)/$(CMDDIR)/$(AGENT)

$(AGENTTARGET): $(GLIDETARGETS) $(PROTOTARGETS) $(AGENTSRC)
	@echo "Compiling $(AGENT) source(s):"
	@echo $?
	@hack/build4alpine $(REPO)/$(AGENTTARGET) $(AGENTPKG) $(LDFLAGS)
	@echo "bin/$(GOOS)/$(GOARCH)/$(AGENT)"

build-agent: $(AGENTTARGET)
	@echo "build $(AGENTIMG)"
	@$(DOCKER_CMD) build -t $(AGENTIMG) $(CMDDIR)/$(AGENT) || (rm -f $(AGENTTARGET); exit 1)

rebuild-agent: clean-agent build-agent

.PHONY: clean-agent
clean-agent:
	@rm -f $(AGENTTARGET)

# =============================================================================
# Quality checks
# =============================================================================
CHECKDIRS := agent api cli cmd data pkg tests
CHECKSRCS := $(shell find $(CHECKDIRS) -type f \( -name '*.go' -and -not -name '*.pb.go' \))

# format and simplify if possible (https://golang.org/cmd/gofmt/#hdr-The_simplify_command)
.PHONY: fmt
fmt:
	@goimports -l $(CHECKDIRS) && goimports -w $(CHECKDIRS)
	@gofmt -s -l -w $(CHECKSRCS)

.PHONY: lint
lint:
	@echo "running lint checks - this will take a while..."
	@gometalinter --deadline=10m --concurrency=1 --enable-gc --vendor --exclude=vendor --exclude=\.pb\.go \
		--sort=path --aggregate \
		--disable-all \
		--enable=deadcode \
		--enable=errcheck \
		--enable=gas \
		--enable=goconst \
		--enable=gocyclo \
		--enable=gofmt \
		--enable=goimports \
		--enable=golint \
		--enable=gosimple \
		--enable=ineffassign \
		--enable=interfacer \
		--enable=staticcheck \
		--enable=structcheck \
		--enable=test \
		--enable=unconvert \
		--enable=unparam \
		--enable=unused \
		--enable=varcheck \
		--enable=vet \
		--enable=vetshadow \
		$(CHECKDIRS)

.PHONY: lint-fast
lint-fast:
	@echo "running subset of lint checks in fast mode"
	@gometalinter \
		--vendored-linters \
		--fast --deadline=10m --concurrency=1 --enable-gc --vendor --exclude=vendor --exclude=\.pb\.go \
		--disable=gotype \
		$(CHECKDIRS)

# =============================================================================
# Misc
# =============================================================================
# Display all the Makefile rules
.PHONY: rules
rules:
	@hack/print-make-rules

# Display pertinent environment variables
.PHONY: env
env:
	@echo "GOOS=$(GOOS)"
	@echo "GOARCH=$(GOARCH)"

# =============================================================================
# Local deployment for development
# =============================================================================
# use well-known guid for local cluster
CID := f573e897-7aa0-4516-a195-42ee91039e97

deploy: build
	@hack/deploy $(CID)

# =============================================================================
# Run check before submitting a pull request!
# =============================================================================

check: fmt buildall lint-fast

