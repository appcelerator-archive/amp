FROM appcelerator/alpine:3.7.0

ENV PROMETHEUS_NATS_EXPORTER_VERSION master

ENV GOLANG_VERSION 1.9.2
ENV GOLANG_SRC_URL https://storage.googleapis.com/golang/go$GOLANG_VERSION.src.tar.gz
ENV GOLANG_SRC_SHA256 665f184bf8ac89986cfd5a4460736976f60b57df6b320ad71ad4cef53bb143dc

RUN apk update && apk upgrade && \
    apk --virtual build-deps add openssl git go musl-dev make gcc && \
    echo "Installing Go" && \
    export GOROOT_BOOTSTRAP="$(go env GOROOT)" && \
    wget -q "$GOLANG_SRC_URL" -O golang.tar.gz && \
    echo "$GOLANG_SRC_SHA256  golang.tar.gz" | sha256sum -c - && \
    tar -C /usr/local -xzf golang.tar.gz && \
    rm golang.tar.gz && \
    cd /usr/local/go/src && \
    ./make.bash && \
    export GOPATH=/go && \
    export PATH=/usr/local/go/bin:$PATH && \
    go version && \
    go get -v github.com/nats-io/prometheus-nats-exporter && \
    cd $GOPATH/src/github.com/nats-io/prometheus-nats-exporter && \
    if [ "x$PROMETHEUS_NATS_EXPORTER_VERSION" != "xmaster" ]; then git checkout -q --detach "${PROMETHEUS_NATS_EXPORTER_VERSION}" ; fi && \
    go build -o /prometheus-nats-exporter && \
    apk del build-deps && \
    cd / && rm -rf /var/cache/apk/* $GOPATH /usr/local/go

EXPOSE 7777

ENTRYPOINT ["/sbin/tini", "--", "/prometheus-nats-exporter"]
