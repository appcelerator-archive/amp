FROM golang:1.9
# This builds a convenience "all-in-one" image for go development.
# It intentionally does not remove any build prerequisites like most of our other
# images since this image is meant strictly for building things.

# The alpine-sdk adds a bunch of essential tools, such as
# ca-certificates, musl, musl-dev, openssl, curl, make, gcc, g++, libstdc++, libc-dev, tar, bzip2, git
# The rest are for building grpc or for any other convenience.
#RUN apk --no-cache add alpine-sdk bash vim autoconf automake libtool bind-tools drill
RUN apt-get update -y && apt-get install --no-install-recommends -y apt-utils autoconf automake bash ca-certificates dnsutils git jq libtool vim tar zip unzip && apt-get clean && rm -rf /var/lib/apt/lists/*

# Helper commands, like go-getter and install-protoc
ADD bin /usr/local/bin/

# Install protoc (download source, build and install -- required by protoc-gen-go and related tools)
RUN install-protoc
RUN go get -u github.com/golang/protobuf/proto
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u google.golang.org/grpc
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

RUN go get -u github.com/alecthomas/gometalinter
RUN gometalinter --install

# Custom settings for interactive shell
ADD etc/profile.d /etc/profile.d

# TODO: might want to use /go as the work directory...
WORKDIR /go/src
CMD [ "echo", "[gotools] specify the command to run" ]
