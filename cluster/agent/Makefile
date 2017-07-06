.PHONY: vendor clean test build-image run

PKG := github.com/appcelerator/amp/cluster/agent/cmd
TARGET := bin/ampctl
IMAGE := appcelerator/ampagent

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

run:
	docker run -it --rm -v /var/run/docker:/var/run/docker $(IMAGE) monitor

