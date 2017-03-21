# Provider


provider "aws" {
  region = "${var.aws_region}"
}


# Data


data "aws_iam_policy_document" "provisioner_role_doc" {
  # version = "2012-10-17"
  statement {
    actions = [ "sts:AssumeRole" ]
    effect = "Allow"
    principals {
      type = "Service"
      identifiers = [ "ec2.amazonaws.com" ]
    }
  }
}

data "aws_iam_policy_document" "provisioner_role_policy_doc" {
  # version = "2012-10-17"
  statement {
    actions = [
      "ec2:Describe*",
      "ec2:Get*",
      "ec2:CreateTags",
    ]
    resources = [ "*" ]
    effect = "Allow"
  }
  statement {
    actions = [
      "ec2:RunInstances",
      "ec2:StartInstances",
      "ec2:StopInstances",
      "ec2:RebootInstances",
      "ec2:TerminateInstances",
      "ec2:AttachVolume",
      "ec2:DetachVolume",
    ]
    resources = [ "*" ]
    effect = "Allow"
  }
  statement {
    actions = [
      "ec2:RunInstances",
      "ec2:StartInstances",
      "ec2:StopInstances",
      "ec2:RebootInstances",
      "ec2:TerminateInstances",
      "ec2:AttachVolume",
      "ec2:DetachVolume",
    ]
    resources = [ "*" ]
    effect = "Allow"
  }
  statement {
    actions = [
      "ec2:RunInstances",
    ]
    resources = [ "*" ]
    effect = "Allow"
  }
  statement {
    actions = [
      "iam:PassRole",
    ]
    resources = [ "*" ]
    effect = "Allow"
  }
}

data "aws_iam_policy_document" "cluster_role_doc" {
  # version = "2012-10-17"
  statement {
    actions = [ "sts:AssumeRole" ]
    effect = "Allow"
    principals {
      type = "Service"
      identifiers = [ "ec2.amazonaws.com" ]
    }
  }
}

data "aws_iam_policy_document" "cluster_role_policy_doc" {
  # version = "2012-10-17"
  statement {
    actions = [
      "ec2:DescribeVolume*",
      "ec2:AttachVolume",
      "ec2:CreateVolume",
      "ec2:CreateTags",
      "ec2:ModifyInstanceAttribute",
    ]
    resources = [ "*" ]
    effect = "Allow"
  }
}

data "template_file" "user_data" {
  template = "${file("user-data.sh")}"

  vars {
    region = "${var.aws_region}"
    name = "${var.aws_name}"
    vpc_id = "${aws_vpc.default.id}"
    subnet_id = "${aws_subnet.default.id}"
    security_group_id = "${aws_security_group.default.id}"
    ami = "${lookup(var.aws_amis, var.aws_region)}"
    instance_type = "${var.bootstrap_instance_type}"
    cluster_instance_profile = "${aws_iam_instance_profile.cluster_instance_profile.id}"
    key_name = "${var.bootstrap_key_name}"
    infrakit_config_base_url = "${var.infrakit_config_base_url}"
    aufs_volume_size = "${var.aufs_volume_size}"
  }
}


# Resources


resource "aws_vpc" "default" {
  cidr_block = "${lookup(var.vpc_cidrs, "vpc")}"
  enable_dns_hostnames = true
  tags {
     Name = "${var.aws_name}-vpc"
  }
}

resource "aws_internet_gateway" "default" {
  depends_on = ["aws_vpc.default"]
  vpc_id  = "${aws_vpc.default.id}"
  tags {
     Name = "${var.aws_name}-igw"
  }
}

resource "aws_route_table" "default" {
  depends_on = ["aws_vpc.default"]
  vpc_id = "${aws_vpc.default.id}"
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.default.id}"
  }
  tags {
    Name = "${var.aws_name}-route-table"
  }
}

resource "aws_subnet" "default" {
  depends_on = ["aws_vpc.default"]
  vpc_id = "${aws_vpc.default.id}"
  cidr_block = "${lookup(var.vpc_cidrs, "subnet1")}"
  map_public_ip_on_launch = true
  tags {
     Name = "${var.aws_name}-subnet"
  }
}

resource "aws_route_table_association" "default" {
  depends_on = ["aws_route_table.default", "aws_subnet.default"]
  route_table_id = "${aws_route_table.default.id}"
  subnet_id = "${aws_subnet.default.id}"
}

resource "aws_security_group" "default" {
  depends_on = ["aws_internet_gateway.default"]
  name = "${var.aws_name}-security-group"
  description = "VPC-wide security group"
  tags {
    Name = "${var.aws_name}-security-group"
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["${lookup(var.vpc_cidrs, "vpc")}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # AMP ports need to be added

  vpc_id = "${aws_vpc.default.id}"
}

resource "aws_iam_role" "provisioner_role" {
  name = "${var.aws_name}-provisioner-role"
  path = "/"
  assume_role_policy = "${data.aws_iam_policy_document.provisioner_role_doc.json}"
}

resource "aws_iam_policy" "provisioner_policy" {
  name = "${var.aws_name}-provisioner-policy"
  path = "/"
  policy = "${data.aws_iam_policy_document.provisioner_role_policy_doc.json}"
}

resource "aws_iam_role_policy_attachment" "provisioner_attachment" {
  role = "${aws_iam_role.provisioner_role.name}"
  policy_arn = "${aws_iam_policy.provisioner_policy.arn}"
}

resource "aws_iam_instance_profile" "provisioner_instance_profile" {
  name = "${var.aws_name}-provisioner-instance-profile"
  path = "/"
  roles = [ "${aws_iam_role.provisioner_role.id}" ]
}

resource "aws_iam_role" "cluster_role" {
  name = "${var.aws_name}-cluster_role"
  path = "/"
  assume_role_policy = "${data.aws_iam_policy_document.cluster_role_doc.json}"
}

resource "aws_iam_policy" "cluster_policy" {
  name = "${var.aws_name}-cluster-policy"
  path = "/"
  policy = "${data.aws_iam_policy_document.cluster_role_policy_doc.json}"
}

resource "aws_iam_role_policy_attachment" "cluster_attachment" {
  role = "${aws_iam_role.cluster_role.name}"
  policy_arn = "${aws_iam_policy.cluster_policy.arn}"
}

resource "aws_iam_instance_profile" "cluster_instance_profile" {
  name = "${var.aws_name}-cluster-instance-profile"
  path = "/"
  roles = [ "${aws_iam_role.cluster_role.id}" ]
}

resource "aws_instance" "m1" {
  depends_on = [ "aws_subnet.default" ]
  vpc_security_group_ids = [ "${aws_security_group.default.id}" ]
  subnet_id = "${aws_subnet.default.id}"
  availability_zone = "${aws_subnet.default.availability_zone}"
  iam_instance_profile = "${aws_iam_instance_profile.provisioner_instance_profile.id}"
  ami = "${lookup(var.aws_amis, var.aws_region)}"
  key_name = "${var.bootstrap_key_name}"
  instance_type = "${var.bootstrap_instance_type}"
  private_ip = "192.168.2.254"
  tags {
    Name = "${var.aws_name}-manager-1"
    infrakit.group = "swarm-managers"
    infrakit.role = "managers"
  }
  user_data = "${data.template_file.user_data.rendered}"
}
