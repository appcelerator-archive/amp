SRC := $(shell find . -type f -name '*.go')

build: aws.alpine
	docker build -t appcelerator/amp-aws-plugin .

aws.alpine: $(SRC)
	docker run -it --rm -v $${PWD}:/go/src/github.com/appcelerator/amp/cluster/plugin/aws \
		-w /go/src/github.com/appcelerator/amp/cluster/plugin/aws \
		go:aws build -o aws.alpine cmd/main.go

.PHONY: compiler
compiler:
	docker build -f Dockerfile.compiler -t go:aws .

.PHONY: clean
clean:
	rm -f aws.alpine

.PHONY: test
test:
	@echo $(KEYPAIR)
	@echo $(REGION)
	docker run -it --rm -v $(HOME)/.aws:/root/.aws -v $(PWD):/go/src/github.com/appcelerator/amp/cluster/plugin/aws \
		-w /go/src/github.com/appcelerator/amp/cluster/plugin/aws \
		-e KEYPAIR=$(KEYPAIR) \
		-e REGION=$(REGION) \
		go:aws test

