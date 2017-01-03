  
.PHONY: all clean proto-clean bin-clean install install-host fmt simplify check version build-image run proto proto-host rules
.PHONY: build build-cli-linux build-cli-darwin build-cli-windows build-server build-server-linux build-server-darwin build-server-windows
.PHONY: dist dist-linux dist-darwin dist-windows
.PHONY: test test-unit test-cli test-integration test-integration-host

SHELL := /bin/bash
BASEDIR := $(shell echo $${PWD})

VERSION_FILE=VERSION
# build variables (provided to binaries by linker LDFLAGS below)
VERSION := $(shell cat $(VERSION_FILE))
BUILD ?= $(shell git rev-parse HEAD | cut -c1-8)

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# ignore vendor directory for go files
SRC := $(shell find . -type f -name '*.go' -not -path './vendor/*' -not -path './.git/*')

# for walking directory tree (like for proto rule)
EXCLUDE_FILES_FILTER := -not -path './vendor/*' -not -path './.git/*' -not -path './.glide/*'
EXCLUDE_DIRS_FILTER := $(EXCLUDE_FILES_FILTER) -not -path '.' -not -path './vendor' -not -path './.git' -not -path './.glide'

# for tests
UNIT_TEST_PACKAGES        := $(shell find .                   -type f -name '*_test.go' -not -path './tests/*' $(EXCLUDE_DIRS_FILTER) -exec dirname {} \; | sort -u)
CLI_TEST_PACKAGES         := $(shell find ./tests/cli         -type f -name '*_test.go'                        $(EXCLUDE_DIRS_FILTER) -exec dirname {} \; | sort -u)
INTEGRATION_TEST_PACKAGES := $(shell find ./tests/integration -type f -name '*_test.go'                        $(EXCLUDE_DIRS_FILTER) -exec dirname {} \; | sort -u)

DIRS = $(shell find . -type d $(EXCLUDE_DIRS_FILTER))

# generated file dependencies for proto rule
PROTOFILES := $(shell find . \( -path ./vendor -o -path ./.git -o -path ./.glide -o -path ./tests \) -prune -o -type f -name '*.proto' -print)
PROTOGWFILES := $(shell find . \( -path ./vendor -o -path ./.git -o -path ./.glide -o -path ./tests \) -prune -o -type f -name '*.proto' -exec grep -l google/api/annotations {} \;)
PROTOTARGETS := $(PROTOFILES:.proto=.pb.go)
PROTOGWTARGETS := $(PROTOGWFILES:.proto=.pb.gw.go)
PROTOALLTARGETS := $(PROTOTARGETS) $(PROTOGWTARGETS)

# generated files that can be cleaned
GENERATED := $(shell find . -type f \( -name '*.pb.go' -o -name '*.pb.gw.go' \) $(EXCLUDE_FILES_FILTER))

# ignore generated files when formatting/linting/vetting
CHECKSRC := $(shell find . -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go' $(EXCLUDE_FILES_FILTER))

OWNER := appcelerator
REPO := github.com/$(OWNER)/amp

CMDDIR := cmd

# Binaries
CLI := amp
SERVER := amplifier
AGENT := amp-agent
LOGWORKER := amp-log-worker
GATEWAY := amplifier-gateway
FUNCTION_LISTENER := amp-function-listener
FUNCTION_WORKER := amp-function-worker
CLUSTERSERVER := adm-server
CLUSTERAGENT := adm-agent
AMPADM := ampadm

TAG ?= latest
IMAGE := $(OWNER)/amp:$(TAG)

# tools
# need UID:GID because files created by containerized tools when mounting
# cwd are set to root:root
UG := $(shell echo "$$(id -u $${USER}):$$(id -g $${USER})")

DOCKER_RUN := docker run -t --rm -u $(UG)

GOTOOLS := appcelerator/gotools2:1.2.0
GOOS := $(shell uname | tr [:upper:] [:lower:])
GOARCH := amd64
GO := $(DOCKER_RUN) --name go -e HOME=$$HOME -v $${HOME}/.ssh:$$HOME/.ssh:ro -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) -e GOOS=$(GOOS) -e GOARCH=$(GOARCH) $(GOTOOLS) go
GOTEST := $(DOCKER_RUN) --name go -e HOME=$$HOME -v $${HOME}/.ssh:$$HOME/.ssh:ro -v $${GOPATH}/bin:/go/bin -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) $(GOTOOLS) go test -v

GLIDE_DIRS := $${HOME}/.glide $${PWD}/.glide vendor
GLIDE := $(DOCKER_RUN) -e HOME=$$HOME -v $$HOME/.ssh:$$HOME/.ssh:ro  -v $$HOME/.gitconfig:$$HOME/.gitconfig -v $$HOME/.glide:$$HOME/.glide -e GLIDE_HOME=$$HOME/.glide -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) $(GOTOOLS) glide $${GLIDE_OPTS}
GLIDE_INSTALL := $(GLIDE) install
GLIDE_UPDATE := $(GLIDE) update

all: version check build

arch:
	@echo $(GOOS)

version:
	@echo "version: $(VERSION) (build: $(BUILD))"

install-deps:
	@$(GLIDE_INSTALL)
# temporary fix to trace conflict
	@rm -rf vendor/github.com/docker/docker/vendor/golang.org/x/net/trace

update-deps:
	@$(GLIDE_UPDATE)
# temporary fix to trace conflict
	@rm -rf vendor/github.com/docker/docker/vendor/golang.org/x/net/trace

# explicit rule to compile protobuf files
%.pb.go: %.proto
	@go run hack/proto.go
%.pb.gw.go: %.proto
	@go run hack/proto.go

# used to install when you're already inside a container
install-host: proto-host
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLI)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(SERVER)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(AGENT)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(LOGWORKER)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(GATEWAY)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(FUNCTION_LISTENER)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(FUNCTION_WORKER)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLUSTERSERVER)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLUSTERAGENT)

proto: $(PROTOALLTARGETS)

# used to run protoc when you're already inside a container
proto-host: $(PROTOFILES)
	@go run hack/proto.go -protoc

proto-clean:
	@rm -rf $(GENERATED)

bin-clean:
	@rm -f $$(which $(CLI)) ./$(CLI)
	@rm -f $$(which $(SERVER)) ./$(SERVER)
	@rm -f coverage.out coverage-all.out
	@rm -f $$(which $(AGENT)) ./$(AGENT)
	@rm -f $$(which $(LOGWORKER)) ./$(LOGWORKER)
	@rm -f $$(which $(GATEWAY)) ./$(GATEWAY)
	@rm -f $$(which $(FUNCTION_LISTENER)) ./$(FUNCTION_LISTENER)
	@rm -f $$(which $(FUNCTION_WORKER)) ./$(FUNCTION_WORKER)
	@rm -f $$(which $(CLUSTERSERVER)) ./$(CLUSTERSERVER)
	@rm -f $$(which $(CLUSTERAGENT)) ./$(CLUSTERAGENT)
	@rm -f *.exe

clean: proto-clean bin-clean

install: install-cli install-server install-agent install-log-worker install-gateway install-fn-listener install-fn-worker install-adm-server install-adm-agent install-ampadm

DATASRC := $(shell find ./data -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
APISRC := $(shell find ./api -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
VENDORSRC := $(shell find ./vendor -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
# binaries are compiled as soon as their sources are newer than the existing binary
CLISRC := $(shell find ./cmd/amp -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
SERVERSRC := $(shell find ./cmd/amplifier -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
AGENTSRC := $(shell find ./cmd/amp-agent -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
LOGWORKERSRC := $(shell find ./cmd/amp-log-worker -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
GATEWAYSRC := $(shell find ./cmd/amplifier-gateway -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
FUNCTIONLISTENERSRC := $(shell find ./cmd/amp-function-listener -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
FUNCTIONWORKERSRC := $(shell find ./cmd/amp-function-worker -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
install-cli: $(CLISRC) $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOALLTARGETS)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLI)
install-server: $(SERVERSRC) $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOALLTARGETS)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(SERVER)
install-agent: $(AGENTSRC) $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOALLTARGETS)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(AGENT)
install-log-worker: $(LOGWORKERSRC) $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOALLTARGETS)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(LOGWORKER)
install-gateway: $(GATEWAYSRC) $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOALLTARGETS)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(GATEWAY)
install-fn-listener: $(FUNCTIONLISTENERSRC) $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOALLTARGETS)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(FUNCTION_LISTENER)
install-fn-worker: $(FUNCTIONWORKERSRC) $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOALLTARGETS)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(FUNCTION_WORKER)
install-adm-server: $(GATEWAYSRC) $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOALLTARGETS)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLUSTERSERVER)
install-adm-agent: $(GATEWAYSRC) $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOALLTARGETS)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLUSTERAGENT)
install-ampadm: $(GATEWAYSRC) $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOALLTARGETS)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(AMPADM)

build: build-cli build-server build-agent build-log-worker build-gateway build-fn-listener build-fn-worker build-adm-server build-adm-agent build-ampadm

build-cli: proto $(CLISRC) $(APISRC) $(VENDORSRC) Makefile
	@hack/build $(CLI)
build-server: proto
	@hack/build $(SERVER)
build-agent: proto
	@hack/build $(AGENT)
build-log-worker: proto
	@hack/build $(LOGWORKER)
build-gateway: proto
	@hack/build $(GATEWAY)
build-adm-server: proto
	@hack/build $(CLUSTERSERVER)
build-adm-agent: proto
	@hack/build $(CLUSTERAGENT)
build-ampadm: proto
	@hack/build $(CLUSTERADM)
build-fn-listener: proto
	@hack/build $(FUNCTION_LISTENER)
build-fn-worker: proto
	@hack/build $(FUNCTION_WORKER)

build-server-image:
	@docker build --build-arg BUILD=$(BUILD) -t appcelerator/$(SERVER):$(TAG) .

build-cli-linux:
	@rm -f $(CLI)
	@env GOOS=linux GOARCH=amd64 VERSION=$(VERSION) BUILD=$(BUILD) hack/build $(CLI)

build-cli-darwin:
	@rm -f $(CLI)
	@env GOOS=darwin GOARCH=amd64 VERSION=$(VERSION) BUILD=$(BUILD) hack/build $(CLI)

build-cli-windows:
	@rm -f $(CLI).exe
	@env GOOS=windows GOARCH=amd64 VERSION=$(VERSION) BUILD=$(BUILD) hack/build $(CLI)

build-server-linux:
	@rm -f $(SERVER)
	@env GOOS=linux GOARCH=amd64 VERSION=$(VERSION) BUILD=$(BUILD) hack/build $(SERVER)

build-server-darwin:
	@rm -f $(SERVER)
	@env GOOS=darwin GOARCH=amd64 VERSION=$(VERSION) BUILD=$(BUILD) hack/build $(SERVER)

build-server-windows:
	@rm -f $(SERVER).exe
	@env GOOS=windows GOARCH=amd64 VERSION=$(VERSION) BUILD=$(BUILD) hack/build $(SERVER)

build-ampadm-linux:
	@rm -f $(AMPADM)
	@env GOOS=linux GOARCH=amd64 VERSION=$(VERSION) BUILD=$(BUILD) hack/build $(AMPADM)

build-ampadm-darwin:
	@rm -f $(AMPADM)
	@env GOOS=darwin GOARCH=amd64 VERSION=$(VERSION) BUILD=$(BUILD) hack/build $(AMPADM)

build-ampadm-windows:
	@rm -f $(AMPADM).exe
	@env GOOS=windows GOARCH=amd64 VERSION=$(VERSION) BUILD=$(BUILD) hack/build $(AMPADM)

dist-linux: build-cli-linux build-server-linux build-ampadm-linux
	@rm -f dist/Linux/x86_64/amp-$(VERSION).tgz
	@mkdir -p dist/Linux/x86_64
	@tar czf dist/Linux/x86_64/amp-$(VERSION).tgz $(CLI) $(SERVER) $(AMPADM)

dist-darwin: build-cli-darwin build-server-darwin build-ampadm-darwin
	@rm -f dist/Darwin/x86_64/amp-$(VERSION).tgz
	@mkdir -p dist/Darwin/x86_64
	@tar czf dist/Darwin/x86_64/amp-$(VERSION).tgz $(CLI) $(SERVER) $(AMPADM)

dist-windows: build-cli-windows build-server-windows build-ampadm-windows
	@rm -f dist/Windows/x86_64/amp-$(VERSION).zip
	@mkdir -p dist/Windows/x86_64
	@zip -q dist/Windows/x86_64/amp-$(VERSION).zip $(CLI).exe $(SERVER).exe $(AMPADM).exe

dist: dist-linux dist-darwin dist-windows

# format and simplify if possible (https://golang.org/cmd/gofmt/#hdr-The_simplify_command)
fmt:
	@gofmt -s -l -w $(CHECKSRC)

check:
	@test -z $(shell gofmt -l ${CHECKSRC} | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@$(DOCKER_RUN) -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) $(GOTOOLS) bash -c 'for p in $$(go list ./... | grep -v /vendor/); do golint $${p} | sed "/pb\.\(gw\.\)*go/d"; done'
	@go tool vet ${CHECKSRC}

build-image:
	@BUILD=$(BUILD) $(PWD)/build-amp-image.sh $(TAG)

run: build-image
	@CID=$(shell docker run --net=host -d --name $(SERVER) $(IMAGE)) && echo $${CID}

test-cli:
	@for pkg in $(CLI_TEST_PACKAGES) ; do \
		go test $$pkg ; \
	done

test-unit:
	@for pkg in $(UNIT_TEST_PACKAGES) ; do \
		go test $$pkg ; \
	done

test-integration:
	@docker service rm amp-integration-test > /dev/null 2>&1 || true
	@docker build -t appcelerator/amp-demo-function examples/functions/demo-function
	@docker build --build-arg BUILD=$(BUILD) -t appcelerator/amp-integration-test .
	@docker service create --network amp-infra --name amp-integration-test --restart-condition none appcelerator/amp-integration-test make BUILD=$(BUILD) test-integration-host
	@containerid=""; \
	while [[ $${containerid} == "" ]] ; do \
		containerid=`docker ps -qf 'name=amp-integration'`; \
		sleep 1 ; \
	done; \
	docker logs -f $$containerid; \
	rc=`docker inspect --format='{{.State.ExitCode}}' $$containerid` ; \
	docker service rm amp-integration-test > /dev/null 2>&1 || true ; \
	exit $$rc

test-integration-host:
	@for pkg in $(INTEGRATION_TEST_PACKAGES) ; do \
		go test $$pkg || exit 1 ; \
	done

test: test-unit test-integration test-cli

cover:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(TEST_PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out

rules:
	@hack/print-make-rules

