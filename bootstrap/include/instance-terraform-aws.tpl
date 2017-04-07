{
  "Plugin": "instance-terraform",
  "Properties": {
    "type": "aws_instance",
    "value": {
      "ami": "${lookup(var.aws_amis, var.aws_region)}",
      "instance_type": "${var.cluster_instance_type}",
      "key_name": "${var.cluster_key_name}",
      "subnet_id": "${var.cluster_subnet_id}",
      "iam_instance_profile": "${var.cluster_iam_instance_profile}",
      "vpc_security_group_ids": [ "${var.cluster_security_group_id}" ],
      "tags": {
        "Name": "{{ ref "/aws/stackname" }}"
      }
    }
  }
}
