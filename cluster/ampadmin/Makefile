.PHONY: vendor clean test build-image

PKG := github.com/appcelerator/amp/cluster/ampadmin/cmd
TARGET := bin/ampadmin
IMAGE := appcelerator/ampadmin

build:
	go build -o $(TARGET) $(PKG)

vendor: vendor.conf
	vndr

clean:
	rm -f bin/envoy

test:
	go test -v -timeout 30m

build-image:
	docker build -t $(IMAGE) .

