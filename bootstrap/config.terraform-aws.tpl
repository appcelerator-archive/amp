{{ source "default.ikt" }}
{{ source "file:///infrakit/env.ikt" }}
{{ $workerSize := ref "/swarm/size/worker" }}
[
  {
    "Plugin": "group",
    "Properties": {
      "ID": "amp-manager-{{ ref "/aws/stackname" }}",
      "Properties": {
        "Allocation": {
          "Size": 3
        },
        "Instance": {
          "Plugin": "instance-terraform",
          "Properties": {
            "type": "aws_instance",
            "value": {
              "ami": "${lookup(var.aws_amis, var.aws_region)}",
              "instance_type": "${var.bootstrap_instance_type}",
              "key_name": "${var.bootstrap_key_name}",
              "subnet_id": "${aws_subnet.default.id}",
              "iam_instance_profile": "${aws_iam_instance_profile.provisioner_instance_profile.id}",
              "vpc_security_group_ids": [ "${aws_security_group.default.id}" ],
              "tags": {
                "Name": "${var.aws_name}-manager",
                "Deployment": "Infrakit",
                "Role" : "manager"
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
                    "apt-get install -y awscli jq"
                  ]
                }
              }, {
                "Plugin": "flavor-swarm/manager",
                "Properties": {
                  "InitScriptTemplateURL": "{{ ref "/script/baseurl" }}/manager-init.tpl",
                  "SwarmJoinIP": "${aws_instance.m1.private_ip}",
                  "Docker" : {
                    {{ if ref "/certificate/ca/service" }}"Host" : "tcp://${aws_instance.m1.private_ip}:{{ ref "/docker/remoteapi/tlsport" }}",
                    "TLS" : {
                      "CAFile": "{{ ref "/docker/remoteapi/cafile" }}",
                      "CertFile": "{{ ref "/docker/remoteapi/certfile" }}",
                      "KeyFile": "{{ ref "/docker/remoteapi/keyfile" }}",
                      "InsecureSkipVerify": false
                    }
                    {{ else }}"Host" : "tcp://${aws_instance.m1.private_ip}:{{ ref "/docker/remoteapi/port" }}"
                    {{ end }}
                  }
                }
              }, {
                "Plugin": "flavor-vanilla",
                "Properties": {
                  "Init": [
                    "# create an overlay network",
                    "docker network inspect {{ ref "/amp/network" }} 2>&1 | grep -q 'No such network' && \\",
                    "  docker network create -d overlay --attachable {{ ref "/amp/network" }}",
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
      "ID": "amp-worker-{{ ref "/aws/stackname" }}",
      "Properties": {
        "Allocation": {
          "Size": {{ $workerSize }}
        },
        "Instance": {
          "Plugin": "instance-terraform",
          "Properties": {
            "type": "aws_instance",
            "value": {
              "ami": "${lookup(var.aws_amis, var.aws_region)}",
              "instance_type": "${var.bootstrap_instance_type}",
              "key_name": "${var.bootstrap_key_name}",
              "subnet_id": "${aws_subnet.default.id}",
              "iam_instance_profile": "${aws_iam_instance_profile.provisioner_instance_profile.id}",
              "vpc_security_group_ids": [ "${aws_security_group.default.id}" ],
              "tags": {
                "Name": "{var.aws_name}-worker",
                "Deployment": "Infrakit",
                "Role" : "worker"
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
                    "apt-get install -y awscli jq"
                  ]
                }
              }, {
                "Plugin": "flavor-swarm/worker",
                "Properties": {
                  "InitScriptTemplateURL": "{{ ref "/script/baseurl" }}/worker-init.tpl",
                  "SwarmJoinIP": "${aws_instance.m1.private_ip}",
                  "Docker" : {
                    {{ if ref "/certificate/ca/service" }}"Host" : "tcp://${aws_instance.m1.private_ip}:{{ ref "/docker/remoteapi/tlsport" }}",
                    "TLS" : {
                      "CAFile": "{{ ref "/docker/remoteapi/cafile" }}",
                      "CertFile": "{{ ref "/docker/remoteapi/certfile" }}",
                      "KeyFile": "{{ ref "/docker/remoteapi/keyfile" }}",
                      "InsecureSkipVerify": false
                    }
                    {{ else }}"Host" : "tcp://${aws_instance.m1.private_ip}:{{ ref "/docker/remoteapi/port" }}"
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
