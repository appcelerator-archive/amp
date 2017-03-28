FROM alpine
RUN apk --no-cache add ca-certificates
COPY amplifier.alpine /usr/local/bin/amplifier
COPY amplifier.yaml /root/.config/amp/amplifier.yaml
ENTRYPOINT [ "amplifier" ]
