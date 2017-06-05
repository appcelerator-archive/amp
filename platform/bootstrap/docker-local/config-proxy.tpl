{{ source "../default.ikt" }}
{{ source "file:///infrakit/env.ikt" }}
  {
      "ID": "amp-proxy-{{ var "/docker/label/cluster/value" }}",
      "Properties": {
        "Allocation": {
          "LogicalIds": [ "amp-proxy" ]
        },
        "Instance": {
          "Plugin": "instance-docker",
          "Properties": {
            "Config": {
              "Image": "subfuzion/dind:17.05.0"{{ if var "/docker/registry/host" }},
              "Cmd": ["--registry-mirror={{ var "/docker/registry/scheme" }}{{ var "/docker/registry-cache/host" }}:{{ var "/docker/registry-cache/port" }}", "--registry-mirror={{ var "/docker/registry/scheme" }}{{ var "/docker/registry/host" }}:{{ var "/docker/registry/port" }}"]{{ end }} {{ if var "/docker/ports/exposed" }},
              "ExposedPorts": {{ var "/docker/ports/exposed" | jsonDecode | jsonEncode }} {{ end }}
            },
            "HostConfig": {
              "AutoRemove": true,
              "Privileged": true{{ if var "/docker/ports/bindings" }},
              "PortBindings": {{ var "/docker/ports/bindings" | jsonDecode | jsonEncode }} {{ end }}
            },
            "NetworkAttachments": [
              {
                "Name": "hostnet"
              }
            ],
            "Tags": {
              "Name": "worker",
              "{{ var "/docker/label/cluster/key" }}": "{{ var "/docker/label/cluster/value" }}",
              "SwarmRole" : "worker",
              "WorkerType" : "proxy",
              "ManagedBy": "InfraKit"
            }
          }
        },
        "Flavor": {
          "Plugin": "flavor-combo",
          "Properties": {
            "Flavors": [
              {
                "Plugin": "flavor-swarm/worker",
                "Properties": {
                  "InitScriptTemplateURL": "{{ var "/script/baseurl" }}/worker-init.tpl",
                  "SwarmJoinIP": "m1",
                  "EngineLabels": { "proxy": "true", "amp.type.metrics": "true" },
                  "Docker" : {
                    "Host" : "tcp://m1:{{ var "/docker/remoteapi/port" }}"
                  },
                  "EngineLabels": {{ var "/swarm/labels/worker/proxy" | jsonDecode | jsonEncode }}
                }
              }
            ]
          }
        }
      }
  }
