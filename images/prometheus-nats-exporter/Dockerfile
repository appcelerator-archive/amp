FROM appcelerator/alpine:3.5.2

#ENV PROMETHEUS_NATS_EXPORTER_VERSION 1.0.0

ENV GOLANG_VERSION 1.8.3
ENV GOLANG_SRC_URL https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz
ENV GOLANG_SRC_SHA256 5f5dea2447e7dcfdc50fa6b94c512e58bfba5673c039259fd843f68829d99fa6

RUN apk update && apk upgrade && \
    apk --virtual build-deps add openssl git gcc musl-dev make binutils patch go && \
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
    if [ $PROMETHEUS_NATS_EXPORTER_VERSION != "master" ]; then git checkout -q --detach "${PROMETHEUS_NATS_EXPORTER_VERSION}" ; fi && \
    go build -o /prometheus-nats-exporter && \
    apk del build-deps && \
    cd / && rm -rf /var/cache/apk/* $GOPATH /usr/local/go

EXPOSE 7777

ENTRYPOINT ["/prometheus-nats-exporter"]
