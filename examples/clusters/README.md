# Cloudformation templates for Docker Swarm cluster creation on AWS

## How to use

Use the aws cli or the console to deploy a swarm cluster.
It is compatible with regions us-east-1, us-east-2, us-west-2, eu-west-1 and ap-southeast-2.

Alternatively it can be used through the aws plugin of AMP:

    docker run -it --rm -v ~/.aws:/root/.aws appcelerator/amp-aws:latest init --region us-west-2 --stackname STACKNAME --parameter KeyName=KEYNAME --template https://s3.amazonaws.com/io-amp-binaries/templates/latest/aws-swarm-asg.yml

## Content

The template will create the infrastructure (1 VPC, 3 subnets for HA on 3 datacenters, security groups, internet gateway, instance profile), 1 autoscaling group for the manager nodes, 1 autoscaling group for the amp core services worker nodes, 1 autoscaling group for the amp monitoring services worker nodes and 1 autoscaling group for the user services worker nodes.
Each autoscaling group run a userdata that initializes or join the swarm depending on the nature of the group and depending on the status of the swarm.
The engine API of all managers are available from all nodes in the VPC, which allow to set the labels on the nodes.

Once the nodes are up and running, it will run the appcelerator/ampadmin image to check the prerequisites, and setup AMP.

## Parameters

| Parameter | Description | Default Value | Example |
| --------- | ----------- | ------------- | ------- |
| KeyName   | Name of an existing EC2 KeyPair to enable SSH access to the instances | - | YOURNAME-REGION |
| ManagerSize | Number of manager nodes, should be 1, 3 or 5 | 3 | |
| CoreWorkerSize | Number of worker nodes for core services | 3 | |
| UserWorkerSize | Number of worker nodes for user services | 3 | |
| MetricsWorkerSize | Number of worker nodes for metrics services | 1 | |
| LinuxDistribution | AMI OS, Debian, Ubuntu or Default | Default | Ubuntu |
| ManagerInstanceType | Instance type for the manager nodes. Must be a valid EC2 HVM instance type | t2.small | m4.large |
| CoreInstanceType | Instance type for the core worker nodes. Must be a valid EC2 HVM instance type | m4.large | c4.large |
| UserInstanceType | Instance type for the user worker nodes. Must be a valid EC2 HVM instance type | t2.medium | m4.large |
| MetricsInstanceType | Instance type for the metrics worker nodes. Must be a valid EC2 HVM instance type | t2.large | m4.large |
| DrainManager | Should we drain services from the manager nodes? | no | yes |
| AufsVolumeSize | Size in GB of the EBS for the /var/lib/docker FS | 26 | 100 |
| OverlayNetworks | name of overlay networks that should be created once swarm is initialized | public core monit | public storage search mq |
| DockerChannel | channel for Docker installation | stable | edge |
| DockerPlugins | space separated list of plugins to install | | rexray/ebs |
| InstallApplication | install AMP | yes | no |
| EnableSystemPrune | Enable Docker system prune | yes | no |

## Output

The output of the stack lists the DNS name of the ELB in front of the manager nodes. It can be used for ssh access, https access to swarm services and configuration of the remote server in the CLI (--server option).

| Output | Description |
| --------- | ----------- |
| VpcId | VPC ID |
| DNSTarget | public facing endpoint for the cluster, It can be used for ssh access, https access to swarm services and configuration of the remote server in the CLI |
| InternalRegistryTarget | internal endpoint for the registry service |
| MetricsURL | URL for cluster health dashboard |

## Custom AMI

the default option for the AMI (Default) is a pre package AMI based on Ubuntu Xenial, with prerequisite packages already installed (in particular Docker).

To build a new version of this image, run the build-ami.sh script. It may take more than 15 min to build the AMI. You'll also need to create a variables.yml file with a content similar to:

```
---
ec2_key_name: "KEY_NAME"
```

Once done, copy the AMI Id in the cloudformation template (aws-swarm-asg.yml).

## Registry

An option of the template is the inclusion of a Docker registry.
It includes a S3 bucket as registry backend, and an autoscaling group of registry containers.
The registry is composed of non swarm nodes and is not part of the swarm.
The registry is only available from the VPC, all Docker swarm nodes are configured with the internal endpoint of the registry as mirror registry.
