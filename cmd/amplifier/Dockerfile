FROM alpine
RUN apk --no-cache add ca-certificates
COPY amplifier.alpine /usr/local/bin/amplifier
COPY amplifier.yml /etc/amp/amplifier.yml
ENTRYPOINT [ "amplifier" ]
