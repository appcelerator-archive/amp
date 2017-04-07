{
  "Plugin": "instance-docker",
  "Properties": {
    "Config": {
      "Image": "subfuzion/dind"{{ if ref "/docker/registry/host" }},
      "Cmd": "--registry-mirror={{ ref "/docker/registry/scheme" }}{{ ref "/docker/registry/host" }}:{{ ref "/docker/registry/port" }}"{{ end }} {{ if ref "/docker/ports/exposed" }},
      "ExposedPorts": {{ ref "/docker/ports/exposed" | to_json }} {{ end }}
    },
    "HostConfig": {
      "AutoRemove": true,
      "Privileged": true{{ if ref "/docker/ports/bindings" }},
      "PortBindings": {{ ref "/docker/ports/bindings" | to_json }} {{ end }}
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
