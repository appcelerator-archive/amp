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

data "template_file" "user_data_leader" {
  template = "${file("user-data-leader")}"

  vars {
    tpl_config_base_url = "${var.infrakit_config_base_url}",
    tpl_infrakit_group_suffix = "${random_id.group_suffix.hex}",
    tpl_aws_name = "${var.aws_name}",
    tpl_instance_type = "${var.cluster_instance_type}",
    tpl_key_name = "${var.bootstrap_key_name}",
    tpl_subnet_id = "${aws_subnet.default.id}",
    tpl_iam_instance_profile = "${aws_iam_instance_profile.cluster_instance_profile.id}",
    tpl_security_group_id = "${aws_security_group.default.id}",
  }
}

data "template_file" "user_data_manager" {
  template = "${file("user-data")}"
  vars {
    tpl_config_base_url = "${var.infrakit_config_base_url}",
    tpl_manager_ip = "${aws_instance.m1.private_ip}",
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
    from_port   = 50101
    to_port     = 50199
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 50101
    to_port     = 50199
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 2375
    to_port     = 2375
    protocol    = "tcp"
    cidr_blocks = ["${var.cidr_remote_api}"]
  }
  ingress {
    from_port   = 5000
    to_port     = 5000
    protocol    = "tcp"
    cidr_blocks = ["${var.cidr_remote_api}"]
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

resource "random_id" "group_suffix" {
  byte_length = 8
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
  tags {
    Name = "${var.aws_name}-manager1"
    SwarmRole = "manager"
    Project = "${var.aws_name}"
  }
  user_data = "${data.template_file.user_data_leader.rendered}"
}
resource "aws_instance" "m2" {
  depends_on = [ "aws_subnet.default" ]
  vpc_security_group_ids = [ "${aws_security_group.default.id}" ]
  subnet_id = "${aws_subnet.default.id}"
  availability_zone = "${aws_subnet.default.availability_zone}"
  iam_instance_profile = "${aws_iam_instance_profile.provisioner_instance_profile.id}"
  ami = "${lookup(var.aws_amis, var.aws_region)}"
  key_name = "${var.bootstrap_key_name}"
  instance_type = "${var.bootstrap_instance_type}"
  tags {
    Name = "${var.aws_name}-manager2"
    SwarmRole = "manager"
    Project = "${var.aws_name}"
  }
  user_data = "${data.template_file.user_data_manager.rendered}"
}
resource "aws_instance" "m3" {
  depends_on = [ "aws_subnet.default" ]
  vpc_security_group_ids = [ "${aws_security_group.default.id}" ]
  subnet_id = "${aws_subnet.default.id}"
  availability_zone = "${aws_subnet.default.availability_zone}"
  iam_instance_profile = "${aws_iam_instance_profile.provisioner_instance_profile.id}"
  ami = "${lookup(var.aws_amis, var.aws_region)}"
  key_name = "${var.bootstrap_key_name}"
  instance_type = "${var.bootstrap_instance_type}"
  tags {
    Name = "${var.aws_name}-manager3"
    SwarmRole = "manager"
    Project = "${var.aws_name}"
  }
  user_data = "${data.template_file.user_data_manager.rendered}"
}

# Outputs

output "leader_ip" {
  value = "${aws_instance.m1.public_ip}"
}
output "manager_ips" {
  value = ["${aws_instance.m1.public_ip}","${aws_instance.m2.public_ip}","${aws_instance.m3.public_ip}"]
}
output "cluster_id" {
    value = "${random_id.group_suffix.hex}"
}
