FROM appcelerator/prometheus:2.0.0_alpine37

COPY promctl.alpine /bin/promctl
COPY config/prometheus.yml  /etc/prometheus/prometheus.yml
COPY config/prometheus.tpl /etc/prometheus/prometheus.tpl

ENTRYPOINT [ "/bin/promctl" ]
CMD [ ]
