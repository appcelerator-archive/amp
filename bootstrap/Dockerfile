FROM docker:17.03

COPY . /usr/local/bootstrap
WORKDIR /usr/local/bootstrap
ENV PATH "/usr/local/bootstrap:/usr/local/bootstrap/hack:$PATH"

ENTRYPOINT [ "dev" ]
