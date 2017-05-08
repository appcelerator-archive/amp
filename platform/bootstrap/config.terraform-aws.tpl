{{ source "default.ikt" }}
{{ source "file:///infrakit/env.ikt" }}
{{ $cattleWorkerSize := var "/swarm/size/worker/cattle" }}
[
  {
    "Plugin": "group",
    "Properties": {
      "ID": "amp-worker-cattle-{{ var "/aws/stackname" }}",
      "Properties": {
        "Allocation": {
          "Size": {{ $cattleWorkerSize }}
        },
        "Instance": {
          "Plugin": "instance-terraform",
          "Properties": {
            "type": "aws_instance",
            "value": {
              "ami": "${lookup(var.aws_amis, var.aws_region)}",
              "instance_type": "${var.cluster_instance_type}",
              "key_name": "${var.cluster_key_name}",
              "subnet_id": "${var.cluster_subnet_id}",
              "iam_instance_profile": "${var.cluster_iam_instance_profile}",
              "vpc_security_group_ids": [ "${var.cluster_security_group_id}" ],
              "tags": {
                "SwarmRole" : "worker",
                "Project": "${var.aws_name}"
              }
            }
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
                    "apt-get install -y awscli jq",
                    "sysctl -w vm.max_map_count=262144"
                  ]
                }
              }, {
                "Plugin": "flavor-swarm/worker",
                "Properties": {
                  "InitScriptTemplateURL": "{{ var "/script/baseurl" }}/worker-init.tpl",
                  "SwarmJoinIP": "{{ var "/bootstrap/ip" }}",
                  "Docker" : {
                    {{ if var "/certificate/ca/service" }}"Host" : "unix:///var/run/docker.sock",
                    "TLS" : {
                      "CAFile": "{{ var "/docker/remoteapi/cafile" }}",
                      "CertFile": "{{ var "/docker/remoteapi/certfile" }}",
                      "KeyFile": "{{ var "/docker/remoteapi/keyfile" }}",
                      "InsecureSkipVerify": false
                    }
                    {{ else }}"Host" : "unix:///var/run/docker.sock"
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
