FROM golang:1.9.2-alpine3.6 AS builder
ENV GOPATH /go
ENV USER root
ENV CFSSL_VERSION d2393674072314fda47d2c7c16cb7fd4cdc16821
RUN set -x && \
    apk --no-cache add git gcc libc-dev && \
    mkdir -p /go/src/github.com/cloudflare && cd /go/src/github.com/cloudflare  &&\
    git clone https://github.com/cloudflare/cfssl && cd cfssl && \
    git checkout ${CFSSL_VERSION} && \
    go get github.com/GeertJohan/go.rice/rice && rice embed-go -i=./cli/serve && \
    mkdir bin && cd bin && \
    go build ../cmd/cfssl && \
    go build ../cmd/cfssljson && \
    go build ../cmd/mkbundle && \
    go build ../cmd/multirootca && \
    echo "Build complete."

FROM appcelerator/alpine:3.7.0
COPY --from=builder /go/src/github.com/cloudflare/cfssl/vendor/github.com/cloudflare/cfssl_trust /etc/cfssl
COPY --from=builder /go/src/github.com/cloudflare/cfssl/bin/ /usr/bin
VOLUME [ "/etc/cfssl" ]
WORKDIR /etc/cfssl
EXPOSE 8888
ENTRYPOINT ["/sbin/tini", "--", "/usr/bin/cfssl"]
CMD ["--help"]
