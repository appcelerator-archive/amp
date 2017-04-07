{{ source "default.ikt" }}
{{ source "file:///infrakit/env.ikt" }}
{{ $workerSize := ref "/swarm/size/worker" }}
{{ global "managerSwarmJoinIP" (ref "/m1/ip") }}
{{ global "managerDockerHostTLS" (print "tcp://" (ref "/m1/ip") ":" (ref "/docker/remoteapi/tlsport")) }}
{{ global "managerDockerHost" (print "tcp://" (ref "/m1/ip") ":" (ref "/docker/remoteapi/port")) }}
{{ global "workerSwarmJoinIP" (ref "/m1/ip") }}
{{ global "workerDockerHostTLS" (print "tcp://" (ref "/m1/ip") ":" (ref "/docker/remoteapi/tlsport")) }}
{{ global "workerDockerHost" (print "tcp://" (ref "/m1/ip") ":" (ref "/docker/remoteapi/tlsport")) }}
[
  {
    "Plugin": "group",
    "Properties": {
      "ID": "amp-manager-{{ ref "/aws/vpcid" }}",
      "Properties": {
        "Allocation": {
          "LogicalIds": [
            "{{ ref "/m1/ip" }}",
            "{{ ref "/m2/ip" }}",
            "{{ ref "/m3/ip" }}"
          ]
        },
        "Instance": {{ include "include/instance-aws.tpl" }},
        "Flavor": {
          "Plugin": "flavor-combo",
          "Properties": {
            "Flavors": [
              {{ include "include/flavor-vanilla-aws-cli.tpl" }},
              {{ include "include/flavor-swarm-manager.tpl" }},
              {{ include "include/flavor-vanilla-create-ampnet.tpl" }}
            ]
          }
        }
      }
    }
  },
  {
    "Plugin": "group",
    "Properties": {
      "ID": "amp-worker-{{ ref "/aws/vpcid" }}",
      "Properties": {
        "Allocation": {
          "Size": {{ $workerSize }}
        },
        "Instance": {{ include "include/instance-aws.tpl" }},
        "Flavor": {
          "Plugin": "flavor-combo",
          "Properties": {
            "Flavors": [
              {{ include "include/flavor-vanilla-aws-cli.tpl" }},
              {{ include "include/flavor-swarm-worker.tpl" }}
            ]
          }
        }
      }
    }
  }
]
