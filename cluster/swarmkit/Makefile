.PHONY: vendor clean test build-image run

PKG := github.com/appcelerator/amp/cluster/swarmkit/cmd
TARGET := bin/swarmmon
IMAGE := appcelerator/swarmkit

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
	docker run -it --rm -v /var/run/docker:/var/run/docker appcelerator/swarmkit swarmmon

