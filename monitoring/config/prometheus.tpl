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
{{/* services with metrics available on all tasks */}}
{{- if eq .Mode "tasks" }}
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
{{/* exporter service in front of a real service */}}
{{- else if eq .Mode "exporter" }}
  - job_name: '{{ .Name }}'
    metrics_path: '{{ .MetricsPath }}'
    static_configs:
{{- range .StaticConfigs }}
      - targets:
        - '{{ .Target }}:{{ .Port }}'
{{- end }}
{{- range .RelabelConfigs }}
    relabel_configs:
      - replacement: '{{ .Replacement }}'
        target_label: {{ .TargetLabel }}
{{- end }}
{{- end }}
{{- end }}
