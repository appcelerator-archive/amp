.PHONY: clean test build-image

MAIN := cmd/main.go
TARGET := bin/aws

build:
	go build -o $(TARGET) $(MAIN)

clean:
	rm -f bin/aws

test:
	go test -v -timeout 30m

build-image:
	docker build -t appcelerator/amp-aws .

