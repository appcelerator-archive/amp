{{ source "default.ikt" }}
{{ source "file:///infrakit/env.ikt" }}
{{ $workerSize := ref "/swarm/size/worker" }}
[
  {
    "Plugin": "group",
    "Properties": {
      "ID": "amp-manager-{{ ref "/docker/label/cluster/value" }}",
      "Properties": {
        "Allocation": {
          "LogicalIds": [
            "m1",
            "m2",
            "m3"
          ]
        },
        "Instance": {
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
              "Name": "manager",
              "{{ ref "/docker/label/cluster/key" }}": "{{ ref "/docker/label/cluster/value" }}",
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
                  "InitScriptTemplateURL": "{{ ref "/script/baseurl" }}/manager-init.tpl",
                  "SwarmJoinIP": "m1",
                  "Docker" : {
                    {{ if ref "/certificate/ca/service" }}"Host" : "tcp://m1:{{ ref "/docker/remoteapi/tlsport" }}",
                    "TLS" : {
                      "CAFile": "{{ ref "/docker/remoteapi/cafile" }}",
                      "CertFile": "{{ ref "/docker/remoteapi/certfile" }}",
                      "KeyFile": "{{ ref "/docker/remoteapi/keyfile" }}",
                      "InsecureSkipVerify": false
                    }
                    {{ else }}"Host" : "tcp://m1:{{ ref "/docker/remoteapi/port" }}"{{ end }}
                  }
                }
              }, {
                "Plugin": "flavor-vanilla",
                "Properties": {
                  "Init": [
                    "# create an overlay network",
                    "docker network inspect {{ ref "/amp/network" }} 2>&1 | grep -q 'No such network' && \\",
                    "  docker network create -d overlay --attachable {{ ref "/amp/network" }}",
                    "docker run --rm --network {{ ref "/amp/network" }} alpine sh -c 'nslookup $(hostname)'",
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
      "ID": "amp-worker-{{ ref "/docker/label/cluster/value" }}",
      "Properties": {
        "Allocation": {
          "Size": {{ $workerSize }}
        },
        "Instance": {
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
              "Name": "worker",
              "Deployment": "Infrakit",
              "{{ ref "/docker/label/cluster/key" }}": "{{ ref "/docker/label/cluster/value" }}",
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
                  "InitScriptTemplateURL": "{{ ref "/script/baseurl" }}/worker-init.tpl",
                  "SwarmJoinIP": "m1",
                  "Docker" : {
                    {{ if ref "/certificate/ca/service" }}"Host" : "tcp://m1:{{ ref "/docker/remoteapi/tlsport" }}",
                    "TLS" : {
                      "CAFile": "{{ ref "/docker/remoteapi/cafile" }}",
                      "CertFile": "{{ ref "/docker/remoteapi/certfile" }}",
                      "KeyFile": "{{ ref "/docker/remoteapi/keyfile" }}",
                      "InsecureSkipVerify": false
                    }
                    {{ else }}"Host" : "tcp://m1:{{ ref "/docker/remoteapi/port" }}"{{ end }}
                  }
                }
              }, {
                "Plugin": "flavor-vanilla",
                "Properties": {
                  "Init": [
                    "docker run --rm  --network {{ ref "/amp/network" }} alpine sh -c 'nslookup $(hostname)'",
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
      "ID": "amp-proxy-{{ ref "/docker/label/cluster/value" }}",
      "Properties": {
        "Allocation": {
          "LogicalIds": [ "amp-proxy" ]
        },
        "Instance": {
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
              "Name": "worker",
              "{{ ref "/docker/label/cluster/key" }}": "{{ ref "/docker/label/cluster/value" }}",
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
                  "InitScriptTemplateURL": "{{ ref "/script/baseurl" }}/worker-init.tpl",
                  "SwarmJoinIP": "m1",
                  "Docker" : {
                    {{ if ref "/certificate/ca/service" }}"Host" : "tcp://m1:{{ ref "/docker/remoteapi/tlsport" }}",
                    "TLS" : {
                      "CAFile": "{{ ref "/docker/remoteapi/cafile" }}",
                      "CertFile": "{{ ref "/docker/remoteapi/certfile" }}",
                      "KeyFile": "{{ ref "/docker/remoteapi/keyfile" }}",
                      "InsecureSkipVerify": false
                    }
                    {{ else }}"Host" : "tcp://m1:{{ ref "/docker/remoteapi/port" }}"{{ end }}
                  }
                }
              }, {
                "Plugin": "flavor-vanilla",
                "Properties": {
                  "Init": [
                    "docker run --rm  --network {{ ref "/amp/network" }} alpine sh -c 'nslookup $(hostname)'",
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
