{{ source "default.ikt" }}
{{ source "file:///infrakit/env.ikt" }}
{ $workerSize := ref "/swarm/size/worker" }}
[
  {
    "Plugin": "group",
    "Properties": {
      "ID": "amp-manager-{{ ref "/terraform/vpcid" }}",
      "Properties": {
        "Allocation": {
          Size: 2
        },
        "Instance": {
          "Plugin": "terraform",
          "Properties": {
            "RunInstancesInput": {
              "ImageId": "{{ ref "/terraform/amiid" }}",
              "InstanceType": "{{ ref "/terraform/instancetype" }}",
              "KeyName": "{{ ref "/terraform/keyname" }}",
              "SubnetId": "{{ ref "/terraform/subnetid" }}",
              {{ if ref "/terraform/instanceprofile" }}"IamInstanceProfile": {
                "Name": "{{ ref "/terraform/instanceprofile" }}"
              },{{ end }}
              "SecurityGroupIds": [ "{{ ref "/terraform/securitygroupid" }}" ]
            },
            "Tags": {
              "Name": "{{ ref "/terraform/stackname" }}-manager",
              "Deployment": "Infrakit",
              "Role" : "manager"
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
                  "SwarmJoinIP": "{{ ref "/m1/ip" }}",
                  "Docker" : {
                    {{ if ref "/certificate/ca/service" }}"Host" : "tcp://{{ ref "/m1/ip" }}:{{ ref "/docker/remoteapi/tlsport" }}",
                    "TLS" : {
                      "CAFile": "{{ ref "/docker/remoteapi/cafile" }}",
                      "CertFile": "{{ ref "/docker/remoteapi/certfile" }}",
                      "KeyFile": "{{ ref "/docker/remoteapi/keyfile" }}",
                      "InsecureSkipVerify": false
                    }
                    {{ else }}"Host" : "tcp://{{ ref "/m1/ip" }}:{{ ref "/docker/remoteapi/port" }}"
                    {{ end }}
                  }
                }
              }, {
                "Plugin": "flavor-vanilla",
                "Properties": {
                  "Init": [
                    "set -o errexit",
                    "docker network inspect {{ ref "/amp/network" }} 2>&1 | grep -q 'No such network' && \\",
                    "  docker network create -d overlay --attachable {{ ref "/amp/network" }}",
                    "docker service ls {{ ref "/amp/network" }} 2>&1 | grep -q 'No such network' && \\",
                    "docker service create --name amplifier --network {{ ref "/amp/network" }} {{ ref "/amp/amplifier/image" }}:{{ ref "/amp/amplifier/version" }} || true"
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
      "ID": "amp-worker-{{ ref "/terraform/vpcid" }}",
      "Properties": {
        "Allocation": {
          "Size": {{ $workerSize }}
        },
        "Instance": {
          "Plugin": "instance-terraform",
          "Properties": {
            "RunInstancesInput": {
              "ImageId": "{{ ref "/terraform/amiid" }}",
              "InstanceType": "{{ ref "/terraform/instancetype" }}",
              "KeyName": "{{ ref "/terraform/keyname" }}",
              "SubnetId": "{{ ref "/terraform/subnetid" }}",
              {{ if ref "/terraform/instanceprofile" }}"IamInstanceProfile": {
                "Name": "{{ ref "/terraform/instanceprofile" }}"
              },{{ end }}
              "SecurityGroupIds": [ "{{ ref "/terraform/securitygroupid" }}" ]
            },
            "Tags": {
              "Name": "{{ ref "/terraform/stackname" }}-worker",
              "Deployment": "Infrakit",
              "Role" : "worker"
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
                  "SwarmJoinIP": "{{ ref "/m1/ip" }}",
                  "Docker" : {
                    {{ if ref "/certificate/ca/service" }}"Host" : "tcp://{{ ref "/m1/ip" }}:{{ ref "/docker/remoteapi/tlsport" }}",
                    "TLS" : {
                      "CAFile": "{{ ref "/docker/remoteapi/cafile" }}",
                      "CertFile": "{{ ref "/docker/remoteapi/certfile" }}",
                      "KeyFile": "{{ ref "/docker/remoteapi/keyfile" }}",
                      "InsecureSkipVerify": false
                    }
                    {{ else }}"Host" : "tcp://{{ ref "/m1/ip" }}:{{ ref "/docker/remoteapi/port" }}"
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
