From appcelerator/alpine:3.6.0
ENV GOPATH          /go
ENV GOBIN           /go/bin
RUN apk update && apk upgrade && \
    apk add go git musl musl-dev && \
    go get -v github.com/nats-io/nats && \
    cd $GOPATH/src/github.com/nats-io/nats && \
    cp examples/*.go /bin/ && \
    rm -rf /var/cache/apk/*
RUN cd $GOPATH/src/github.com/nats-io/nats/examples && \
    mkdir $GOBIN && \
    go install nats-pub.go

COPY ./test.sh /bin
CMD ["/bin/test.sh"]
