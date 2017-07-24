FROM alpine
RUN apk --no-cache add ca-certificates
COPY amplifier.alpine /usr/local/bin/amplifier
ENTRYPOINT [ "amplifier" ]
