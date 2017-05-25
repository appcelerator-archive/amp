## AMP labels

Here are the labels used in AMP for tagging resources:

Resource | Name | Values | Description
--- | --- | --- | ---
container | `io.amp.role` | [`infrastructure`, `tool` ] | Set on amp containers to describe their role |
service | `io.amp.role` | [`infrastructure`, `tool` ] | Set on amp containers to describe their role |
daemon | `atomiq.clusterid` | Cluster ID | Docker engine label for the local bootstrap cluster |
daemon | `infrakit.group` |  InfraKit group name | Docker engine label for the local bootstrap cluster (is an EC2 tag on AWS cluster deployment) |

Resource type: https://docs.docker.com/engine/userguide/labels-custom-metadata/
