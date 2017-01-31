FROM appcelerator/alpine:3.5.1

RUN apk update && apk upgrade && apk --no-cache add openjdk8-jre

ENV PATH /bin:/opt/elasticsearch/bin:$PATH
ENV ELASTIC_VERSION 5.1.2

RUN curl -L https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-$ELASTIC_VERSION.tar.gz -o /tmp/elasticsearch-$ELASTIC_VERSION.tar.gz && \
    mkdir /opt && \
    tar xzf /tmp/elasticsearch-$ELASTIC_VERSION.tar.gz -C /opt && \
    ln -s /opt/elasticsearch-$ELASTIC_VERSION /opt/elasticsearch && \
    rm -f /opt/elasticsearch/bin/elasticsearch*exe /opt/elasticsearch/bin/elasticsearch*bat && \
    rm /tmp/elasticsearch-$ELASTIC_VERSION.tar.gz

WORKDIR /opt/elasticsearch

COPY config/java.policy /opt/elasticsearch/config/java.policy
COPY config/elasticsearch.yml /opt/elasticsearch/config/elasticsearch.yml.tpl
COPY config/log4j2.properties /opt/elasticsearch/config/
COPY /bin/docker-entrypoint.sh /bin/

RUN mkdir -p /opt/elasticsearch/config/scripts
RUN adduser -D -h /opt/elasticsearch -s /sbin/nologin elastico
RUN chown -R elastico:elastico /opt/elasticsearch

VOLUME /opt/elasticsearch-$ELASTIC_VERSION/data

EXPOSE 9200 9300
#ENV JAVA_HEAP_SIZE=256

HEALTHCHECK --interval=15s --retries=3 --timeout=5s CMD curl -s 127.0.0.1:9200 | jq .version.number | grep -qv null

ENTRYPOINT ["/bin/docker-entrypoint.sh"]
CMD ["elasticsearch"]
