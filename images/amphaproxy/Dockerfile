FROM haproxy:1.7-alpine

ENV GOPATH /go
ENV PATH $PATH:/go/bin

RUN echo "@edge http://nl.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories && \
    echo "@testing http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
    echo "@community http://nl.alpinelinux.org/alpine/v3.5/community" >> /etc/apk/repositories
RUN apk update && apk upgrade && \
    apk -v add curl && \
    cd / && rm -rf /go/src /go/pkg /var/cache/apk/* /root/.cache /root/.glide

COPY haproxy.cfg.tpt /usr/local/etc/haproxy/haproxy.cfg.tpt
# amp-haproxy should be built outside from this Dockerfile, see ../build-images.sh
COPY amphaproxy.alpine /go/bin/amphaproxy

HEALTHCHECK --interval=5s --timeout=10s --retries=12 CMD curl http://127.0.0.1/healthcheck

CMD ["/go/bin/amphaproxy"]
