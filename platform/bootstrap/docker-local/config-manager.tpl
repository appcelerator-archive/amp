{{ source "../default.ikt" }}
{{ source "file:///infrakit/env.ikt" }}
  {
      "ID": "amp-manager-{{ var "/docker/label/cluster/value" }}",
      "Properties": {
        "Allocation": {
          "LogicalIds": [
            "m1"
          ]
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
              "Name": "manager",
              "{{ var "/docker/label/cluster/key" }}": "{{ var "/docker/label/cluster/value" }}",
              "SwarmRole" : "manager",
              "ManagedBy": "InfraKit"
            }
          }
        },
        "Flavor": {
          "Plugin": "flavor-combo",
          "Properties": {
            "Flavors": [
              {
                "Plugin": "flavor-swarm/manager",
                "Properties": {
                  "InitScriptTemplateURL": "{{ var "/script/baseurl" }}/manager-init.tpl",
                  "SwarmJoinIP": "m1",
                  "Docker" : {
                    "Host" : "tcp://m1:{{ var "/docker/remoteapi/port" }}"
                  },
                  "EngineLabels": {{ var "/swarm/labels/manager" | jsonDecode | jsonEncode }}
                }
              }, {
                "Plugin": "flavor-vanilla",
                "Properties": {
                  "Init": [
                    "# create an overlay network",
                    "docker network inspect {{ var "/amp/network" }} 2>&1 | grep -q 'No such network' && \\",
                    "  docker network create -d overlay --attachable {{ var "/amp/network" }}",
                    "exit 0"
                  ]
                }
              }
            ]
          }
        }
      }
  }
