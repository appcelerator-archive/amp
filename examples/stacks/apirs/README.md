APIRS
=====

This document outlines the steps to deploy the APIRS stack on AMP.  

## Overview

The deployment consists of 3 stages:

* Stage 1 starts: 
  * mongodb
  * consul
  * redis
  * registry

* Stage 2 initializes:
  * the replica set in mongodb
  * KV store in consul

* Stage 3 starts:
  * acs
  * dashboard
  * admin
  * stratus
  * app-stats-monitor
  * push-dispatcher
  * registry-auth

> NOTE: *The mongodb configuration and the consul configuration have to be created as a Docker config beforehand.*

## Prerequisite

> NOTE: Make sure you follow the [prerequisites](docs/README#prerequisites) specific to your OS before getting started.

Login to Portus with your Portus credentials using Docker login:

    $ docker login services-registry.cloudapp-enterprise-preprod.appctest.com:5000

Build the `kvcodec` tool in order to prepare your consul KV configuration and run in this directory:

    $ go build -o kvcodec .
    
Make sure your AMP CLI is at least version [0.17.0](https://github.com/appcelerator/amp/releases/tag/v0.17.0).

If you plan to deploy on AWS, prepare your AWS credentials. The recommended way is to use the [~/.aws/credentials file](http://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html) but you can also pass the access key id and the secret access key as parameter to the CLI, as described in `amp cluster create -h`.

## Local Deployment

For local deployments you won't have S3 buckets or NFS servers, unless you set it by your own means.

### Cluster creation & sign up
    amp cluster create
    amp -k user signup

No need to verify the login, because it's a local deployment, you're logged in by default (relaxed security for dev environments).

### Configuration
The `consul.json` file is configured for local deployments. You only need to generate your Consul KV configuration file using the following command:

    ./kvcodec encode consul.json > kv.json

Then you need to create the following Docker `config`s:

    amp -k config create mongo_dbinit mongo_dbinit.sh
    amp -k config create consul_kv kv.json

### Deploying the stacks

Run in this directory:

    amp -k stack deploy --with-registry-auth -c apirs.stage1.yml
    # Check point: amp -k service ps apirs_mongo-primary

    amp -k stack deploy --with-registry-auth -c apirs.stage2.yml
    # Check point: amp -k service ps apirs_mongo-init

    amp -k stack deploy --with-registry-auth -c apirs.stage3.yml

## Cloud Deployment

Full instructions on cluster deployment on AWS is available in the [AWS cluster documentation](docs/awscluster.md), the following section is a quick guide adapted to the APIRS use case.

### Cluster Configuration

    export STACK_NAME=amp-$USER
    export KEYPAIR_NAME=YOUR_KEYPAIR
    export REGION=YOUR_REGION (for instance: us-west-2)
    
#### Create cluster

    amp cluster create --provider aws --aws-region $REGION --aws-stackname $STACK_NAME --aws-parameter KeyName=$KEYPAIR_NAME --aws-parameter NFSEndpoint=true

The output will show the URL of the cluster. We will refer to this as $REMOTE_SERVER.
It will also show the NFS endpoint. We will refer to this as $NFS_ENDPOINT.

    export REMOTE_SERVER="URL of the cluster"
    export NFS_ENDPOINT="URL of the NFS server"
    
If you don't want to store your credentials in the ~/.aws/ folder, you can alternatively use CLI arguments to pass the credentials:

    amp cluster create --provider aws --aws-region $REGION --aws-stackname $STACK_NAME --aws-parameter KeyName=$KEYPAIR_NAME --aws-parameter NFSEndpoint=true --aws-access-key-id YOUR_ACCESS_KEY_ID --aws-secret-access-key YOUR_SECRET_ACCESS_KEY

#### Enable user registration

This is a cloud cluster creation, the user registration will go through email notifications and will require a user verification.
In order to use this feature, the cluster requires a SendGrid key, without which the user registration will fail.
First, obtain a Sendgrid key. You can get a 30 days trial here: https://app.sendgrid.com/signup?id=71713987-9f01-4dea-b3d4-8d0bcd9d53ed

SSH into a manager node and update the amplifier configuration with the following lines:

    cat > amplifier.yml << EOF
    EmailKey: _SENDGRID_KEY_
    SUPassword: _SUPER_USER_PASSWORD_
    JWTSecretKey: _JWT_RANDOM_STRING_
    EOF
    cat amplifier.yml | docker secret create amplifier_yml_custom -
    docker service update --secret-rm amplifier_yml --secret-add source=amplifier_yml_custom,target=/run/secrets/amplifier.yml amp_amplifier

Register your user:

    amp -k -s $REMOTE_SERVER user signup
    # followed by user verification and login (as detailed in the email notification)


#### Adding object storage

You can now add object storage for APIRS. In the following command, we will use `BUCKET1` and `BUCKET2`. You should use your own unique names.

Create the new buckets:

    amp -k -s $REMOTE_SERVER object-store create BUCKET1
    amp -k -s $REMOTE_SERVER object-store create BUCKET2

List the existing buckets available in the cluster, make sure you see BUCKET1 and BUCKET2:

    amp -k -s $REMOTE_SERVER object-store list

### Configuration

The `consul.json` file is configured for local deployments. You need to replace all the domain names occurrences 
in `consul.json` from `apirs.local.appcelerator.io` to the domain you're going to use. For instance:

    sed s/apirs.local.appcelerator.io/apirs.aws.appcelerator.io/g consul.json > cloud.json

You then need to generate a Consul KV configuration file by running he following command.

    ./kvcodec encode cloud.json > kv.json
    
Then you need to create the following Docker `config`s:

    amp -k -s $REMOTE_SERVER config create mongo_dbinit mongo_dbinit.sh
    amp -k -s $REMOTE_SERVER config create consul_kv kv.json

### Deploying the stacks

    amp -k -s $REMOTE_SERVER stack deploy -c apirs.stage1.yml
    # Check point: amp -k -s $REMOTE_SERVER service ps apirs_mongo-primary

    amp -k -s $REMOTE_SERVER stack deploy -c apirs.stage2.yml
    # Check point: amp -k -s $REMOTE_SERVER service ps apirs_mongo-init
    # Have a look at consul, and update the key value store with your domain name. There's also the global/nfs_server and global/nfs_server_ip that should be updated if you have one.

    amp -k -s $REMOTE_SERVER stack deploy --with-registry-auth -c apirs.stage3.yml

    curl -k -H Host:acs.apirs.local.appcelerator.io https://$REMOTE_SERVER/v1/admins/ping.json

This last line simulates a DNS alias by using a HTTP header, but if you set a real DNS record as a CNAME resolving to $REMOTE_SERVER, it can be used directly in curl or a browser.

## Links

* Mongo: https://mongo.apirs.local.appcelerator.io
* Consul: https://consul.apirs.local.appcelerator.io
* ACS: 
    * https://acs.apirs.local.appcelerator.io/v1/admins/ping.json 
    * https://acs.apirs.local.appcelerator.io/v1/admins/pingdb.json
* Dashboard: https://dashboard.apirs.local.appcelerator.io

## APIRS Dashboard

In order to login to **Dashboard**, use the following credentials:

* login: admin@dashboard.local
* password: changeme

## AMP Dashboard

For local deployment, connect to http://dashboard.local.appcelerator.io/

For cloud deployment, check the output of the cluster creation for the link to the dashboards.

An extra dashboard sample with panels specific to APIRS can be imported in Grafana: `apirs.dashboard.json`.
