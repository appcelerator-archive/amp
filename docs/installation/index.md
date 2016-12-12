# Install AMP

This guide documents how to install AMP and its prerequisites (mainly Docker Swarm).

You can install AMP on any cloud platform that runs an operating system (OS) that AMP supports. This includes many flavors and versions of Linux, along with Mac and Windows.

> **Note:** AMP targets deployments on all system types but only linux based installation are documented at this point

## Install AMP From binaries

The recommended way to install AMP is to use the [precompiled binaries](./binaries.md).
You'll only need as prerequisite a running Docker 1.12.3+ with Swarm mode enabled.

## Deploy an AMP cluster on a cloud platform

* [Choose how to Install](./cloud/overview.md)

If you prefer instead to build AMP (to test the development branch), follow the steps below.

* [Example: Manual install on a cloud provider](./cloud/cloud-ex-aws.md)

## On a local Linux server

* Build AMP on [Ubuntu](./linux/ubuntulinux.md), you can easily adapt to other Linux distribution

## On Mac OS

* Install [Docker for Mac](https://docs.docker.com/docker-for-mac/)
* Initialize Swarm ```docker swarm init```

## On Windows

todo
