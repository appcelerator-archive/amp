FROM scratch

COPY defaults /defaults
COPY stacks /stacks
COPY ampagent.alpine /usr/local/bin/ampctl

ENTRYPOINT [ "ampctl" ]
