FROM appcelerator/grafana:grafana-4.1.1-1
ENV INFLUXDB_HOST	influxdb
ENV INFLUXDB_PASS	changeme
ENV FORCE_HOSTNAME	auto
ENV GRAFANA_PLUGIN_LIST  "grafana-piechart-panel"

COPY config /etc/extra-config/grafana
