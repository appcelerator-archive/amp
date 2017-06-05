{
    "name":"prometheus",
    "type":"prometheus",
    "url":"http://{{ .PROMETHEUS_HOST | default "prometheus" }}:{{ .PROMETHEUS_PORT | default "9090" }}",
    "access":"{{ .PROMETHEUS_ACCESS | default "proxy" }}",
    "isDefault":true
}
