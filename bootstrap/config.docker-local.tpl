{{ source "default.ikt" }}
{{ source "file:///infrakit/env.ikt" }}
{{ $workerSize := ref "/swarm/size/worker" }}
{{ global "managerSwarmJoinIP" "m1" }}
{{ global "managerDockerHostTLS" (print "tcp://m1:" (ref "/docker/remoteapi/tlsport")) }}
{{ global "managerDockerHost" (print "tcp://m1:" (ref "/docker/remoteapi/port")) }}
{{ global "workerSwarmJoinIP" "m1" }}
{{ global "workerDockerHostTLS" (print "tcp://m1:" (ref "/docker/remoteapi/tlsport")) }}
{{ global "workerDockerHost" (print "tcp://m1:" (ref "/docker/remoteapi/port")) }}
[
  {
    "Plugin": "group",
    "Properties": {
      "ID": "amp-manager-{{ ref "/docker/label/cluster/value" }}",
      "Properties": {
        "Allocation": {
          "LogicalIds": [
            "m1"
          ]
        },
        "Instance": {{ include "include/instance-docker.tpl" }},
        "Flavor": {
          "Plugin": "flavor-combo",
          "Properties": {
            "Flavors": [
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
      "ID": "amp-worker-{{ ref "/docker/label/cluster/value" }}",
      "Properties": {
        "Allocation": {
          "Size": {{ $workerSize }}
        },
        "Instance": {{ include "include/instance-docker.tpl" }},
        "Flavor": {
          "Plugin": "flavor-combo",
          "Properties": {
            "Flavors": [
              {{ include "include/flavor-swarm-worker.tpl" }},
              {{ include "include/flavor-vanilla-verify-ampnet.tpl" }}
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
        "Instance": {{ include "include/instance-docker-amp-proxy.tpl" }},
        "Flavor": {
          "Plugin": "flavor-combo",
          "Properties": {
            "Flavors": [
              {{ include "include/flavor-swarm-worker.tpl" }},
              {{ include "include/flavor-vanilla-verify-ampnet.tpl" }}
            ]
          }
        }
      }
    }
  }
]
