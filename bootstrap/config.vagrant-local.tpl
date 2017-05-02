{{ source "default.ikt" }}
{{ source "env.ikt" }}
{{ $workerSize := var "/swarm/size/worker" }}
[
  {
    "Plugin": "group",
    "Properties": {
      "ID": "amp-manager",
      "Properties": {
        "Allocation": {
          "LogicalIds": [
            "{{ var "/m1/ip" }}",
            "{{ var "/m2/ip" }}",
            "{{ var "/m3/ip" }}"
          ]
        },
        "Instance": {
          "Plugin": "instance-vagrant",
          "Properties": {
            "Box": "ubuntu/xenial64"
          }
        },
        "Flavor": {
          "Plugin": "flavor-combo",
          "Properties": {
            "Flavors": [
              {
                "Plugin": "flavor-vanilla",
                "Properties": {
                  "Init": [
                    "#!/bin/bash",
                    "apt-get install -y awscli jq"
                  ]
                }
              }, {
                "Plugin": "flavor-swarm/manager",
                "Properties": {
                  "InitScriptTemplateURL": "{{ var "/script/baseurl" }}/manager-init.vagrant.tpl",
                  "SwarmJoinIP": "{{ var "/m1/ip" }}",
                  "Docker" : {
                    {{ if var "/certificate/ca/service" }}"Host" : "tcp://{{ var "/m1/ip" }}:{{ var "/docker/remoteapi/tlsport" }}",
                    "TLS" : {
                      "CAFile": "{{ var "/docker/remoteapi/cafile" }}",
                      "CertFile": "{{ var "/docker/remoteapi/certfile" }}",
                      "KeyFile": "{{ var "/docker/remoteapi/keyfile" }}",
                      "InsecureSkipVerify": false
                    }
                    {{ else }}"Host" : "tcp://{{ var "/m1/ip" }}:{{ var "/docker/remoteapi/port" }}"
                    {{ end }}
                  }
                }
              }, {
                "Plugin": "flavor-vanilla",
                "Properties": {
                  "Init": [
                    "set -o errexit",
                    "docker network inspect {{ var "/amp/network" }} 2>&1 | grep -q 'No such network' && \\",
                    "  docker network create -d overlay --attachable {{ var "/amp/network" }}",
                    "docker service ls {{ var "/amp/network" }} 2>&1 | grep -q 'No such network' && \\",
                    "docker service create --name amplifier --network {{ var "/amp/network" }} {{ var "/amp/amplifier/image" }}:{{ var "/amp/amplifier/version" }} || true"
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
      "ID": "amp-worker",
      "Properties": {
        "Allocation": {
          "Size": {{ $workerSize }}
        },
        "Instance": {
          "Plugin": "instance-vagrant",
          "Properties": {
            "Box": "ubuntu/xenial64"
          }
        },
        "Flavor": {
          "Plugin": "flavor-combo",
          "Properties": {
            "Flavors": [
              {
                "Plugin": "flavor-vanilla",
                "Properties": {
                  "Init": [
                    "#!/bin/bash",
                    "apt-get install -y awscli jq"
                  ]
                }
              }, {
                "Plugin": "flavor-swarm/worker",
                "Properties": {
                  "InitScriptTemplateURL": "{{ var "/script/baseurl" }}/worker-init.vagrant.tpl",
                  "SwarmJoinIP": "{{ var "/m1/ip" }}",
                  "Docker" : {
                    {{ if var "/certificate/ca/service" }}"Host" : "tcp://{{ var "/m1/ip" }}:{{ var "/docker/remoteapi/tlsport" }}",
                    "TLS" : {
                      "CAFile": "{{ var "/docker/remoteapi/cafile" }}",
                      "CertFile": "{{ var "/docker/remoteapi/certfile" }}",
                      "KeyFile": "{{ var "/docker/remoteapi/keyfile" }}",
                      "InsecureSkipVerify": false
                    }
                    {{ else }}"Host" : "tcp://{{ var "/m1/ip" }}:{{ var "/docker/remoteapi/port" }}"
                    {{ end }}
                  }
                }
              }
            ]
          }
        }
      }
    }
  }
]
