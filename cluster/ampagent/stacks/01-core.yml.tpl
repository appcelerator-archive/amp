version: "3.3"

networks:
  public:
    external: true
  core:
    external: true
  monit:
    external: true

secrets:
  amplifier_yml:
    external: true
  certificate_amp:
    external: true

volumes:
  etcd-data:

services:

  etcd:
    image: appcelerator/etcd:3.2.13
    networks:
      - core
      - monit
    volumes:
      - etcd-data:/data
    environment:
      SERVICE_NAME: "etcd"
      MIN_SEEDS_COUNT: {{ if eq .DeploymentMode "cluster" }}3{{ else }}1{{ end }}
    command:
      - "--advertise-client-urls"
      - "http://etcd:2379"
    labels:
      io.amp.role: "infrastructure"
    deploy:
      mode: replicated
{{- if eq .DeploymentMode "cluster" }}
      replicas: 3
      update_config:
        parallelism: 1
        delay: 30s
      restart_policy:
        condition: any
        delay: 25s
        window: 15s
{{- else }}
      replicas: 1
{{- end }}
      labels:
        io.amp.role: "infrastructure"
        io.amp.metrics.port: "2379"
      placement:
        constraints:
        - node.labels.amp.type.kv == true
{{- if eq .DeploymentMode "cluster" }}
      resources:
        reservations:
          memory: '80M'
{{- end }}

  amplifier:
    image: appcelerator/amplifier:${TAG:-latest}
    networks:
      - core
      - monit
    environment:
      REGISTRATION: ${REGISTRATION:-email}
      NOTIFICATIONS: ${NOTIFICATIONS:-true}
{{- if .EnableTLS }}
      DOCKER_HOST: "${AMP_HOST:-unix:///var/run/docker.sock}"
      DOCKER_TLS_VERIFY: "${AMP_TLS_VERIFY}"
      DOCKER_CERT_PATH: "/root/.docker"
{{- end }}
    ports:
      - "50101:50101"
    volumes:
{{- if .EnableTLS }}
      - ${AMP_CERT_PATH:-/root/.docker}:/root/.docker:ro
{{- end }}
      - "/var/run/docker.sock:/var/run/docker.sock"
    labels:
      io.amp.role: "infrastructure"
    deploy:
      mode: global
      labels:
        io.amp.role: "infrastructure"
        io.amp.metrics.port: "5100"
      restart_policy:
        condition: on-failure
      placement:
        constraints:
        - node.labels.amp.type.api == true
      resources:
        limits:
          memory: '500M'
{{- if eq .DeploymentMode "cluster" }}
        reservations:
          cpus: '0.15'
          memory: '150M'
{{- end }}
    secrets:
      - source: amplifier_yml
        target: amplifier.yml
        mode: 0400
      - source: certificate_amp
        target: cert0.pem
        mode: 0400

  gateway:
    image: appcelerator/gateway:${TAG:-latest}
    networks:
      - core
      - public
    labels:
      io.amp.role: "infrastructure"
    environment:
      SERVICE_PORTS: 80
      VIRTUAL_HOST: "https://gw.*,http://gw.*"
    deploy:
      mode: global
      labels:
        io.amp.role: "infrastructure"
      restart_policy:
        condition: on-failure
      placement:
        constraints:
        - node.labels.amp.type.core == true
      resources:
        limits:
          memory: '100M'
{{- if eq .DeploymentMode "cluster" }}
        reservations:
          cpus: '0.05'
          memory: '20M'
{{- end }}
