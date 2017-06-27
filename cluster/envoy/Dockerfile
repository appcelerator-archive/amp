FROM alpine
RUN apk add --no-cache ca-certificates
COPY envoy.alpine /usr/local/bin/envoy
ENTRYPOINT [ "envoy" ]
