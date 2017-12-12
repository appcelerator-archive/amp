FROM alpine:3.7
RUN apk --no-cache add ca-certificates
COPY amplifier.alpine /usr/local/bin/amplifier
ENTRYPOINT [ "amplifier" ]
