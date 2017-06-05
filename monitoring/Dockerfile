FROM appcelerator/alpine:3.5.2

ENV PROMETHEUS_VERSION 1.6.3

ENV GOLANG_VERSION 1.8.3
ENV GOLANG_SRC_URL https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz
ENV GOLANG_SRC_SHA256 5f5dea2447e7dcfdc50fa6b94c512e58bfba5673c039259fd843f68829d99fa6

RUN apk update && apk upgrade && \
    apk --virtual build-deps add build-base openssl git gcc musl musl-dev make binutils patch go && \
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
    mkdir -p /go/src/github.com/prometheus && cd /go/src/github.com/prometheus && \
    git clone https://github.com/prometheus/prometheus.git -b v${PROMETHEUS_VERSION} && \
    cd prometheus && \
    make build && \
    cp prometheus /bin/prometheus && \
    cp promtool /bin/promtool && \
    cp -pr console_libraries/ /usr/share/prometheus/ && \
    cp -pr consoles/ /usr/share/prometheus/ && \
    mkdir /etc/prometheus && \
    ln -s /usr/share/prometheus/console_libraries /usr/share/prometheus/consoles/ /etc/prometheus/ && \
    apk del build-deps && cd / && rm -rf /var/cache/apk/* $GOPATH /usr/local/go

COPY promctl.alpine /bin/promctl
COPY config/prometheus.yml  /etc/prometheus/prometheus.yml
COPY config/prometheus.tpl /etc/prometheus/prometheus.tpl

EXPOSE     9090
VOLUME     [ "/prometheus" ]
WORKDIR    /prometheus
ENTRYPOINT [ "/bin/promctl" ]
