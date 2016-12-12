# Choose how to install

You can install AMP on any cloud platform that runs an operating system (OS) that AMP supports. This includes many flavors and versions of Linux, along with Mac and Windows.

You have two options for installing:

* Manually install on the cloud (create cloud hosts, then install AMP on them)
* Use AMP provisioning capacities to provision cloud infrastructure components

## Use AMP deployment scripts to provision cloud hosts

A set of scripts makes it easy to deploy a full cluster on AWS with AMP running on it.
For each version of AMP (starting with v0.3.0), a related tag is available for these scripts.

create a custom.yaml file with your choice of deployment (VPC id, region, ami id, number of nodes) and run ```ansible-play swarm.yaml```.

Find the detailed information on the [dedicated repo](https://github.com/appcelerator/amp-swarm-deploy)

## Manually install Docker Engine on a cloud host

Alternatively, if you prefer to learn how to deploy AMP step by step, or adapt it to a different cloud provider,

1. Create an account with the cloud provider, and read cloud provider documentation to understand their process for creating hosts.

2. Decide which OS you want to run on the cloud host.

3. Understand AMP prerequisites and install process for the chosen OS. See [Install AMP](../index.md) for a list of supported systems and links to the install guides.

4. Create a host with an AMP supported OS, and install AMP per the instructions for that OS.

[Example (AWS): Manual install on a cloud provider](cloud-ex-aws.md) shows how to create an <a href="https://aws.amazon.com/" target="_blank"> Amazon Web Services (AWS)</a> EC2 instance, and install AMP on it.
