## AMP labels

Here are the labels used in AMP for tagging resources:

Resource | Name | Values | Description
--- | --- | --- | ---
container | `io.amp.role` | [`infrastructure`, `tool` ] | Set on amp containers to describe their role |
service | `io.amp.role` | [`infrastructure`, `tool` ] | Set on amp containers to describe their role |
daemon | `atomiq.clusterid` | Cluster ID | Docker engine label for the local bootstrap cluster |
daemon | `infrakit.group` |  InfraKit group name | Docker engine label for the local bootstrap cluster (is an EC2 tag on AWS cluster deployment) |
node | `amp.type.api` | true | Docker node label for service scheduling (api server) |
node | `amp.type.route` | true | Docker node label for service scheduling (proxy service) |
node | `amp.type.core` | true | Docker node label for service scheduling (core services) |
node | `amp.type.search` | true | Docker node label for service scheduling (log database service) |
node | `amp.type.kv` | true | Docker node label for service scheduling (storage service) |
node | `amp.type.mq` | true | Docker node label for service scheduling (messaging service) |
node | `amp.type.metrics` | true | Docker node label for service scheduling (monitoring service) |
node | `amp.type.user` | true | Docker node label for service scheduling (user services) |

Resource type: https://docs.docker.com/engine/userguide/labels-custom-metadata/
