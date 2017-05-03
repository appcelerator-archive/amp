FROM docker:17.05.0-ce-rc1

RUN apk --update --no-cache add bash openssl

WORKDIR /usr/local
COPY . /usr/local/bootstrap
RUN mv /usr/local/bootstrap/hack /usr/local/hack && mv /usr/local/bootstrap/stacks /usr/local/stacks

ENV PATH "/usr/local/bootstrap:/usr/local/hack:$PATH"
CMD [ "hack/deploy" ]
