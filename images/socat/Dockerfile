# based on https://github.com/bvis/docker-socat
FROM alpine:3.7
RUN echo "@community http://nl.alpinelinux.org/alpine/v3.7/community" >> /etc/apk/repositories
ENV IN "9323"
ENV OUT "4999"
RUN apk add --no-cache socat tini@community
COPY entrypoint.sh /bin/
ENTRYPOINT ["/sbin/tini", "--", "/bin/entrypoint.sh"]
