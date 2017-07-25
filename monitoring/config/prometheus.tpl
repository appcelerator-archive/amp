{{ $dockerPort := .DockerEngineMetricsPort -}}
{{ $systemPort := .SystemMetricsPort -}}
global:
  scrape_interval:     15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

  # Attach these labels to any time series or alerts when communicating with
  # external systems (federation, remote storage, Alertmanager).
  external_labels:
      monitor: 'amp'

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  - "/etc/prometheus/*.rules"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets:
        - localhost:9090
  - job_name: 'etcd'
    dns_sd_configs:
      - names:
        - 'tasks.etcd'
        type: 'A'
        port: 2379
  - job_name: 'haproxy'
    static_configs:
      - targets:
        - haproxy_exporter:9101
  - job_name: 'nats'
    static_configs:
      - targets:
        - nats_exporter:7777
  - job_name: 'elasticsearch'
    metrics_path: "/_prometheus/metrics"
    dns_sd_configs:
      - names:
        - 'tasks.elasticsearch'
        type: 'A'
        port: 9200
  - job_name: 'amplifier'
    dns_sd_configs:
      - names:
        - 'tasks.amplifier'
        type: 'A'
        port: 5100
{{- if .Hostnames }}
  - job_name: 'docker-engine'
    static_configs:
      - targets:
{{- range .Hostnames }}
        - '{{ . }}:{{ $dockerPort }}'
{{- end }}
{{- end }}
{{- if .Hostnames }}
  - job_name: 'nodes'
    static_configs:
      - targets:
{{- range .Hostnames }}
        - '{{ . }}:{{ $systemPort }}'
{{- end }}
{{- end }}
