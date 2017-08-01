.PHONY: cleanproto cleancli rebuildcli cleanserver rebuildserver cleanall rebuildall deploy

export VERSION := $(shell cat VERSION)
export BUILD := $(shell git rev-parse HEAD | cut -c1-8)
export LDFLAGS := "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
export OWNER := appcelerator
export REPO := github.com/$(OWNER)/amp

VENDOR := vendor
COMMONDIRS := pkg

# ===============================================
# general
# ===============================================

all: protoc server cli

clean: cleanserver cleancli

cleanall: cleanproto cleanserver cleancli

rebuildall: rebuildserver rebuildcli

build: all

# ===============================================
# protocol buffers
# ===============================================
PROTODIRS := api cmd data tests $(COMMONDIRS)
PROTOFILES := $(shell find $(PROTODIRS) -type f -name '*.proto')
PROTOTARGETS := $(PROTOFILES:.proto=.pb.go)
PROTOGWFILES := $(shell find $(PROTODIRS) -type f -name '*.proto' -exec grep -l 'google.api.http' {} \;)
PROTOGWTARGETS := $(PROTOGWFILES:.proto=.pb.gw.go) $(PROTOGWFILES:.proto=.swagger.json)
PROTOALLTARGETS := $(PROTOTARGETS) $(PROTOGWTARGETS)

# build any proto target (.pb.go, .pb.gw.go, .swagger.json) that is missing or not newer than .proto
protoc: $(PROTOALLTARGETS)

# build any proto target - use dockerized protobuf toolchain
%.pb.go %.pb.gw.go %.swagger.json: %.proto
	@echo $<
	@echo "compile proto files"
	@ docker run -it --rm -v $${PWD}:/go/src/github.com/appcelerator/amp -w /go/src/github.com/appcelerator/amp appcelerator/amptools make protoc

cleanproto:
	rm -f $(PROTOALLTARGETS)

# ===============================================
# cli
# ===============================================
AMPPATH := cmd/amp
AMPDIRS := $(AMPPATH) cli $(COMMONDIRS)
AMPSRC := $(shell find $(AMPDIRS) -type f -name '*.go')
AMPTARGET := bin/darwin/amd64/amp
AMPPKG := $(REPO)/$(AMPPATH)

$(AMPTARGET): $(VENDOR) $(PROTOTARGETS) $(AMPSRC)
	@ echo "building amp"
	@go build -ldflags $(LDFLAGS) -o $(AMPTARGET) $(AMPPKG)

cli: $(AMPTARGET)

cleancli:
	@rm -f $(AMPTARGET)

rebuildcli: cleancli cli

# ===============================================
# server
# ===============================================
AMPLPATH := cmd/amplifier
AMPLDIRS := $(AMPLPATH) api data $(COMMONDIRS)
AMPLSRC := $(shell find $(AMPLDIRS) -type f -name '*.go')
AMPLTARGET := $(AMPLPATH)/amplifier.alpine
AMPLPKG := $(REPO)/$(AMPLPATH)
AMPLIMG := $(OWNER)/amplifier:$(VERSION)

$(AMPLTARGET): $(VENDOR) $(PROTOTARGETS) $(AMPLSRC)
	@ echo "building amplifier"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags $(LDFLAGS) -o $(AMPLTARGET) $(AMPLPKG)
	@docker build -t $(AMPLIMG) $(AMPLPATH)

server: $(AMPLTARGET)

cleanserver:
	@rm -f $(AMPLTARGET)

rebuildserver: cleanserver server

# ===============================================
# services
# ===============================================

deploy: $(AMPLTARGET)
	@cd cluster/agent && TAG=$(VERSION) docker stack deploy -c stacksamples/amplifier-lite.yml amplifier

