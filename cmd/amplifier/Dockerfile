FROM alpine
RUN apk --no-cache add ca-certificates
COPY amplifier.alpine /usr/local/bin/amplifier
COPY amplifier.yaml /etc/atomiq/amplifier.yaml
ENTRYPOINT [ "amplifier" ]
