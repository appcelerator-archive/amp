.PHONY: vendor clean test build-image run

PKG := github.com/appcelerator/amp/cluster/agent/cmd
TARGET := bin/ampctl
VERSION ?= latest
IMAGE := appcelerator/ampagent:$(VERSION)

build:
	go build -o $(TARGET) $(PKG)

vendor: vendor.conf
	vndr

clean:
	rm -f $(TARGET)

test:
	go test -v -timeout 30m $(PKG)/../admin

image:
	docker build -t $(IMAGE) .

run:
	docker run -t --rm -v /var/run/docker:/var/run/docker:ro $(IMAGE) monitor

run-checks:
	docker run -t --rm -v /var/run/docker:/var/run/docker:ro -v /var/run/docker.sock:/var/run/docker.sock:ro appcelerator/ampagent check -a

