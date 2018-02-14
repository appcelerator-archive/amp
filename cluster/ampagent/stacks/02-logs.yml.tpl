version: "3.3"

networks:
  public:
    external: true
  monit:
    external: true
  core:
    external: true

volumes:
  elasticsearch-data:
  ampagent:

services:

  elasticsearch:
    image: appcelerator/elasticsearch-amp:6.2.1
    networks:
      - monit
      - core
    volumes:
      - elasticsearch-data:/opt/elasticsearch/data
    labels:
      io.amp.role: "infrastructure"
    environment:
{{- if eq .DeploymentMode "cluster" }}
      MIN_MASTER_NODES: 2
      UNICAST_HOSTS: "tasks.elasticsearch"
{{- end }}
      NETWORK_HOST: "_site_"
      JAVA_HEAP_SIZE: "${ES_JAVA_HEAP_SIZE:-1024}"
    deploy:
      mode: replicated
{{- if eq .DeploymentMode "cluster" }}
      replicas: 3
      update_config:
        parallelism: 1
        delay: 120s
      restart_policy:
        condition: any
        delay: 5s
        window: 25s
{{- else }}
      replicas: 1
{{- end }}
      labels:
        io.amp.role: "infrastructure"
        io.amp.metrics.port: "9200"
        io.amp.metrics.path: "/_prometheus/metrics"
      placement:
        constraints:
        - node.labels.amp.type.search == true
{{- if eq .DeploymentMode "cluster" }}
      resources:
        limits:
          cpus: '1'
        reservations:
          cpus: '0.5'
          memory: '2G'
{{- end }}

  nats:
    image: appcelerator/amp-nats-streaming:v0.7.0
    networks:
      - core
    labels:
      io.amp.role: "infrastructure"
    deploy:
      mode: replicated
      replicas: 1
      labels:
        io.amp.role: "infrastructure"
      placement:
        constraints:
        - node.labels.amp.type.mq == true
{{- if eq .DeploymentMode "cluster" }}
      resources:
        limits:
          cpus: '1.5'
        reservations:
          cpus: '0.4'
          memory: '512M'
{{- end }}

  ampbeat:
    image: appcelerator/ampbeat:${TAG:-latest}
    networks:
      - core
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
          cpus: '0.01'
          memory: '20M'
    labels:
      io.amp.role: "infrastructure"

  kibana:
    image: appcelerator/kibana:6.2.1
    networks:
      - core
      - public
    deploy:
      mode: replicated
      replicas: 1
      labels:
        io.amp.role: "infrastructure"
        io.amp.mapping: "kibana:5601"
      placement:
        constraints:
        - node.labels.amp.type.core == true
      resources:
        limits:
          cpus: '1'
          memory: 200M
{{- if eq .DeploymentMode "cluster" }}
        reservations:
          cpus: '0.05'
          memory: 200M
{{- end }}
    labels:
      io.amp.role: "infrastructure"
    environment:
      ELASTICSEARCH_URL: "http://elasticsearch:9200"
      SERVICE_PORTS: 5601
      VIRTUAL_HOST: "http://kibana.*,https://kibana.*"

  agent:
    image: appcelerator/agent:${TAG:-latest}
    networks:
      - core
    deploy:
      mode: global
      labels:
        io.amp.role: "infrastructure"
      resources:
        limits:
          cpus: '0.05'
          memory: 15M
    labels:
      io.amp.role: "infrastructure"
    volumes:
      - ampagent:/containers
      - /var/run/docker.sock:/var/run/docker.sock
