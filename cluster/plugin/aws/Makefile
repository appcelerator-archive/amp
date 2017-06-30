.PHONY: vendor clean test build-image

PKG := github.com/appcelerator/amp/cluster/plugin/aws
TARGET := bin/aws
IMAGE := appcelerator/amp-aws

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

