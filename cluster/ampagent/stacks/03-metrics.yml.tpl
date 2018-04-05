version: "3.3"

networks:
  public:
    external: true
  monit:
    external: true
  core:
    external: true

volumes:
  alertmanager-data:
  grafana-data:
  prometheus-data:

configs:
  prometheus_alerts_rules:
    external: true

secrets:
  alertmanager_yml:
    external: true

services:

  prometheus:
    image: appcelerator/amp-prometheus:${TAG:-latest}
    networks:
      - public
      - monit
      - core
    volumes:
      - prometheus-data:/prometheus
{{- if  .EnableTLS }}
      - ${AMP_CERT_PATH:-/root/.docker}:/root/.docker:ro
{{- else }}
      - /var/run/docker.sock:/var/run/docker.sock:ro
{{- end }}
    environment:
      SERVICE_PORTS: 9090
      VIRTUAL_HOST: "http://alerts.*,https://alerts.*"
      PROMETHEUS_EXTERNAL_URL: "${PROMETHEUS_EXTERNAL_URL:-https://alerts.local.appcelerator.io}"
{{- if .EnableTLS }}
      DOCKER_HOST: ${AMP_HOST}
      DOCKER_TLS_VERIFY: ${AMP_TLS_VERIFY}
      DOCKER_CERT_PATH: "/root/.docker"
{{- end }}
    ports:
      - "9090:9090"
    labels:
      io.amp.role: "infrastructure"
    deploy:
      mode: replicated
      replicas: 1
      labels:
        io.amp.role: "infrastructure"
      placement:
        constraints:
        - node.labels.amp.type.metrics == true
{{- if eq .DeploymentMode "cluster" }}
      resources:
        reservations:
          cpus: '0.5'
          memory: 1.5G
{{- end }}
    configs:
      - source: prometheus_alerts_rules
        target: /etc/prometheus/alerts.rules
        mode: 0400

  cadvisor:
    image: google/cadvisor:v0.29.0
    networks:
      - core
      - monit
    labels:
      io.amp.role: "infrastructure"
    volumes:
      - /:/rootfs:ro
      - /var/run:/var/run:rw
      - /sys:/sys:ro
      - /var/lib/docker/:/var/lib/docker:ro
      #- /dev/disk:/dev/disk:ro
    # ports:
    #  - "8080:8080"
    deploy:
      mode: global
      labels:
        io.amp.role: "infrastructure"
        io.amp.metrics.port: "8080"
        io.amp.metrics.drop: "container_network_tcp_usage_total|container_network_udp_usage_total"
      resources:
        limits:
          cpus: '0.10'
          memory: 200M

  docker-engine:
    image: appcelerator/socat:1.0.0
    networks:
      - core
      - monit
    #ports:
    #  - "4999:4999"
    labels:
      io.amp.role: "infrastructure"
    deploy:
      mode: global
      labels:
        io.amp.role: "infrastructure"
        io.amp.metrics.port: "4999"
      resources:
        limits:
          cpus: '0.02'
          memory: 15M

  haproxy_exporter:
    image: prom/haproxy-exporter:v0.9.0
    command: ["--haproxy.scrape-uri", "http://stats:stats@proxy:1936/haproxy?stats;csv"]
    networks:
      - monit
      - core
    #ports:
      #- target: 9101
      #- published: 9101
    labels:
      io.amp.role: "infrastructure"
    deploy:
      mode: replicated
      replicas: 1
      labels:
        io.amp.role: "infrastructure"
        io.amp.metrics.port: "9101"
        io.amp.metrics.mode: "exporter"
      placement:
        constraints:
        - node.labels.amp.type.core == true
      resources:
        limits:
          cpus: '0.01'
          memory: 10M

  nats_exporter:
    image: appcelerator/prometheus-nats-exporter:latest
    networks:
      - monit
      - core
    command: ["-varz", "-routez", "-connz", "-subz", "nats,http://nats:8222"]
    #ports:
      #- target: 7777
      #- published: 7777
    labels:
      io.amp.role: "infrastructure"
    deploy:
      mode: replicated
      replicas: 1
      labels:
        io.amp.role: "infrastructure"
        io.amp.metrics.port: "7777"
        io.amp.metrics.mode: "exporter"
      placement:
        constraints:
        - node.labels.amp.type.core == true
      resources:
        limits:
          cpus: '0.01'
          memory: 10M

  nodes:
    image: prom/node-exporter:v0.15.2
    networks:
      - monit
      - core
    volumes:
      - /proc:/host/proc:ro
      - /sys:/host/sys:ro
      - /:/rootfs
      - /var/run/docker.sock:/var/run/docker.sock:ro
    #ports:
    #  - "9100:9100"
    command: [ "--path.procfs", "/host/proc", "--path.sysfs", "/host/sys", "--collector.filesystem.ignored-mount-points", "^/(sys|proc|dev|host|etc|var|rootfs/var/lib/docker|rootfs/run/docker/netns|rootfs/sys/kernel/debug)($$|/)"]
    labels:
      io.amp.role: "infrastructure"
    deploy:
      mode: global
      labels:
        io.amp.role: "infrastructure"
        io.amp.metrics.port: "9100"
      resources:
        limits:
          cpus: '0.03'
          memory: 15M

  alertmanager:
    image: prom/alertmanager:v0.14.0
    networks:
      - core
    volumes:
      - alertmanager-data:/alertmanager
    ports:
      - "9093:9093"
    environment:
      VIRTUAL_HOST: "https://alertmanager.*,alertmanager.*"
      SERVICE_PORTS: "9093"
    labels:
      io.amp.role: "infrastructure"
    deploy:
      mode: replicated
      replicas: 1
      labels:
        io.amp.role: "infrastructure"
      placement:
        constraints:
        - node.labels.amp.type.core == true
      resources:
        limits:
          cpus: '0.1'
          memory: 30M
    secrets:
      - source: alertmanager_yml
        target: alertmanager.yml
        mode: 0400
    command: [ "--config.file=/run/secrets/alertmanager.yml",
             "--storage.path=/alertmanager",
             "--web.external-url=http://localhost:9093" ]

  grafana:
    image: appcelerator/grafana-amp:1.2.13
    networks:
      - core
      - public
    volumes:
      - grafana-data:/var/lib/grafana
    environment:
      SERVICE_PORTS: 3000
      VIRTUAL_HOST: "http://dashboard.*,https://dashboard.*"
    ports:
      - "3000:3000"
    labels:
      io.amp.role: "infrastructure"
    deploy:
      mode: replicated
      replicas: 1
      labels:
        io.amp.role: "infrastructure"
      placement:
        constraints:
        - node.labels.amp.type.core == true
      resources:
        limits:
          cpus: '1'
          memory: 512M
{{- if eq .DeploymentMode "cluster" }}
        reservations:
          cpus: '0.05'
          memory: 50M
{{- end }}
