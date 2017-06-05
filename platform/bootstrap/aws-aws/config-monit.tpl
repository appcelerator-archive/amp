{{ source "../default.ikt" }}
{{ source "file:///infrakit/env.ikt" }}
  {
      "ID": "amp-worker-monit-{{ var "/aws/stackname" }}",
      "Properties": {
        "Allocation": {
          "Size": {{ var "/swarm/size/worker/monit" }}
        },
        "Instance": {
          "Plugin": "instance-aws/ec2-instance",
          "Properties": {
            "RunInstancesInput": {
              "ImageId": "{{ var "/aws/amiid" }}",
              "InstanceType": "{{ var "/aws/instancetype/monit" }}",
              "KeyName": "{{ var "/aws/keyname" }}",
              "SubnetId": "{{ var "/aws/subnetid1" }}",
              {{ if var "/aws/instanceprofile" }}"IamInstanceProfile": {
                "Name": "{{ var "/aws/instanceprofile" }}"
              },{{ end }}
              "SecurityGroupIds": [ "{{ var "/aws/securitygroupid" }}" ]
            },
            "Tags": {
              "Name": "{{ var "/aws/stackname" }}-worker-monit",
              "{{ var "/docker/label/cluster/key" }}": "{{ var "/docker/label/cluster/value" }}",
              "SwarmRole" : "worker",
              "WorkerType": "monit",
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
                    "apt-get install -y awscli jq sysstat iotop",
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
                  },
                  "EngineLabels": {{ var "/swarm/labels/worker/monit" | jsonDecode | jsonEncode }}
                }
              }
            ]
          }
        }
      }
  }
