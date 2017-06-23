FROM docker:17.06.0-ce-rc5

RUN apk --update --no-cache add bash openssl

COPY . /opt/amp/platform
WORKDIR /opt/amp/platform
ENV PATH "/opt/amp/platform/bin:/opt/amp/platform/bootstrap:$PATH"
CMD [ "/opt/amp/platform/bin/deploy" ]

