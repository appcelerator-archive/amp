FROM docker:17.03

# TODO: remove when no longer required for bootstrap
RUN apk update && apk add bash
COPY . /usr/local/bootstrap
RUN ln -s /usr/local/bootstrap/dev /usr/local/bin/dev
RUN ln -s /usr/local/bootstrap/deploy /usr/local/bin/deploy
RUN ln -s /usr/local/bootstrap/bootstrap /usr/local/bin/bootstrap

ENTRYPOINT [ "dev" ]
