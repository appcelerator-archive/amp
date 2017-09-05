.PHONY: vendor clean test build-image

PKG := github.com/appcelerator/amp/cluster/plugin/aws/cmd
TARGET := bin/aws
VERSION ?= $(shell cat ../../../VERSION)
BUILD := $(shell git rev-parse HEAD | cut -c1-8)
IMAGE := appcelerator/amp-aws
export LDFLAGS := "-s -X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

build:
	go build -ldflags $(LDFLAGS) -o $(TARGET) $(PKG)

vendor: vendor.conf
	vndr

clean:
	rm -f $(TARGET)
	@docker image rm $(IMAGE) 2>/dev/null || true

test:
	go test -v -timeout 30m

image:
	docker build --build-arg LDFLAGS=$(LDFLAGS) -t $(IMAGE) .

