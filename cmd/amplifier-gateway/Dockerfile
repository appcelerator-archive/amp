FROM alpine
RUN apk --no-cache add ca-certificates
COPY amplifier-gateway.alpine /usr/local/bin/amplifier-gateway
ENTRYPOINT [ "amplifier-gateway" ]
