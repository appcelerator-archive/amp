# Cloudformation templates for Docker Swarm cluster creation on AWS

## How to use

Use the aws cli or the console to deploy a swarm cluster.
It is compatible with regions us-east-1, us-east-2, us-west-2, eu-west-1 and ap-southeast-2.

Alternatively it can be used through the aws plugin of AMP:

    docker run -it --rm -v ~/.aws:/root/.aws appcelerator/amp-aws:latest init --region us-west-2 --stackname STACKNAME --parameter KeyName=KEYNAME --parameter ConfigurationURL=https://raw.githubusercontent.com/appcelerator/amp/cloudformation-for-amp/examples/clusters --parameter OverlayNetworks=public --parameter DockerChannel=edge --template https://s3.amazonaws.com/io-amp-binaries/templates/latest/aws-swarm-asg.yml

## Content

The template will create the infrastructure (1 VPC, 3 subnets for HA on 3 datacenters, security groups, internet gateway, instance profile), 1 autoscaling group for the manager nodes, 1 autoscaling group for the amp core services worker nodes, 1 autoscaling group for the amp monitoring services worker nodes and 1 autoscaling group for the user services worker nodes.
Each autoscaling group run a userdata that initializes or join the swarm depending on the nature of the group and depending on the status of the swarm.
The engine API of all managers are available from all nodes in the VPC, which allow to set the labels on the nodes.

It is ready for the deployment of AMP, with the help of the CLI on one of the manager nodes:

    amp -s localhost -p local

or the development version:

    amp -s localhost -p local --tag latest

## Parameters

| Parameter | Description | Default Value | Example |
| --------- | ----------- | ------------- | ------- |
| KeyName   | Name of an existing EC2 KeyPair to enable SSH access to the instances | - | YOURNAME-REGION |
| ManagerSize | Number of manager nodes, should be 1, 3 or 5 | 3 | |
| CoreWorkerSize | Number of worker nodes for core services | 3 | |
| UserWorkerSize | Number of worker nodes for user services | 3 | |
| MetricsWorkerSize | Number of worker nodes for metrics services | 1 | |
| LinuxDistribution | AMI OS, Debian or Ubuntu | Ubuntu | Debian |
| ManagerInstanceType | Instance type for the manager nodes. Must be a valid EC2 HVM instance type | t2.small | m4.large |
| CoreInstanceType | Instance type for the core worker nodes. Must be a valid EC2 HVM instance type | m4.large | c4.large |
| UserInstanceType | Instance type for the user worker nodes. Must be a valid EC2 HVM instance type | t2.medium | m4.large |
| MetricsInstanceType | Instance type for the metrics worker nodes. Must be a valid EC2 HVM instance type | t2.large | m4.large |
| DrainManager | Should we drain services from the manager nodes? | false | true |
| ConfigurationURL | URL for manager and worker userdata configuration | https://raw.githubusercontent.com/appcelerator/amp/master/examples/clusters | |
| AufsVolumeSize | Size in GB of the EBS for the /var/lib/docker FS | 26 | 100 |
| OverlayNetworks | name of overlay networks that should be created once swarm is initialized | ampnet | public storage search mq |
| DockerChannel | channel for Docker installation | stable | edge |
| DockerPlugins | space separated list of plugins to install | | rexray/ebs |

## Output

The output of the stack lists the DNS name of the ELB in front of the manager nodes. It can be used for ssh access, https access to swarm services and configuration of the remote server in the CLI (--server option).

| Output | Description | 
| --------- | ----------- |
| DNSTarget | public facing endpoint for the cluster, It can be used for ssh access, https access to swarm services and configuration of the remote server in the CLI |
| MetricsURL | URL for cluster health dashboard |
