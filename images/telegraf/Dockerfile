FROM appcelerator/alpine:3.5.1
MAINTAINER Nicolas Degory <ndegory@axway.com>

ENV TELEGRAF_VERSION 1.1.2

RUN apk update && apk upgrade && \
    apk --virtual build-deps add go git gcc musl-dev make binutils patch go && \
    export GOPATH=/go && \
    go get -v github.com/influxdata/telegraf && \
    cd $GOPATH/src/github.com/influxdata/telegraf && \
    if [ $TELEGRAF_VERSION != "master" ]; then git checkout -q --detach "${TELEGRAF_VERSION}" ; fi && \
    make && \
    chmod +x $GOPATH/bin/* && \
    mv $GOPATH/bin/* /bin/ && \
    apk del build-deps && \
    cd / && rm -rf /var/cache/apk/* $GOPATH && \
    mkdir -p /etc/telegraf

EXPOSE 8094

ENV INFLUXDB_URL http://localhost:8086
ENV INTERVAL 10s
ENV OUTPUT_INFLUXDB_ENABLED     true
ENV OUTPUT_CLOUDWATCH_ENABLED   false
ENV OUTPUT_KAFKA_ENABLED        false
ENV OUTPUT_NATS_ENABLED         false
ENV OUTPUT_FILE_ENABLED         false
ENV INPUT_KAFKA_ENABLED         false
ENV INPUT_NATS_ENABLED          false
ENV INPUT_CPU_ENABLED           true
ENV INPUT_DISK_ENABLED          true
ENV INPUT_DISKIO_ENABLED        true
ENV INPUT_KERNEL_ENABLED        false
ENV INPUT_MEM_ENABLED           true
ENV INPUT_PROCESS_ENABLED       true
ENV INPUT_SWAP_ENABLED          true
ENV INPUT_SYSTEM_ENABLED        true
ENV INPUT_NET_ENABLED           true
ENV INPUT_NETSTAT_ENABLED       false
ENV INPUT_DOCKER_ENABLED        true
ENV INPUT_LISTENER_ENABLED      false
ENV INPUT_HAPROXY_ENABLED       false

COPY telegraf.conf.tpl /etc/telegraf/telegraf.conf.tpl
COPY run.sh /run.sh

ENTRYPOINT ["/run.sh"]
CMD []

HEALTHCHECK --interval=5s --retries=3 --timeout=3s CMD pidof telegraf
