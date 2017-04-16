FROM docker:17.03

RUN apk --update --no-cache add bash openssl

WORKDIR /usr/local
COPY . /usr/local/bootstrap
RUN mv bootstrap/hack hack
RUN mv bootstrap/stacks stacks

ENV PATH "/usr/local/bootstrap:/usr/local/hack:$PATH"
ENTRYPOINT [ "hack/dev" ]

