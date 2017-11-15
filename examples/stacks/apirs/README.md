APIRS
=====

## Overview

The deployment consists in 3 stages:

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

The mongodb configuration and the consul configuration have to be created as Docker config beforehand.

## Prerequisite

Login to Portus with your Portus credentials using docker login:

    $ docker login services-registry.cloudapp-enterprise-preprod.appctest.com:5000

Build the `kvcodec` tool in order to prepare your consul KV configuration and run in this directory:

    $ go build -o kvcodec .
    
## Local Deployment

### Cluster creation & sign up
    amp cluster create
    amp -k user signup

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

### Cluster creation & sign up

    export STACK_NAME=amp-$USER
    export STACK_NAME=YOUR_KEYPAIR
    
    amp cluster create --provider aws --aws-region us-west-2 --aws-stackname $STACK_NAME --aws-parameter KeyName=$KEYPAIR_NAME --aws-sync --aws-parameter NFSEndpoint=true
    # the output will show the URL of the cluster. It should be used below for the REMOTE_SERVER variable.
    export REMOTE_SERVER="URL of the cluster"
    
    amp -k -s $REMOTE_SERVER user signup
    # followed by user verification and login

### Configuration
The `consul.json` file is configured for local deployments. You need to replace the all the domain names occurrences 
in `consul.json` from `apirs.local.appcelerator.io` to the domain you're going to use. For instance:

    sed s/apirs.local.appcelerator.io/apirs.aws.appcelerator.io/g consul.json > cloud.json

You then need to generate a Consul KV configuration file by running he following command.

    ./kvcodec encode cloud.json > kv.json
    
Then you need to create the following Docker `config`s:

    amp -k -s $REMOTE_SERVER config create mongo_dbinit mongo_dbinit.sh
    amp -k -s $REMOTE_SERVER config create consul_kv kv.json

### Deploying the stacks
    amp -k -s $REMOTE_SERVER config stack deploy -c apirs.stage1.yml
    # Check point: amp -k -s $REMOTE_SERVER service ps apirs_mongo-primary

    amp -k -s $REMOTE_SERVER config stack deploy -c apirs.stage2.yml
    # Check point: amp -k -s $REMOTE_SERVER service ps apirs_mongo-init
    # Have a look at consul, and update the key value store with your domain name. There's also the global/nfs_server and global/nfs_server_ip that should be updated if you have one.

    amp -k -s $REMOTE_SERVER config stack deploy --with-registry-auth -c apirs.stage3.yml

    curl -k -H Host:acs.apirs.local.appcelerator.io https://$REMOTE_SERVER/v1/admins/ping.json

## Links

* Mongo: https://mongo.apirs.local.appcelerator.io
* Consul: https://consul.apirs.local.appcelerator.io
* ACS: 
    * https://acs.apirs.local.appcelerator.io/v1/admins/ping.json 
    * https://acs.apirs.local.appcelerator.io/v1/admins/pingdb.json
* Dashboard: https://dashboard.apirs.local.appcelerator.io

## Apirs Dashboard

In order to login to **Dashboard**, use the following credentials:

* login: admin@dashboard.local
* password: changeme

## AMP Dashboard

For local deployment, connect to http://dashboard.local.appcelerator.io/

For cloud deployment, check the output of the cluster creation for the link to the dashboards.

An extra dashboard sample with panels specific to Apirs can be imported in Grafana: `apirs.dashboard.json`.
