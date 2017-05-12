{{ source "default.ikt" }}
{{ source "file:///infrakit/env.ikt" }}
{{ $workerSize := var "/swarm/size/worker" }}
[
  {
    "Plugin": "group",
    "Properties": {
      "ID": "amp-worker-{{ var "/aws/stackname" }}",
      "Properties": {
        "Allocation": {
          "Size": {{ $workerSize }}
        },
        "Instance": {
          "Plugin": "instance-aws/ec2-instance",
          "Properties": {
            "RunInstancesInput": {
              "ImageId": "{{ var "/aws/amiid" }}",
              "InstanceType": "{{ var "/aws/instancetype" }}",
              "KeyName": "{{ var "/aws/keyname" }}",
              "SubnetId": "{{ var "/aws/subnetid1" }}",
              {{ if var "/aws/instanceprofile" }}"IamInstanceProfile": {
                "Name": "{{ var "/aws/instanceprofile" }}"
              },{{ end }}
              "SecurityGroupIds": [ "{{ var "/aws/securitygroupid" }}" ]
            },
            "Tags": {
              "Name": "{{ var "/aws/stackname" }}-worker",
              "{{ var "/docker/label/cluster/key" }}": "{{ var "/docker/label/cluster/value" }}",
              "SwarmRole" : "worker",
              "ManagedBy": "InfraKit"
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
                    "sysctl -w vm.max_map_count=262144",
                    "echo 'vm.max_map_count = 262144' > /etc/sysctl.d/99-amp.conf"
                  ]
                }
              }, {
                "Plugin": "flavor-swarm/worker",
                "Properties": {
                  "InitScriptTemplateURL": "{{ var "/script/baseurl" }}/worker-init.tpl",
                  "SwarmJoinIP": "{{ var "/docker/manager/host" }}",
                  "Docker" : {
                    "Host" : "unix:///var/run/docker.sock"
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
