{{ source "default.ikt" }}
{{ source "file:///infrakit/env.ikt" }}
{{ $workerSize := var "/swarm/size/worker" }}
[
  {
    "Plugin": "group",
    "Properties": {
      "ID": "amp-manager-{{ var "/aws/vpcid" }}",
      "Properties": {
        "Allocation": {
          "LogicalIds": [
            "{{ var "/m1/ip" }}",
            "{{ var "/m2/ip" }}",
            "{{ var "/m3/ip" }}"
          ]
        },
        "Instance": {
          "Plugin": "instance-aws",
          "Properties": {
            "RunInstancesInput": {
              "ImageId": "{{ var "/aws/amiid" }}",
              "InstanceType": "{{ var "/aws/instancetype" }}",
              "KeyName": "{{ var "/aws/keyname" }}",
              "SubnetId": "{{ var "/aws/subnetid" }}",
              {{ if var "/aws/instanceprofile" }}"IamInstanceProfile": {
                "Name": "{{ var "/aws/instanceprofile" }}"
              },{{ end }}
              "SecurityGroupIds": [ "{{ var "/aws/securitygroupid" }}" ]
            },
            "Tags": {
              "Name": "{{ var "/aws/stackname" }}-manager",
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
                  "InitScriptTemplateURL": "{{ var "/script/baseurl" }}/manager-init.tpl",
                  "SwarmJoinIP": "{{ var "/m1/ip" }}",
                  "Docker" : {
                    "Host" : "tcp://{{ var "/m1/ip" }}:{{ var "/docker/remoteapi/port" }}"
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
      "ID": "amp-worker-{{ var "/aws/vpcid" }}",
      "Properties": {
        "Allocation": {
          "Size": {{ $workerSize }}
        },
        "Instance": {
          "Plugin": "instance-aws",
          "Properties": {
            "RunInstancesInput": {
              "ImageId": "{{ var "/aws/amiid" }}",
              "InstanceType": "{{ var "/aws/instancetype" }}",
              "KeyName": "{{ var "/aws/keyname" }}",
              "SubnetId": "{{ var "/aws/subnetid" }}",
              {{ if var "/aws/instanceprofile" }}"IamInstanceProfile": {
                "Name": "{{ var "/aws/instanceprofile" }}"
              },{{ end }}
              "SecurityGroupIds": [ "{{ var "/aws/securitygroupid" }}" ]
            },
            "Tags": {
              "Name": "{{ var "/aws/stackname" }}-worker",
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
                  "InitScriptTemplateURL": "{{ var "/script/baseurl" }}/worker-init.tpl",
                  "SwarmJoinIP": "{{ var "/m1/ip" }}",
                  "Docker" : {
                    "Host" : "tcp://{{ var "/m1/ip" }}:{{ var "/docker/remoteapi/port" }}"
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
