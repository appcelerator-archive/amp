FROM alpine
RUN apk --no-cache add ca-certificates
COPY gateway.alpine /usr/local/bin/gateway
ENTRYPOINT [ "gateway" ]
