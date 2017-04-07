{{ source "default.ikt" }}
{{ source "file:///infrakit/env.ikt" }}
{{ $workerSize := ref "/swarm/size/worker" }}
{{ global "workerSwarmJoinIP" (ref "/bootstrap/ip") }}
{{ global "workerDockerHostTLS" "unix:///var/run/docker.sock" }}
{{ global "workerDockerHost" "unix:///var/run/docker.sock" }}
[
  {
    "Plugin": "group",
    "Properties": {
      "ID": "amp-worker-{{ ref "/aws/stackname" }}",
      "Properties": {
        "Allocation": {
          "Size": {{ $workerSize }}
        },
        "Instance": {{ include "include/instance-terraform-aws.tpl" }},
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
