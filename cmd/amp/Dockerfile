FROM docker:1.13

# TODO: remove when no longer required for local-bootstrap
RUN apk update && apk add bash
COPY local-bootstrap /usr/local/bin/

COPY amp.alpine /usr/local/bin/amp
ENTRYPOINT [ "amp" ]

