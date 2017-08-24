# based on https://github.com/bvis/docker-socat
FROM alpine:3.6
ENV IN "9323"
ENV OUT "4999"
RUN apk add --no-cache socat
COPY entrypoint.sh /bin/
ENTRYPOINT /bin/entrypoint.sh
