.PHONY: vendor clean test build-image run

PKG := github.com/appcelerator/amp/cluster/agent/cmd
TARGET := bin/ampctl
VERSION ?= $(shell cat ../../VERSION)
IMAGE := appcelerator/ampagent:$(VERSION)
BUILD := $(shell git rev-parse HEAD | cut -c1-8)
export LDFLAGS := "-s -X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

build:
	go build -ldflags $(LDFLAGS) -o $(TARGET) $(PKG)

vendor: vendor.conf
	vndr

clean:
	rm -f $(TARGET)

test:
	go test -v -timeout 30m $(PKG)/../admin

image:
	docker build --build-arg LDFLAGS=$(LDFLAGS) -t $(IMAGE) .

run:
	docker run -t --rm -v /var/run/docker:/var/run/docker:ro $(IMAGE) monitor

run-checks:
	docker run -t --rm -v /var/run/docker:/var/run/docker:ro -v /var/run/docker.sock:/var/run/docker.sock:ro appcelerator/ampagent check -a

