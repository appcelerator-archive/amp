FROM alpine
RUN apk --no-cache add ca-certificates
COPY amp-local.alpine /usr/local/bin/amp-local
ENTRYPOINT [ "amp-local" ]
