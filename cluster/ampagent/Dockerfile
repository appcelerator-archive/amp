FROM alpine

RUN apk --no-cache add ca-certificates

COPY defaults /defaults
COPY stacks /stacks
COPY ampagent.alpine /usr/local/bin/ampctl

ENTRYPOINT [ "ampctl" ]
