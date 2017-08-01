FROM alpine:3.6
COPY ampbeat.alpine /beat/ampbeat
COPY ampbeat.*.json /beat/
COPY ampbeat.yml /beat/ampbeat.yml
RUN chmod go-w /beat/ampbeat.yml
COPY _meta /beat/_meta/
WORKDIR /beat
ENTRYPOINT [ "/beat/ampbeat", "-e" ]
