FROM appcelerator/alpine:3.5.2

RUN apk --no-cache add nodejs

# Kibana installation
ENV KIBANA_MAJOR 5.3
ENV KIBANA_VERSION 5.3.0
RUN curl -LO https://artifacts.elastic.co/downloads/kibana/kibana-${KIBANA_VERSION}-linux-x86_64.tar.gz \
    && mkdir /opt \
    && tar xzf /kibana-${KIBANA_VERSION}-linux-x86_64.tar.gz -C /opt \
    && mv /opt/kibana-${KIBANA_VERSION}-linux-x86_64 /opt/kibana \
    && rm /opt/kibana/node/bin/node \
    && rm /opt/kibana/node/bin/npm \
    && ln -s /usr/bin/node /opt/kibana/node/bin/node \
    && ln -s /usr/bin/npm /opt/kibana/node/bin/npm \
    && rm /kibana-${KIBANA_VERSION}-linux-x86_64.tar.gz

ENV PATH /opt/kibana/bin:$PATH
ENV ELASTICSEARCH_URL http://elasticsearch:9200

COPY kibana.yml.tpl /opt/kibana/config/kibana.yml.tpl
RUN rm -v /opt/kibana/config/kibana.yml

#COPY preconfiguration.sh /
COPY run.sh /

EXPOSE 5601
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["/run.sh"]

HEALTHCHECK --interval=5s --retries=24 --timeout=1s CMD curl -s "127.0.0.1:5601/api/status" | jq -r '.status.overall.state' | grep -q green
