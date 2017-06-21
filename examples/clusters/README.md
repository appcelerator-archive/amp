# Cloudformation templates for Docker Swarm cluster creation on AWS

## How to use

Use the aws cli or the console to deploy a swarm cluster.
It is compatible with regions us-east-1, us-east-2, us-west-2, eu-west-1 and ap-southeast-2.

## Content

The template will create the infrastructure (1 VPC, 3 subnets for HA on 3 datacenters, security groups, internet gateway, instance profile), 1 autoscaling group for the manager nodes and 1 autoscaling group for the worker nodes.
Each autoscaling group run a userdata that initialize or join the swarm depending on the nature of the group and depending on the status of the swarm.
The engine API of all managers are available from all nodes in the VPC, which allow to set the labels on the nodes.

It is ready for the deployment of AMP, with the help of the CLI on one of the manager nodes:

    amp -s localhost -p local

or the development version:

    amp -s localhost -p local --tag latest

you can add the option --secrets-dir and --domain if you have a domain and a certificate available, please check the CLI help for more information.

## Output

The output of the stack lists the DNS name of the ELB in front of the manager nodes. It can be used for ssh access, https access to swarm services and configuration of the remote server in the CLI (--server option).
