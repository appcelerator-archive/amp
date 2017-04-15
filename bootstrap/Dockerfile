FROM docker:17.03

# TODO: remove when no longer required for bootstrap
RUN apk update && apk add bash
COPY . /usr/local/bootstrap
RUN ln -s /usr/local/bootstrap/dev /usr/local/bin/dev

ENTRYPOINT [ "dev" ]
