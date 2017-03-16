FROM docker:17.03

# TODO: remove when no longer required for bootstrap
RUN apk update && apk add bash
COPY . /usr/local/bootstrap
RUN ln -s /usr/local/bootstrap/bootstrap /usr/local/bin/bootstrap

ENTRYPOINT [ "bootstrap" ]
