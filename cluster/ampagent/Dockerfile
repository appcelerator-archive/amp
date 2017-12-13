FROM appcelerator/alpine:3.7.0

COPY defaults /defaults
COPY stacks /stacks
COPY ampagent.alpine /usr/local/bin/ampctl

ENTRYPOINT [ "ampctl" ]
