FROM docker:17.06.1-ce-rc1

RUN apk --update --no-cache add bash openssl curl

COPY . /opt/amp/platform
WORKDIR /opt/amp/platform
ENV PATH "/opt/amp/platform/bin:/opt/amp/platform/bootstrap:$PATH"
CMD [ "/opt/amp/platform/bin/deploy" ]

