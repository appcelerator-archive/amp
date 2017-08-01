FROM golang:1.8-alpine as buildenv2yaml

RUN apk add --no-cache build-base ca-certificates git
RUN go get gopkg.in/yaml.v2
COPY env2yaml /go/src/env2yaml
RUN cd /go/src && go build -o /bin/env2yaml env2yaml

FROM appcelerator/alpine:3.6.0 as buildlogstash
ENV ELASTIC_VERSION 5.5.0
RUN curl -L https://artifacts.elastic.co/downloads/logstash/logstash-$ELASTIC_VERSION.tar.gz -o /tmp/logstash-$ELASTIC_VERSION.tar.gz && \
    tar xzf /tmp/logstash-$ELASTIC_VERSION.tar.gz -C /usr/share && \
    mv /usr/share/logstash-$ELASTIC_VERSION /usr/share/logstash && \
    rm -f /usr/share/logstash/bin/logstash*exe /usr/share/logstash/bin/logstash*bat && \
    rm /tmp/logstash-$ELASTIC_VERSION.tar.gz
ADD pipeline/ /usr/share/logstash/pipeline/
ADD config/ /usr/share/logstash/config/
ADD template/ /usr/share/logstash/template/
RUN adduser -D -h /usr/share/logstash -s /sbin/nologin elastico
RUN chown -R elastico:elastico /usr/share/logstash

FROM appcelerator/alpine:3.6.0
RUN apk update && apk upgrade && apk add openjdk8-jre bind-tools && rm -rf /var/cache/apk/*
ENV PATH /bin:/usr/share/logstash/bin:$PATH
ENV ELASTIC_VERSION 5.5.0
RUN adduser -D -h /usr/share/logstash -s /sbin/nologin elastico
RUN chown -R elastico:elastico /usr/share/logstash
COPY --from=buildenv2yaml /bin/env2yaml /usr/local/bin/env2yaml
COPY --from=buildlogstash /usr/share/logstash /usr/share/logstash
WORKDIR /usr/share/logstash
COPY /bin/docker-entrypoint.sh /usr/local/bin/
VOLUME /usr/share/logstash/data
EXPOSE 5044 9600
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["-f", "/usr/share/logstash/pipeline"]
