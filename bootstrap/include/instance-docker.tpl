{
  "Plugin": "instance-docker",
  "Properties": {
    "Config": {
      "Image": "subfuzion/dind"{{ if ref "/docker/registry/host" }},
      "Cmd": ["--registry-mirror={{ ref "/docker/registry/scheme" }}{{ ref "/docker/registry/host" }}:{{ ref "/docker/registry/port" }}"]{{ end }}
    },
    "HostConfig": {
      "AutoRemove": true,
      "Privileged": true
    },
    "NetworkAttachments": [
      {
        "Name": "hostnet"
      }
    ],
    "Tags": {
      "{{ ref "/docker/label/cluster/key" }}": "{{ ref "/docker/label/cluster/value" }}"
    }
  }
}
