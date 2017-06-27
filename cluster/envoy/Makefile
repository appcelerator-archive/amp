.PHONY: clean test build-image

SRC := $(shell find . -name vendor -prune -o -type f -name '*.go')
MAIN := cmd/main.go
TARGET := bin/envoy

build: $(TARGET)

$(TARGET): $(SRC) vendor
	go build -o $(TARGET) $(MAIN)

vendor: vendor.conf
	vndr

clean:
	rm -f bin/envoy

test:
	go test -v -timeout 30m

build-image:
	docker build -t appcelerator/amp-envoy .
