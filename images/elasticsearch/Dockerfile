FROM appcelerator/alpine:3.7.0

RUN apk --no-cache add openjdk8-jre bind-tools

ENV PATH /bin:/opt/elasticsearch/bin:$PATH
ENV ELASTIC_CONTAINER true
ENV ELASTIC_VERSION 6.2.1

RUN mkdir -p /opt/elasticsearch && adduser -D -h /opt/elasticsearch -s /sbin/nologin elastico && rm -rf /opt

COPY config /_config
ENV ES_TMPDIR /tmp/_elasticsearch${ELASTIC_VERSION}

RUN curl -L https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-$ELASTIC_VERSION.tar.gz -o /tmp/elasticsearch-$ELASTIC_VERSION.tar.gz && \
    mkdir /opt && \
    echo "install elasticsearch" && \
    tar xzf /tmp/elasticsearch-$ELASTIC_VERSION.tar.gz -C /opt && \
    mv /opt/elasticsearch-$ELASTIC_VERSION /opt/elasticsearch && \
    rm -f /opt/elasticsearch/bin/elasticsearch*exe /opt/elasticsearch/bin/elasticsearch*bat && \
    rm /tmp/elasticsearch-$ELASTIC_VERSION.tar.gz && \
    mv /_config/* /opt/elasticsearch/config/ && rm -rf /_config && \
    mkdir -p /opt/elasticsearch/config/scripts && \
    echo "install prometheus plugin" && \
    /opt/elasticsearch/bin/elasticsearch-plugin install -b https://distfiles.compuscene.net/elasticsearch/elasticsearch-prometheus-exporter-${ELASTIC_VERSION}.0.zip && \
    chown -R elastico:elastico /opt/elasticsearch

COPY /bin/docker-entrypoint.sh /bin/

VOLUME /opt/elasticsearch/data

EXPOSE 9200 9300
#ENV JAVA_HEAP_SIZE=256
ENV MIN_MASTER_NODES 1

#HEALTHCHECK --interval=15s --retries=3 --timeout=5s CMD curl -s 127.0.0.1:9200 | jq .version.number | grep -qv null

ENTRYPOINT ["/sbin/tini", "--", "/bin/docker-entrypoint.sh"]
CMD ["elasticsearch"]
