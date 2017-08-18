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
{{- range .Jobs }}
  - job_name: '{{ .Name }}'
    metrics_path: '{{ .MetricsPath }}'
    static_configs:
{{- range .StaticConfigs }}
      - targets:
        - '{{ .Target }}:{{ .Port }}'
        labels:
{{- range $key, $value := .Labels }}
          {{ $key }}: '{{ $value }}'
{{- end }}
{{- end }}
{{- range .RelabelConfigs }}
    relabel_configs:
      - source_labels: [ {{ StringsJoin .SourceLabels ", " }} ]
        separator: '{{ .Separator }}'
        target_label: {{ .TargetLabel }}
{{- end }}
{{- end }}
  - job_name: 'haproxy'
    static_configs:
      - targets:
        - haproxy_exporter:9101
    relabel_configs:
      - replacement: haproxy
        target_label: instance
  - job_name: 'nats'
    static_configs:
      - targets:
        - nats_exporter:7777
    relabel_configs:
      - replacement: nats
        target_label: instance
{{- if .Hostnames }}
  - job_name: 'docker-engine'
    static_configs:
      - targets:
{{- range .Hostnames }}
        - '{{ . }}:{{ $dockerPort }}'
{{- end }}
    relabel_configs:
      - source_labels: [__address__]
        regex: (.*):.*
        replacement: $1
        target_label: instance
{{- end }}
{{- if .Hostnames }}
  - job_name: 'nodes'
    static_configs:
      - targets:
{{- range .Hostnames }}
        - '{{ . }}:{{ $systemPort }}'
{{- end }}
    relabel_configs:
      - source_labels: [__address__]
        regex: (.*):.*
        replacement: $1
        target_label: instance
{{- end }}
