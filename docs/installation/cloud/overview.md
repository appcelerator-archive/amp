# Choose how to install

You can install AMP on any cloud platform that runs an operating system (OS) that AMP supports. This includes many flavors and versions of Linux, along with Mac and Windows.

You have two options for installing:

* Manually install on the cloud (create cloud hosts, then install AMP on them)
* Use AMP provisioning capacities to provision cloud infrastructure components

## Manually install Docker Engine on a cloud host

To install on a cloud provider:

1. Create an account with the cloud provider, and read cloud provider documentation to understand their process for creating hosts.

2. Decide which OS you want to run on the cloud host.

3. Understand AMP prerequisites and install process for the chosen OS. See [Install AMP](../index.md) for a list of supported systems and links to the install guides.

4. Create a host with an AMP supported OS, and install AMP per the instructions for that OS.

[Example (AWS): Manual install on a cloud provider](cloud-ex-aws.md) shows how to create an <a href="https://aws.amazon.com/" target="_blank"> Amazon Web Services (AWS)</a> EC2 instance, and install AMP on it.


## Use AMP deployment scripts to provision cloud hosts

Alternatively, you can also check out latest initiative to fully deploy all required components through dedicated scripts.
Find all related information [here](https://github.com/appcelerator/amp-swarm-deploy)
