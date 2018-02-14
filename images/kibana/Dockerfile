FROM appcelerator/alpine:3.7.0

RUN apk --no-cache add nodejs-current freetype-dev fontconfig-dev

ENV ELASTIC_CONTAINER true
ENV ELASTICSEARCH_URL http://elasticsearch:9200
ENV KIBANA_MAJOR 6.2
ENV KIBANA_VERSION 6.2.1

# Kibana installation
RUN mkdir -p /opt/kibana && adduser -D -h /opt/kibana -s /sbin/nologin elastico && rm -rf /opt/kibana
RUN curl -LO https://artifacts.elastic.co/downloads/kibana/kibana-${KIBANA_VERSION}-linux-x86_64.tar.gz \
    && tar xzf /kibana-${KIBANA_VERSION}-linux-x86_64.tar.gz -C /opt \
    && mv /opt/kibana-${KIBANA_VERSION}-linux-x86_64 /opt/kibana \
    && rm /opt/kibana/node/bin/node \
    && rm /opt/kibana/node/bin/npm \
    && ln -s /usr/bin/node /opt/kibana/node/bin/node \
    && ln -s /usr/bin/npm /opt/kibana/node/bin/npm \
    && chown -R elastico:elastico /opt/kibana \
    && rm /kibana-${KIBANA_VERSION}-linux-x86_64.tar.gz /opt/kibana/config/kibana.yml
ENV PATH /opt/kibana/bin:$PATH

COPY kibana.yml.tpl /opt/kibana/config/kibana.yml.tpl
COPY run.sh /

EXPOSE 5601

ENTRYPOINT ["/sbin/tini", "--", "/run.sh"]
CMD ["kibana"]

#HEALTHCHECK --interval=5s --retries=24 --timeout=1s CMD curl -s "127.0.0.1:5601/api/status" | jq -r '.status.overall.state' | grep -q green
