network.host: {{ .NETWORK_HOST | default "0.0.0.0" }}
{{ if .PUBLISH_HOST }}network.publish_host: {{ .PUBLISH_HOST }}{{ end }}
http.compression: true
{{ if .UNICAST_HOSTS }}discovery.zen.minimum_master_nodes: {{ .MIN_MASTER_NODES }}
discovery.zen.ping.unicast.hosts: {{ .UNICAST_HOSTS }}{{ end }}
bootstrap.memory_lock: {{ .MEMORY_LOCK | default "false" }}
