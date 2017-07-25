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

# =============================================================================
# COMMON DIRECTORIES
# =============================================================================
COMMONDIRS := pkg
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
	@rm -rf vendor/github.com/docker/swarmkit/vendor/golang.org/x/net/trace

install-deps: $(GLIDETARGETS)

.PHONY: update-deps
update-deps:
	@glide update
# TODO: temporary fix for trace conflict, remove when resolved
	@rm -rf vendor/github.com/docker/docker/vendor/golang.org/x/net/trace
	@rm -rf vendor/github.com/docker/swarmkit/vendor/golang.org/x/net/trace

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
PROTODIRS := api cmd data tests $(COMMONDIRS)

# standard protobuf files
PROTOFILES := $(shell find $(PROTODIRS) -type f -name '*.proto')
PROTOTARGETS := $(PROTOFILES:.proto=.pb.go)

# grpc rest gateway protobuf files
PROTOGWFILES := $(shell find $(PROTODIRS) -type f -name '*.proto' -exec grep -l 'google.api.http' {} \;)
PROTOGWTARGETS := $(PROTOGWFILES:.proto=.pb.gw.go) $(PROTOGWFILES:.pb.gw.go=.swagger.json)

PROTOOPTS := -I$(GOPATH)/src/ \
	-I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:$(GOPATH)/src/ \
	--grpc-gateway_out=logtostderr=true:$(GOPATH)/src \
	--swagger_out=logtostderr=true:$(GOPATH)/src/

PROTOALLTARGETS := $(PROTOTARGETS) $(PROTOGWTARGETS)

%.pb.go %.pb.gw.go %.swagger.json: %.proto
	@echo $<
	@protoc $(PROTOOPTS) $(GOPATH)/src/$(REPO)/$<

protoc: $(PROTOALLTARGETS)

.PHONY: clean-protoc
clean-protoc:
	@find $(PROTODIRS) \( -name '*.pb.go' -o -name '*.pb.gw.go' -o -name '*.swagger.json' \) -type f -delete

# =============================================================================
# CLEAN
# =============================================================================
.PHONY: clean cleanall
# clean doesn't remove the vendor directory since installing is time-intensive;
# you can do this explicitly: `ampmake clean-deps clean`

clean: clean-protoc cleanall-cli clean-server clean-beat clean-agent clean-monit
cleanall: clean cleanall-deps

# =============================================================================
# BUILD
# =============================================================================
# When running in the amptools container, set DOCKER_CMD="sudo docker"
DOCKER_CMD ?= "docker"

build-base: install-deps protoc build-server build-gateway build-beat build-agent build-bootstrap build-monit
build: build-base buildall-cli

# =============================================================================
# BUILD CLI (`amp`)
# Saves binary to `cmd/amp/amp.alpine`, then builds `appcelerator/amp` image
# =============================================================================
AMP := amp
AMPTARGET := bin/$(GOOS)/$(GOARCH)/$(AMP)
AMPDIRS := $(CMDDIR)/$(AMP) cli $(COMMONDIRS)
AMPSRC := $(shell find $(AMPDIRS) -type f -name '*.go')
AMPPKG := $(REPO)/$(CMDDIR)/$(AMP)

$(AMPTARGET): $(GLIDETARGETS) $(PROTOTARGETS) $(AMPSRC)
	@echo "Compiling $(AMP) source(s) ($(GOOS)/$(GOARCH))"
	@echo $?
	@GOOS=$(GOOS) GOARCH=$(GOARCH) hack/lbuild $(REPO)/bin $(AMP) $(AMPPKG) $(LDFLAGS)
	@echo "bin/$(GOOS)/$(GOARCH)/$(AMP)"

# Warning: this only builds the CLI for the current OS, so when building under `ampmake`,
# the binary will be created under `bin/linux/amd64`.
build-cli: $(AMPTARGET) build-bootstrap

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
AMPLTAG := $(VERSION)
AMPLIMG := appcelerator/$(AMPL):$(AMPLTAG)
AMPLTARGET := $(CMDDIR)/$(AMPL)/$(AMPLBINARY)
AMPLDIRS := $(CMDDIR)/$(AMPL) api data $(COMMONDIRS)
AMPLSRC := $(shell find $(AMPLDIRS) -type f -name '*.go')
AMPLPKG := $(REPO)/$(CMDDIR)/$(AMPL)

$(AMPLTARGET): $(GLIDETARGETS) $(PROTOTARGETS) $(AMPLSRC)
	@echo "Compiling $(AMPL) source(s):"
	@echo $?
	@hack/build4alpine $(REPO)/$(AMPLTARGET) $(AMPLPKG) $(LDFLAGS)
	@echo "bin/$(GOOS)/$(GOARCH)/$(AMPL)"

build-server: $(AMPLTARGET)
	@echo "build $(AMPLIMG)"
	@$(DOCKER_CMD) build -t $(AMPLIMG) $(CMDDIR)/$(AMPL) || (rm -f $(AMPLTARGET); exit 1)

rebuild-server: clean-server build-server

.PHONY: clean-server
clean-server:
	@rm -f $(AMPLTARGET)


# =============================================================================
# BUILD GATEWAY (`gateway`)
# Saves binary to `cmd/gateway/gateway.alpine`,
# then builds `appcelerator/gateway` image
# =============================================================================
GW := gateway
GWBINARY=$(GW).alpine
GWTAG := $(VERSION)
GWIMG := appcelerator/$(GW):$(GWTAG)
GWTARGET := $(CMDDIR)/$(GW)/$(GWBINARY)
GWDIRS := $(CMDDIR)/$(GW) api data $(COMMONDIRS)
GWSRC := $(shell find $(GWDIRS) -type f -name '*.go')
GWPKG := $(REPO)/$(CMDDIR)/$(GW)

$(GWTARGET): $(GLIDETARGETS) $(PROTOTARGETS) $(GWSRC)
	@echo "Compiling $(GW) source(s):"
	@echo $?
	@hack/build4alpine $(REPO)/$(GWTARGET) $(GWPKG) $(LDFLAGS)
	@echo "bin/$(GOOS)/$(GOARCH)/$(GW)"

build-gateway: $(GWTARGET)
	@echo "build $(GWIMG)"
	@$(DOCKER_CMD) build -t $(GWIMG) $(CMDDIR)/$(GW) || (rm -f $(GWTARGET); exit 1)

rebuild-gateway: clean-gateway build-gateway

.PHONY: clean-gateway
clean-gateway:
	@rm -f $(GWTARGET)

# =============================================================================
# BUILD monitoring (`promctl`)
# Saves binary to `cmd/monitoring/promctl.alpine`,
# then builds `appcelerator/amp-prometheus` image
# =============================================================================
MONIT := monitoring
MONITBINARY=promctl.alpine
MONITTAG := $(VERSION)
MONITIMG := appcelerator/amp-prometheus:$(MONITTAG)
MONITTARGET := $(MONIT)/$(MONITBINARY)
MONITDIRS := $(MONIT)/promctl $(COMMONDIRS)
MONITSRC := $(shell find $(MONITDIRS) -type f -name '*.go')
MONITPKG := $(REPO)/$(MONIT)/promctl

$(MONITTARGET): $(GLIDETARGETS) $(PROTOTARGETS) $(MONITSRC)
	@echo "Compiling $(MONIT) source(s):"
	@echo $?
	@hack/build4alpine $(REPO)/$(MONITTARGET) $(MONITPKG) $(LDFLAGS)
	@echo "bin/$(GOOS)/$(GOARCH)/$(MONIT)"

build-monit: $(MONITTARGET)
	@echo "build $(MONITIMG)"
	@$(DOCKER_CMD) build -t $(MONITIMG) $(MONIT) || (rm -f $(MONITTARGET); exit 1)

rebuild-monit: clean-monit build-monit

.PHONY: clean-monit
clean-monit:
	@rm -f $(MONITTARGET)

# =============================================================================
# BUILD BEAT (`ampbeat`)
# Saves binary to `cmd/ampbeat/ampbeat.alpine`,
# then builds `appcelerator/ampbeat` image
# =============================================================================
BEAT := ampbeat
BEATBINARY=$(BEAT).alpine
BEATTAG := $(VERSION)
BEATIMG := appcelerator/$(BEAT):$(BEATTAG)
BEATTARGET := $(CMDDIR)/$(BEAT)/$(BEATBINARY)
BEATDIRS := $(CMDDIR)/$(BEAT) api data $(COMMONDIRS)
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
AGENTTAG := $(VERSION)
AGENTIMG := appcelerator/$(AGENT):$(AGENTTAG)
AGENTTARGET := $(CMDDIR)/$(AGENT)/$(AGENTBINARY)
AGENTDIRS := $(CMDDIR)/$(AGENT) agent api $(COMMONDIRS)
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
# BUILD BOOTSTRAP (`amp-bootstrap`)
# Bootstrap local amp cluster
# =============================================================================
AMPBOOTDIR := platform
AMPBOOTBIN := platform
AMPBOOTIMG := appcelerator/amp-bootstrap
AMPBOOTVER ?= $(VERSION)
AMPBOOTSRC := platform/bin/deploy platform/bin/dev $(shell find $(AMPBOOTDIR) -type f)

.PHONY: build-bootstrap
build-bootstrap:
	@echo "Building $(AMPBOOTIMG):$(AMPBOOTVER)"
	@rm -f $(AMPBOOTDIR)/stacks/*.pem
	@$(DOCKER_CMD) build -t $(AMPBOOTIMG):$(AMPBOOTVER) $(AMPBOOTDIR) >/dev/null

.PHONY: push-bootstrap
push-bootstrap:
	@echo "Pushing $(AMPBOOTIMG)"
	@$(AMPBOOTDIR)/build

# =============================================================================
# BUILD CLUSTER PLUGINS (`amp-aws`)
# Needed for `amp cluster init` CLI command support.
# Plugins are located under `cluster/plugin`
# =============================================================================
CPDIR := cluster/plugin
CPAWSDIR := $(CPDIR)/aws

.PHONY: build-amp-aws-compiler
build-amp-aws-compiler:
	@cd $(CPAWSDIR) && $(MAKE) compiler

.PHONY: build-amp-aws
build-amp-aws: build-amp-aws-compiler
	@cd $(CPAWSDIR) && $(MAKE) build

# WARNING:
# If the environment variables for $KEYNAME and $REGION are not set, the
# test will fail with misleading error output:
# docker run -it --rm -v /root/.aws:/root/.aws -v /go/src/github.com/appcelerator/amp/cluster/plugin/aws:/go/src/github.com/appcelerator/amp/cluster/plugin/aws \
#          -w /go/src/github.com/appcelerator/amp/cluster/plugin/aws \
#          -e KEYNAME= \
#          -e REGION= \
#          appcelerator/amp-aws-compiler test -v -timeout 30m
#  can't load package: package github.com/appcelerator/amp/cluster/plugin/aws: no buildable Go source files in /go/src/github.com/appcelerator/amp/cluster/plugin/aws
#  Makefile:21: recipe for target 'test' failed
#
# To succeed, ensure that these environment variables have values, like this:
# $ KEYNAME=tony-amp-dev REGION=us-west-2 make test-amp-aws
#
.PHONY: test-amp-aws
test-amp-aws: build-amp-aws
	@cd $(CPAWSDIR) && $(MAKE) test

# =============================================================================
# BUILD AMPADMIN (`ampadmin`)
# =============================================================================
AMPADMINDIR := cluster/ampadmin

.PHONY: build-ampadmin
build-ampadmin: build-ampadmin
	@cd $(AMPADMINDIR) && $(MAKE) build

.PHONY: test-ampadmin
test-ampadmin:
	@cd $(AMPADMINDIR) && $(MAKE) test

# =============================================================================
# Quality checks
# =============================================================================
CHECKDIRS := agent api cli cmd data monitoring tests $(COMMONDIRS)
CHECKSRCS := $(shell find $(CHECKDIRS) -type f \( -name '*.go' -and -not -name '*.pb.go' -and -not -name '*.pb.gw.go'  \))

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
# Run check before submitting a pull request!
# =============================================================================

check: fmt buildall lint-fast

