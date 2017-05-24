{{ source "../default.ikt" }}
{{ source "file:///infrakit/env.ikt" }}
  {
      "ID": "amp-worker-cattle-{{ var "/docker/label/cluster/value" }}",
      "Properties": {
        "Allocation": {
          "Size": {{ var "/swarm/size/worker/gp" }}
        },
        "Instance": {
          "Plugin": "instance-docker",
          "Properties": {
            "Config": {
              "Image": "subfuzion/dind:17.05.0"{{ if var "/docker/registry/host" }},
              "Cmd": ["--registry-mirror={{ var "/docker/registry/scheme" }}{{ var "/docker/registry-cache/host" }}:{{ var "/docker/registry-cache/port" }}", "--registry-mirror={{ var "/docker/registry/scheme" }}{{ var "/docker/registry/host" }}:{{ var "/docker/registry/port" }}"]{{ end }}
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
              "Name": "worker-gp",
              "WorkerType": "gp",
              "ManagedBy": "Infrakit",
              "{{ var "/docker/label/cluster/key" }}": "{{ var "/docker/label/cluster/value" }}",
              "SwarmRole" : "worker"
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
                  "Docker" : {
                    "Host" : "tcp://m1:{{ var "/docker/remoteapi/port" }}"
                  },
                  "EngineLabels": {{ var "/swarm/labels/worker/gp" | jsonDecode | jsonEncode }}
                }
              }
            ]
          }
        }
      }
  }
