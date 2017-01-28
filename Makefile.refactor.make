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
GOTOOLS := appcelerator/gotools:1.3.0

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
.PHONY: protoc protoc-host protoc-clean

PROTOFILES := $(shell find . -type f -name '*.proto' $(EXCLUDE_DIRS_FILTER))
PROTOGWFILES := $(shell find . -type f -name '*.proto' $(EXCLUDE_DIRS_FILTER) -exec grep -l 'google.api.http' {} \;)
# Generate swagger.json files for protobuf types even if only exposed over gRPC, not REST API
PROTOTARGETS := $(PROTOFILES:.proto=.pb.go) $(PROTOGWFILES:.proto=.pb.gw.go) $(PROTOFILES:.proto=.swagger.json)

%.pb.go %.pb.gw.go %.swagger.json: %.proto
	@echo $@
	@go run hack/proto.go "/go/src/$(REPO)/$<"

protoc: $(PROTOTARGETS)

# can be used when already in a container 
protoc-host: $(PROTOFILES)
	@go run hack/proto.go -protoc

protoc-clean:
	@find . \( -name "*.pb.go" -o -name "*.pb.gw.go" -o -name "*.swagger.json" \) \
			$(EXCLUDE_DIRS_FILTER) -type f -delete

# =============================================================================
# CLEAN
# =============================================================================
.PHONY: clean
clean: bin-clean protoc-clean

# =============================================================================
# CHECK
# formatting, linting and vetting checks on source files
# =============================================================================
.PHONY: check fmt
# ignore generated files when formatting/linting/vetting
CHECKSRC := $(shell find . -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go' $(EXCLUDE_DIRS_FILTER))

check:
	@test -z $(shell $(DOCKER_RUN_CMD) -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) $(GOTOOLS) gofmt -l ${CHECKSRC} | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@$(DOCKER_RUN_CMD) -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) $(GOTOOLS) bash -c 'for p in $$(go list ./... | grep -v -e /vendor/ -e /images/ -e /.git/); do golint $${p} | sed "/pb\.\(gw\.\)*go/d"; done'
	@$(DOCKER_RUN_CMD) -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) $(GOTOOLS) go tool vet ${CHECKSRC}

# format and simplify if possible (https://golang.org/cmd/gofmt/#hdr-The_simplify_command)
fmt:
	@$(DOCKER_RUN_CMD) -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) $(GOTOOLS) gofmt -s -l -w $(CHECKSRC)

# =============================================================================
# BUILD
# =============================================================================

.PHONY: bin-clean
.PHONY: build build-cli build-server build-agent build-log-worker build-gateway build-adm-server build-adm-agent build-ampadm build-fn-listener build-fn-worker
.PHONY: install install-cli install-server install-agent install-log-worker install-gateway install-adm-server install-adm-agent install-ampadm install-fn-listener install-fn-worker
.PHONY: install-host

CMDDIR := cmd
# Binaries
CLI := amp
SERVER := amplifier
AGENT := amp-agent
LOGWORKER := amp-log-worker
GATEWAY := amplifier-gateway
FUNCTION_LISTENER := amp-function-listener
FUNCTION_WORKER := amp-function-worker
ADMSERVER := adm-server
ADMAGENT := adm-agent
AMPADM := ampadm

# binaries are compiled as soon as their sources are newer than the existing binary
# generic source files
DATASRC := $(shell find ./data -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
APISRC := $(shell find ./api -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
VENDORSRC := $(shell find ./vendor -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
# binary specific source files
CLISRC := $(shell find ./cmd/amp -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
SERVERSRC := $(shell find ./cmd/amplifier -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
AGENTSRC := $(shell find ./cmd/amp-agent -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
LOGWORKERSRC := $(shell find ./cmd/amp-log-worker -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
GATEWAYSRC := $(shell find ./cmd/amplifier-gateway -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
ADMSERVERSRC := $(shell find ./cmd/adm-server -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
ADMAGENTSRC := $(shell find ./cmd/adm-agent -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
AMPADMSRC := $(shell find ./cmd/ampadm -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
FUNCTIONLISTENERSRC := $(shell find ./cmd/amp-function-listener -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')
FUNCTIONWORKERSRC := $(shell find ./cmd/amp-function-worker -type f -name '*.go' -not -name '*.pb.go' -not -name '*.pb.gw.go')

# build binaries on a Docker host
# leaves the binaries in the repo (not in "$GOPATH/bin")
build-cli: protoc $(CLISRC)                       $(APISRC) $(VENDORSRC) $(DATASRC) $(PROTOTARGETS)
	@env VERSION=$(VERSION) hack/build $(CLI)
build-server: protoc $(SERVERSRC)                 $(APISRC) $(VENDORSRC) $(DATASRC) $(PROTOTARGETS)
	@env VERSION=$(VERSION) hack/build $(SERVER)
build-agent: protoc $(AGENTSRC)                   $(APISRC) $(VENDORSRC) $(DATASRC) $(PROTOTARGETS)
	@env VERSION=$(VERSION) hack/build $(AGENT)
build-log-worker: protoc $(LOGWORKERSRC)          $(APISRC) $(VENDORSRC) $(DATASRC) $(PROTOTARGETS)
	@env VERSION=$(VERSION) hack/build $(LOGWORKER)
build-gateway: protoc $(GATEWAYSRC)               $(APISRC) $(VENDORSRC) $(DATASRC) $(PROTOTARGETS)
	@env VERSION=$(VERSION) hack/build $(GATEWAY)
build-adm-server: protoc $(ADMSERVERSRC)          $(APISRC) $(VENDORSRC) $(DATASRC) $(PROTOTARGETS)
	@env VERSION=$(VERSION) hack/build $(ADMSERVER)
build-adm-agent: protoc $(ADMAGENTSRC)            $(APISRC) $(VENDORSRC) $(DATASRC) $(PROTOTARGETS)
	@env VERSION=$(VERSION) hack/build $(ADMAGENT)
build-ampadm: protoc $(AMPADMSRC)                 $(APISRC) $(VENDORSRC) $(DATASRC) $(PROTOTARGETS)
	@env VERSION=$(VERSION) hack/build $(AMPADM)
build-fn-listener: protoc $(FUNCTIONLISTENERSRC)  $(APISRC) $(VENDORSRC) $(DATASRC) $(PROTOTARGETS)
	@env VERSION=$(VERSION) hack/build $(FUNCTION_LISTENER)
build-fn-worker: protoc $(FUNCTIONWORKERSRC)      $(APISRC) $(VENDORSRC) $(DATASRC) $(PROTOTARGETS)
	@env VERSION=$(VERSION) hack/build $(FUNCTION_WORKER)

build: build-cli build-server build-agent build-log-worker build-gateway build-adm-server build-adm-agent build-ampadm build-fn-listener build-fn-worker

install-cli: $(CLISRC)                            $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOTARGETS)
	@go install $(LDFLAGS)               $(REPO)/$(CMDDIR)/$(CLI)
install-server: $(SERVERSRC)                      $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOTARGETS)
	@go install $(LDFLAGS)            $(REPO)/$(CMDDIR)/$(SERVER)
install-agent: $(AGENTSRC) $(DATASRC)                        $(APISRC) $(VENDORSRC) $(PROTOTARGETS)
	@go install $(LDFLAGS)             $(REPO)/$(CMDDIR)/$(AGENT)
install-log-worker: $(LOGWORKERSRC)               $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOTARGETS)
	@go install $(LDFLAGS)         $(REPO)/$(CMDDIR)/$(LOGWORKER)
install-gateway: $(GATEWAYSRC)                    $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOTARGETS)
	@go install $(LDFLAGS)           $(REPO)/$(CMDDIR)/$(GATEWAY)
install-fn-listener: $(FUNCTIONLISTENERSRC)       $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOTARGETS)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(FUNCTION_LISTENER)
install-fn-worker: $(FUNCTIONWORKERSRC)           $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOTARGETS)
	@go install $(LDFLAGS)   $(REPO)/$(CMDDIR)/$(FUNCTION_WORKER)
install-adm-server: $(ADMSERVERSRC)               $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOTARGETS)
	@go install $(LDFLAGS)         $(REPO)/$(CMDDIR)/$(ADMSERVER)
install-adm-agent: $(ADMAGENTSRC)                 $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOTARGETS)
	@go install $(LDFLAGS)          $(REPO)/$(CMDDIR)/$(ADMAGENT)
install-ampadm: $(AMPADMSRC)                      $(DATASRC) $(APISRC) $(VENDORSRC) $(PROTOTARGETS)
	@go install $(LDFLAGS)            $(REPO)/$(CMDDIR)/$(AMPADM)

install: install-cli install-server install-agent install-log-worker install-gateway install-adm-server install-adm-agent install-ampadm install-fn-listener install-fn-worker

# can be used when already in a container to install all binaries
# (used by the Dockerfile in same directory)
install-host: protoc-host
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(CLI)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(SERVER)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(AGENT)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(LOGWORKER)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(GATEWAY)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(ADMSERVER)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(ADMAGENT)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(AMPADM)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(FUNCTION_LISTENER)
	@go install $(LDFLAGS) $(REPO)/$(CMDDIR)/$(FUNCTION_WORKER)

bin-clean:
	@rm -f $$(which $(CLI)) ./$(CLI)
	@rm -f $$(which $(SERVER)) ./$(SERVER)
	@rm -f $$(which $(AGENT)) ./$(AGENT)
	@rm -f $$(which $(LOGWORKER)) ./$(LOGWORKER)
	@rm -f $$(which $(GATEWAY)) ./$(GATEWAY)
	@rm -f $$(which $(FUNCTION_LISTENER)) ./$(FUNCTION_LISTENER)
	@rm -f $$(which $(FUNCTION_WORKER)) ./$(FUNCTION_WORKER)
	@rm -f $$(which $(ADMSERVER)) ./$(ADMSERVER)
	@rm -f $$(which $(ADMAGENT)) ./$(ADMAGENT)
	@rm -f $$(which $(AMPADM)) ./$(AMPADM)
	@rm -f *.exe
	@rm -f coverage.out coverage-all.out

# =============================================================================
# TESTS
# =============================================================================
.PHONY: test-cli test-unit test-integration test-integration-host test cover

test-cli:
	@$(DOCKER_RUN_CMD) -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) $(GOTOOLS) bash -c \
	'for pkg in $(CLI_TEST_PACKAGES) ; do \
		go test $$pkg ; \
	done'

test-unit:
	@$(DOCKER_RUN_CMD) -v $${PWD}:/go/src/$(REPO) -w /go/src/$(REPO) $(GOTOOLS) bash -c \
	'for pkg in $(UNIT_TEST_PACKAGES) ; do \
		go test $$pkg ; \
	done'

test-integration-host:
	@for pkg in $(INTEGRATION_TEST_PACKAGES) ; do \
		go test $$pkg || exit 1 ; \
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

test: test-unit test-integration test-cli

cover:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(TEST_PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out

# =============================================================================
# MISC
# =============================================================================
# TODO: used for debugging makefile, will ultimately this remove when all finished
.PHONY: dump
dump:
	@echo $(SRCDIRS)


