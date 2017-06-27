SRC := $(shell find . -type f -name '*.go')

build: envoy.alpine
	docker build -t appcelerator/amp-envoy .

envoy.alpine: $(SRC)
	docker run -it --rm -v $${PWD}:/go/src/github.com/appcelerator/amp/cluster/envoy \
		-w /go/src/github.com/appcelerator/amp/cluster/envoy \
		appcelerator/amp-envoy-compiler build -o envoy.alpine cmd/main.go

.PHONY: compiler
compiler:
	docker build -f Dockerfile.compiler -t appcelerator/amp-envoy-compiler .

.PHONY: clean
clean:
	rm -f envoy.alpine

.PHONY: test
test: build
# TODO wip
	docker run -it --rm -v $(PWD):/go/src/github.com/appcelerator/amp/cluster/envoy \
		-w /go/src/github.com/appcelerator/amp/cluster/envoy \
		appcelerator/amp-envoy-compiler test -v -timeout 30m

