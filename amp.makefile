.PHONY: cleancli rebuildcli cleanserver rebuildserver cleanall rebuildall run-amplifier

export VERSION := $(shell cat VERSION)
export BUILD := $(shell git rev-parse HEAD | cut -c1-8)
export LDFLAGS := "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
export OWNER := appcelerator
export REPO := github.com/$(OWNER)/amp

VENDOR := vendor
COMMONDIRS := pkg

PROTODIRS := api cmd data tests $(COMMONDIRS)
PROTOFILES := $(shell find $(PROTODIRS) -type f -name '*.proto')
PROTOTARGETS := $(PROTOFILES:.proto=.pb.go)
PROTOGWFILES := $(shell find $(PROTODIRS) -type f -name '*.proto' -exec grep -l 'google.api.http' {} \;)
PROTOGWTARGETS := $(PROTOGWFILES:.proto=.pb.gw.go) $(PROTOGWFILES:.pb.gw.go=.swagger.json)
PROTOALLTARGETS := $(PROTOTARGETS) $(PROTOGWTARGETS)

AMPDIRS := cmd/amp cli $(COMMONDIRS)
AMPSRC := $(shell find $(AMPDIRS) -type f -name '*.go')

AMPLDIRS := cmd/amplifier api data $(COMMONDIRS)
AMPLSRC := $(shell find $(AMPLDIRS) -type f -name '*.go')

AMPTARGET := bin/darwin/amd64/amp
AMPLTARGET := cmd/amplifier/amplifier.alpine

all: $(PROTOALLTARGETS) $(AMPLTARGET) $(AMPTARGET)

$(AMPTARGET): $(VENDOR) $(PROTOTARGETS) $(AMPSRC)
	@ echo "building amp"
	@go build -ldflags $(LDFLAGS) -o bin/darwin/amd64/amp github.com/appcelerator/amp/cmd/amp

$(AMPLTARGET): $(VENDOR) $(PROTOTARGETS) $(AMPLSRC)
	@ echo "building amplifier"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags $(LDFLAGS) -o cmd/amplifier/amplifier.alpine github.com/appcelerator/amp/cmd/amplifier
	@docker build -t appcelerator/amplifier:$(VERSION) cmd/amplifier || (rm -rf $(AMPLTARGET); exit 1)

%.pb.go %.pb.gw.go %.swagger.json: %.proto
	@echo $<
	@echo "compile proto files"
	@ docker run -it --rm -v $${PWD}:/go/src/github.com/appcelerator/amp -w /go/src/github.com/appcelerator/amp appcelerator/amptools make protoc

protoc: $(PROTOALLTARGETS)

cli: $(AMPTARGET)

cleancli:
	@rm -f $(AMPTARGET)

rebuildcli: cleancli $(AMPTARGET)

server: $(AMPLTARGET)

cleanserver:
	@rm -f $(AMPLTARGET)

rebuildserver: cleanserver $(AMPLTARGET)

cleanll: cleanserver cleancli

rebuildall: rebuildserver rebuildcli

build: all

runamplifier: $(AMPLTARGET)
	@cd cluster/agent && TAG=$(VERSION) docker stack deploy -c stacksamples/amplifier-lite.yml amplifier

