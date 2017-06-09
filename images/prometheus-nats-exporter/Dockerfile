FROM appcelerator/alpine:3.6.0

#ENV PROMETHEUS_NATS_EXPORTER_VERSION 1.0.0

RUN apk update && apk upgrade && \
    apk --virtual build-deps add openssl git go musl-dev && \
    export GOPATH=/go && \
    go version && \
    go get -v github.com/nats-io/prometheus-nats-exporter && \
    cd $GOPATH/src/github.com/nats-io/prometheus-nats-exporter && \
    if [ $PROMETHEUS_NATS_EXPORTER_VERSION != "master" ]; then git checkout -q --detach "${PROMETHEUS_NATS_EXPORTER_VERSION}" ; fi && \
    go build -o /prometheus-nats-exporter && \
    apk del build-deps && \
    cd / && rm -rf /var/cache/apk/* $GOPATH

EXPOSE 7777

ENTRYPOINT ["/prometheus-nats-exporter"]
