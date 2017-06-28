.PHONY: clean test build-image

SRC := $(shell find . -type f -name '*.go')
MAIN := cmd/main.go
TARGET := bin/envoy

build: $(TARGET)

$(TARGET): $(SRC)
	go build -o $(TARGET) $(MAIN)

clean:
	rm -f bin/envoy

test:
	go test -v -timeout 30m

build-image:
	@: "$${PACKAGE:?PACKAGE environment variable null, empty, or not set}"
	docker build -t appcelerator/amp-envoy --build-arg package=$${PACKAGE} .
