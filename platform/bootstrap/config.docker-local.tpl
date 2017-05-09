{{ source "default.ikt" }}
{{ source "file:///infrakit/env.ikt" }}
{{ $cattleWorkerSize := var "/swarm/size/worker/cattle" }}
[
  {
    "Plugin": "group",
    "Properties": {
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
              "Image": "subfuzion/dind:17.05-ce-rc1"{{ if var "/docker/registry/host" }},
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
              "SwarmRole" : "manager"
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
                  }
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
  },
  {
    "Plugin": "group",
    "Properties": {
      "ID": "amp-worker-cattle-{{ var "/docker/label/cluster/value" }}",
      "Properties": {
        "Allocation": {
          "Size": {{ $cattleWorkerSize }}
        },
        "Instance": {
          "Plugin": "instance-docker",
          "Properties": {
            "Config": {
              "Image": "subfuzion/dind:17.05-ce-rc1"{{ if var "/docker/registry/host" }},
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
              "Name": "worker",
              "Deployment": "Infrakit",
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
                  }
                }
              }, {
                "Plugin": "flavor-vanilla",
                "Properties": {
                  "Init": [
                    "docker run --rm  --network {{ var "/amp/network" }} alpine sh -c 'nslookup $(hostname)'",
                    "if [ $? -ne 0 ]; then echo 'Docker Swarm DNS check failed'; exit 1; fi",
                    "exit 0"
                  ]
                }
              }
            ]
          }
        }
      }
    }
  },
  {
    "Plugin": "group",
    "Properties": {
      "ID": "amp-proxy-{{ var "/docker/label/cluster/value" }}",
      "Properties": {
        "Allocation": {
          "LogicalIds": [ "amp-proxy" ]
        },
        "Instance": {
          "Plugin": "instance-docker",
          "Properties": {
            "Config": {
              "Image": "subfuzion/dind:17.05-ce-rc1"{{ if var "/docker/registry/host" }},
              "Cmd": ["--registry-mirror={{ var "/docker/registry/scheme" }}{{ var "/docker/registry-cache/host" }}:{{ var "/docker/registry-cache/port" }}", "--registry-mirror={{ var "/docker/registry/scheme" }}{{ var "/docker/registry/host" }}:{{ var "/docker/registry/port" }}"]{{ end }} {{ if var "/docker/ports/exposed" }},
              "ExposedPorts": {{ var "/docker/ports/exposed" | jsonEncode }} {{ end }}
            },
            "HostConfig": {
              "AutoRemove": true,
              "Privileged": true{{ if var "/docker/ports/bindings" }},
              "PortBindings": {{ var "/docker/ports/bindings" | jsonEncode }} {{ end }}
            },
            "NetworkAttachments": [
              {
                "Name": "hostnet"
              }
            ],
            "Tags": {
              "Name": "worker",
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
                  "EngineLabels": { "proxy": "true" },
                  "Docker" : {
                    "Host" : "tcp://m1:{{ var "/docker/remoteapi/port" }}"
                  }
                }
              }, {
                "Plugin": "flavor-vanilla",
                "Properties": {
                  "Init": [
                    "docker run --rm  --network {{ var "/amp/network" }} alpine sh -c 'nslookup $(hostname)'",
                    "if [ $? -ne 0 ]; then echo 'Docker Swarm DNS check failed'; exit 1; fi",
                    "exit 0"
                  ]
                }
              }
            ]
          }
        }
      }
    }
  }
]
