# Default variables used by the *.tf files
# This file shouldn't be modified, override the variables in a terraform.tfvars file instead

variable "aws_region" {
  description = "AWS region to launch servers"
  default     = "us-west-2"
}

variable "aws_profile" {
  description = "AWS credentials profile"
  default     = "default"
}

variable "aws_amis" {
  default = {
    eu-central-1 = "ami-a85480c7"
    eu-west-1 = "ami-971238f1"
    us-east-1 = "ami-2757f631"
    us-west-1 = "ami-44613824"
    us-west-2 = "ami-7ac6491a"
  }
}

variable "vpc_cidrs" {
  default {
    vpc = "192.168.0.0/16"
    subnet1 = "192.168.2.0/24"
  }
}

variable "aws_name" {
  default = "ikt"
}

variable "bootstrap_instance_type" {
  type = "string"
  description = "EC2 HVM instance type (t2.micro, m3.medium, etc)"
  default = "t2.micro"
}

variable "bootstrap_key_name" {
  type = "string"
  description = "Name of an existing EC2 KeyPair to enable SSH access to the instances"
  default = ""
}

variable "cluster_instance_type" {
  type = "string"
  description = "EC2 HVM instance type (t2.micro, m3.medium, etc)"
  default = "t2.small"
}

variable "cluster_key_name" {
  type = "string"
  description = "Name of an existing EC2 KeyPair to enable SSH access to the instances"
  default = ""
}

variable "cluster_subnet_id" {
  type = "string"
  description = "Subnet ID for the cluster instances"
  default = ""
}

variable "cluster_iam_instance_profile" {
  type = "string"
  description = "IAM instance profile ID for the cluster instances"
  default = ""
}

variable "cluster_security_group_id" {
  type = "string"
  description = "Security group ID for the cluster instances"
  default = ""
}

variable "infrakit_config_base_url" {
  type = "string"
  description = "Base URL for InfraKit configuration. there should be a bootstrap.sh, a variables.ikt and a config.tpl file"
  default = "https://raw.githubusercontent.com/appcelerator/amp/master/bootstrap"
}

variable "aufs_volume_size" {
  description = "Size in GB of the EBS volume for the Docker AUFS storage on each node"
  default = 26
}

variable "cidr_remote_api" {
  description = "CIDR to open for Docker remote API access (port 2375)"
  default = "0.0.0.0/0"
}
