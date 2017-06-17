SRC := $(shell find . -type f -name '*.go')

build: aws.alpine
	docker build -t appcelerator/amp-aws .

aws.alpine: $(SRC)
	docker run -it --rm -v $${PWD}:/go/src/github.com/appcelerator/amp/cluster/plugin/aws \
		-w /go/src/github.com/appcelerator/amp/cluster/plugin/aws \
		appcelerator/amp-aws-compiler build -o aws.alpine cmd/main.go

.PHONY: compiler
compiler:
	docker build -f Dockerfile.compiler -t appcelerator/amp-aws-compiler .

.PHONY: clean
clean:
	rm -f aws.alpine

.PHONY: test
test: build
	docker run -it --rm -v $(HOME)/.aws:/root/.aws -v $(PWD):/go/src/github.com/appcelerator/amp/cluster/plugin/aws \
		-w /go/src/github.com/appcelerator/amp/cluster/plugin/aws \
		-e KEYNAME=$(KEYNAME) \
		-e REGION=$(REGION) \
		appcelerator/amp-aws-compiler test -v -timeout 30m

pwd:
	echo $(PWD)
