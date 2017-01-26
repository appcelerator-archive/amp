From appcelerator/alpine:3.5.1
RUN apk --update add docker bind-tools@edge
RUN apk add go git musl-dev && \
    export GOPATH=/go && \
    go get -v github.com/nats-io/nats && \
    cd $GOPATH/src/github.com/nats-io/nats && \
    cp examples/*.go /bin/
ENV GOPATH=/go
ENV GOBIN=/go/bin
RUN cd $GOPATH/src/github.com/nats-io/nats/examples && \
    mkdir $GOBIN && \
    go install nats-pub.go
COPY ./test*.sh /bin/
CMD ["/bin/test.sh"]
