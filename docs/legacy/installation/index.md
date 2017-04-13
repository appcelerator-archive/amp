# Install AMP

This guide documents how to install AMP and its prerequisites (mainly Docker Swarm).

## Install AMP From binaries

The recommended way to install AMP is to use the [precompiled binaries](./binaries.md).
You'll only need as prerequisite a running Docker 1.12.5+ with Swarm mode enabled.

## Deploy an AMP cluster on a cloud platform

You can install AMP on any cloud platform that runs an operating system (OS) that AMP supports. This includes many flavors and versions of Linux, along with Mac and Windows.

> **Note:** AMP targets deployments on all system types but only linux based installation are fully documented at this point

You have two options for installing:

* Manually install on the cloud (create cloud hosts, then install AMP on them)
* Use AMP provisioning capacities to provision cloud infrastructure components

## Use AMP deployment scripts to provision cloud hosts

In next version of AMP the deployment will be possible without any additional tool.

For now, a set of scripts makes it easier to deploy a full cluster on AWS with all the prerequisites and a fully fonctional AMP.
For each version of AMP (starting with v0.3.0), a corresponding tag is available for these scripts, you can also deploy the development version if you set the configuration accordingly.

The deployment consists in the configuration of a yaml file (custom.yaml) with the characteristics of your deployment (VPC id, region, ami id, number of nodes) and the execution of ```ansible-play swarm.yaml```.

Find the detailed information on the [dedicated github repo](https://github.com/appcelerator/amp-swarm-deploy)

## Manually install Docker Engine on a cloud host

Alternatively, if you prefer to learn how to deploy AMP step by step, or adapt it to a different cloud provider,

1. Create an account with the cloud provider, and read the cloud provider documentation to understand their process for creating hosts.

2. Decide which OS you want to run on the cloud host (any recent Linux distribution will do, deployments on CentOS 7 and Ubuntu 16.04 have been extensively tested).

3. Understand AMP prerequisites and install process for the chosen OS.

4. Create a host with an AMP supported OS, and install AMP per the instructions for that OS.

[Example (AWS): Manual install on a cloud provider](cloud/cloud-ex-aws.md) shows how to create an <a href="https://aws.amazon.com/" target="_blank"> Amazon Web Services (AWS)</a> EC2 instance, and install AMP on it.

## On a local Linux server

* Build AMP on [Ubuntu](./linux/ubuntulinux.md), you can easily adapt to other Linux distributions

## On Mac OS

* Install [Docker for Mac](https://docs.docker.com/docker-for-mac/)
* Initialize Swarm ```docker swarm init```

## On Windows

todo
