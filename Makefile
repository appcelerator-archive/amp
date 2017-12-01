SHELL := /bin/bash
BASEDIR := $(shell echo $${PWD})

# =============================================================================
# BUILD MANAGEMENT
# Variables declared here are used by this Makefile *and* are exported to
# override default values used by supporting scripts in the hack directory
# =============================================================================
export UG := $(shell echo "$$(id -u):$$(id -g)")

export VERSION ?= $(shell cat VERSION)
export BUILD := $(shell git rev-parse HEAD | cut -c1-8)
export LDFLAGS := "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -s -w"

export OWNER := appcelerator
export REPO := github.com/$(OWNER)/amp

export GOOS := $(shell go env | grep GOOS | sed 's/"//g' | cut -c6-)
export GOARCH := $(shell go env | grep GOARCH | sed 's/"//g' | cut -c8-)

# =============================================================================
# COMMON DIRECTORIES
# =============================================================================
COMMONDIRS := pkg docker
CMDDIR := cmd

# =============================================================================
# DEFAULT TARGET
# =============================================================================
all: build

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
	-I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway \
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
.PHONY: clean
clean: clean-protoc cleanall-cli clean-server clean-beat clean-agent clean-monit clean-plugin-local clean-plugin-aws clean-ampagent

# =============================================================================
# VENDOR
# =============================================================================
.PHONY: clean-vendor install-vendor

clean-vendor:
	@rm -rf vendor

# not named `vendor` to avoid inopportune triggering
install-vendor:
	@dep ensure -vendor-only
	@dep prune

# =============================================================================
# BUILD
# =============================================================================
# When running in the amptools container, set DOCKER_CMD="sudo docker"
DOCKER_CMD ?= "docker"

build-base: protoc swagger build-server build-gateway build-beat build-agent build-monit
build-plugins: build-plugin-local build-plugin-aws build-ampagent
build: build-base build-plugins buildall-cli

# =============================================================================
# BUILD CLI (`amp`)
# =============================================================================
AMP := amp
AMPTARGET := bin/$(GOOS)/$(GOARCH)/$(AMP)
PLUGINDIRS := cluster/plugin/aws/plugin cluster/plugin/local/plugin
AMPDIRS := $(CMDDIR)/$(AMP) cli $(COMMONDIRS) $(PLUGINDIRS)
AMPSRC := $(shell find $(AMPDIRS) -type f -name '*.go')
AMPPKG := $(REPO)/$(CMDDIR)/$(AMP)

$(AMPTARGET): $(PROTOTARGETS) $(AMPSRC) VERSION vendor
	@echo "Compiling $(AMP) source(s) ($(GOOS)/$(GOARCH))"
	@echo $?
	@GOOS=$(GOOS) GOARCH=$(GOARCH) hack/lbuild $(REPO)/bin $(AMP) $(AMPPKG) $(LDFLAGS)
	@echo "bin/$(GOOS)/$(GOARCH)/$(AMP)"

# Warning: this only builds the CLI for the current OS, so when building under `ampmake`,
# the binary will be created under `bin/linux/amd64`.
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
buildall-cli: $(AMPTARGET) VERSION
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

$(AMPLTARGET): $(PROTOTARGETS) $(AMPLSRC) VERSION vendor
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
	@docker image rm $(AMPLIMG) 2>/dev/null || true


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

$(GWTARGET): $(PROTOTARGETS) $(GWSRC) VERSION vendor
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
	@docker image rm $(GWIMG) 2>/dev/null || true

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

$(MONITTARGET): $(PROTOTARGETS) $(MONITSRC) VERSION vendor
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
	@docker image rm $(MONITIMG) 2>/dev/null || true

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

$(BEATTARGET): $(PROTOTARGETS) $(BEATSRC) VERSION vendor
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
	@docker image rm $(BEATIMG) 2>/dev/null || true

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

$(AGENTTARGET): $(PROTOTARGETS) $(AGENTSRC) VERSION vendor
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
	@docker image rm $(AGENTIMG) 2>/dev/null || true

# =============================================================================
# CLUSTER PLUGINS
# =============================================================================
CPDIR := cluster/plugin

# =============================================================================
# BUILD AWS CLUSTER PLUGIN (`amp-aws`)
# Saves binary to `cluster/plugin/aws/aws.alpine`,
# then builds `appcelerator/amp-aws` image
# =============================================================================
CPAWS := amp-aws
CPAWSBINARY=$(CPAWS).alpine
CPAWSTAG := $(VERSION)
CPAWSIMG := appcelerator/$(CPAWS):$(CPAWSTAG)
CPAWSDIR := $(CPDIR)/aws
CPAWSDIRS := $(CPAWSDIR) $(COMMONDIRS)
CPAWSTARGET := $(CPAWSDIR)/$(CPAWSBINARY)
CPAWSSRC := $(shell find $(CPAWSDIRS) -type f -name '*.go')
CPAWSPKG := $(REPO)/$(CPAWSDIR)

$(CPAWSTARGET): $(PROTOTARGETS) $(CPAWSSRC) VERSION vendor
	@echo "Compiling $(CPAWS) source(s):"
	@echo $?
	@hack/build4alpine $(REPO)/$(CPAWSTARGET) $(CPAWSPKG) $(LDFLAGS)
	@echo "bin/$(GOOS)/$(GOARCH)/$(CPAWS)"

build-plugin-aws: $(CPAWSTARGET)
	@echo "build $(CPAWSIMG)"
	@$(DOCKER_CMD) build -t $(CPAWSIMG) $(CPAWSDIR) || (rm -f $(CPAWSTARGET); exit 1)

rebuild-cpaws: clean-plugin-aws build-plugin-aws

.PHONY: clean-plugin-aws
clean-plugin-aws:
	@rm -f $(CPAWSTARGET)
	@docker image rm $(CPAWSIMG) 2>/dev/null || true

# =============================================================================
# BUILD LOCAL CLUSTER PLUGIN (`amp-local`)
# Saves binary to `cluster/plugin/local/local.alpine`,
# then builds `appcelerator/amp-local` image
# =============================================================================
CPLOCAL := amp-local
CPLOCALBINARY=$(CPLOCAL).alpine
CPLOCALTAG := $(VERSION)
CPLOCALIMG := appcelerator/$(CPLOCAL):$(CPLOCALTAG)
CPLOCALDIR := $(CPDIR)/local
CPLOCALDIRS := $(CPLOCALDIR) $(COMMONDIRS)
CPLOCALTARGET := $(CPLOCALDIR)/$(CPLOCALBINARY)
CPLOCALSRC := $(shell find $(CPLOCALDIRS) -type f -name '*.go')
CPLOCALPKG := $(REPO)/$(CPLOCALDIR)

$(CPLOCALTARGET): $(PROTOTARGETS) $(CPLOCALSRC) VERSION vendor
	@echo "Compiling $(CPLOCAL) source(s):"
	@echo $?
	@hack/build4alpine $(REPO)/$(CPLOCALTARGET) $(CPLOCALPKG) $(LDFLAGS)
	@echo "bin/$(GOOS)/$(GOARCH)/$(CPLOCAL)"

build-plugin-local: $(CPLOCALTARGET)
	@echo "build $(CPLOCALIMG)"
	@$(DOCKER_CMD) build -t $(CPLOCALIMG) $(CPLOCALDIR) || (rm -f $(CPLOCALTARGET); exit 1)

rebuild-cplocal: clean-plugin-local build-plugin-local

.PHONY: clean-plugin-local
clean-plugin-local:
	@rm -f $(CPLOCALTARGET)
	@docker image rm $(CPLOCALIMG) 2>/dev/null || true

# =============================================================================
# BUILD AMPAGENT (`ampagent`)
# Saves binary to `cluster/ampagent/ampagent.alpine`,
# then builds `appcelerator/ampagent` image
# =============================================================================
AMPAGENT := ampagent
AMPAGENTBINARY=$(AMPAGENT).alpine
AMPAGENTTAG := $(VERSION)
AMPAGENTIMG := appcelerator/$(AMPAGENT):$(AMPAGENTTAG)
AMPAGENTDIR := cluster/$(AMPAGENT)
AMPAGENTDIRS := $(AMPAGENTDIR) $(COMMONDIRS)
AMPAGENTTARGET := $(AMPAGENTDIR)/$(AMPAGENTBINARY)
AMPAGENTSRC := $(shell find $(AMPAGENTDIRS) -type f -name '*.go')
AMPAGENTPKG := $(REPO)/$(AMPAGENTDIR)

$(AMPAGENTTARGET): $(PROTOTARGETS) $(AMPAGENTSRC) VERSION vendor
	@echo "Compiling $(AMPAGENT) source(s):"
	@echo $?
	@hack/build4alpine $(REPO)/$(AMPAGENTTARGET) $(AMPAGENTPKG) $(LDFLAGS)
	@echo "bin/$(GOOS)/$(GOARCH)/$(AMPAGENT)"

build-ampagent: $(AMPAGENTTARGET)
	@echo "build $(AMPAGENTIMG)"
	@$(DOCKER_CMD) build -t $(AMPAGENTIMG) $(AMPAGENTDIR) || (rm -f $(AMPAGENTTARGET); exit 1)

rebuild-ampagent: clean-ampagent build-ampagent

.PHONY: clean-ampagent
clean-ampagent:
	@rm -f $(AMPAGENTTARGET)
	@docker image rm $(AMPAGENTIMG) 2>/dev/null || true

# =============================================================================
# Swagger combine
# =============================================================================
swagger: protoc
	@rm -f swagger.json
	@node swagger-combine.js | jq . > swagger.json

# =============================================================================
# Quality checks
# =============================================================================
CHECKDIRS := agent api cli cluster cmd data monitoring tests $(COMMONDIRS)
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

