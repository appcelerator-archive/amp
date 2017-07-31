.PHONY: vendor clean test build-image

PKG := github.com/appcelerator/amp/cluster/plugin/local/cmd
TARGET := bin/amp-local
VERSION ?= $(shell cat ../../../VERSION)
BUILD := $(shell git rev-parse HEAD | cut -c1-8)
IMAGE := appcelerator/amp-local
export LDFLAGS := "-s -X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

build:
	go build -o $(TARGET) $(PKG)

vendor: vendor.conf
	vndr

clean:
	rm -f $(TARGET)

test:
	go test -v -timeout 30m

image:
	docker build -t $(IMAGE) .

