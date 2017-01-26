{
    "name":"telegraf",
    "type":"influxdb",
    "url":"{{ .INFLUXDB_PROTO | default "http" }}://{{ .INFLUXDB_HOST | default "localhost" }}:{{ .INFLUXDB_PORT | default "8086" }}",
    "access":"{{ .INFLUXDB_ACCESS | default "proxy" }}",
    "isDefault":true,
    "database":"{{ .GRAFANA_DB | default "telegraf" }}",
    "user":"{{ .INFLUXDB_USER | default "admin" }}",
    "password":"{{ .INFLUXDB_PASS | default "secret" }}"
}
