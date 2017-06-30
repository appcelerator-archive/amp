.PHONY: vendor clean test build-image

MAIN := cmd/main.go
TARGET := bin/aws

build:
	go build -o $(TARGET) $(MAIN)

vendor: vendor.conf
	vndr

clean:
	rm -f bin/aws

test:
	go test -v -timeout 30m

build-image:
	docker build -t appcelerator/amp-aws .

