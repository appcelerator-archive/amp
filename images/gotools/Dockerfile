FROM golang:1.9-alpine
# This builds a convenience "all-in-one" image for go development.
# It intentionally does not remove any build prerequisites like most of our other
# images since this image is meant strictly for building things.

RUN rm -rf /etc/apk/repositories
RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories
RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
RUN apk update
RUN apk add bash git curl protobuf-dev

WORKDIR /tmp

# Version of the packages must be in sync with those in Gopkg.toml/Gopkg.lock
RUN mkdir -p $GOPATH/src/github.com/golang/dep && \
    git clone https://github.com/golang/dep.git $GOPATH/src/github.com/golang/dep && \
    cd $GOPATH/src/github.com/golang/dep && \
    git checkout tags/v0.3.2

RUN mkdir -p $GOPATH/src/github.com/alecthomas/gometalinter && \
    git clone https://github.com/alecthomas/gometalinter.git $GOPATH/src/github.com/alecthomas/gometalinter && \
    cd $GOPATH/src/github.com/alecthomas/gometalinter && \
    git checkout tags/v1.2.1

RUN mkdir -p $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway && \
    git clone https://github.com/grpc-ecosystem/grpc-gateway.git $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway && \
    cd $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway && \
    git checkout tags/v1.3.0

RUN mkdir -p $GOPATH/src/github.com/golang/glog && \
    git clone https://github.com/golang/glog.git $GOPATH/src/github.com/golang/glog && \
    cd $GOPATH/src/github.com/golang/glog && \
    git checkout 23def4e6c14b4da8ac2ed8007337bc5eb5007998

RUN mkdir -p $GOPATH/src/github.com/golang/protobuf && \
    git clone https://github.com/golang/protobuf.git $GOPATH/src/github.com/golang/protobuf && \
    cd $GOPATH/src/github.com/golang/protobuf && \
    git checkout 130e6b02ab059e7b717a096f397c5b60111cae74

RUN mkdir -p $GOPATH/src/google.golang.org/genproto && \
    git clone https://github.com/google/go-genproto.git $GOPATH/src/google.golang.org/genproto && \
    cd $GOPATH/src/google.golang.org/genproto && \
    git checkout f676e0f3ac6395ff1a529ae59a6670878a8371a6

RUN mkdir -p $GOPATH/src/google.golang.org/grpc/grpc-go && \
    git clone https://github.com/grpc/grpc-go.git $GOPATH/src/google.golang.org/grpc/grpc-go && \
    cd $GOPATH/src/google.golang.org/grpc/grpc-go && \
    git checkout tags/v1.5.2

RUN go install github.com/golang/dep/cmd/dep
RUN go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
RUN go install github.com/golang/protobuf/protoc-gen-go
RUN go install github.com/alecthomas/gometalinter

WORKDIR /go

CMD [ "echo", "[gotools] specify the command to run" ]
