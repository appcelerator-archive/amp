
build: aws.alpine
	docker build -t appcelerator/amp-aws-plugin .

aws.alpine: main.go
	docker run -it --rm -v $${PWD}:/go/src/github.com/appcelerator/amp/cluster/plugin/aws \
		-w /go/src/github.com/appcelerator/amp/cluster/plugin/aws \
		awssdk go build -o aws.alpine .

.PHONY: compiler
compiler:
	docker build -f Dockerfile.compiler -t awssdk .

.PHONY: clean
clean:
	rm -f aws.alpine

