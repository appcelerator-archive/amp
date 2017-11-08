MBaaS
=====

## Prerequisite

Login to Portus with your Portus credentials using docker login:

    $ docker login services-registry.cloudapp-enterprise-preprod.appctest.com:5000

## Deployment

The deployment consists in multiple stages:

* Stage 1: Starts mongodb, consul & redis
* Stage 2: Initializes the replica set in mongodb and the KV store in consul
* Stage 3: Starts ACS & Dashboard

The mongodb configuration and the consul configuration have to be created as Docker config beforehand.

### Local Deployment

    amp cluster create

    amp -k user signup

    amp -k config create mongo_dbinit mongo_dbinit.sh
    amp -k config create consul_kv consul.local.json

    amp -k stack deploy --with-registry-auth -c mbaas.stage1.yml
    # Check point: amp -k service ps mbaas_mongo-primary

    amp -k stack deploy --with-registry-auth -c mbaas.stage2.yml
    # Check point: amp -k service ps mbaas_mongo-init

    amp -k stack deploy --with-registry-auth -c mbaas.stage3.yml

### Cloud Deployment

    export STACK_NAME=amp-$USER
    export STACK_NAME=YOUR_KEYPAIR
    amp cluster create --provider aws --aws-region us-west-2 --aws-stackname $STACK_NAME --aws-parameter KeyName=$KEYPAIR_NAME --aws-sync --aws-parameter NFSEndpoint=true
    # the output will show the URL of the cluster. It should be used below for the REMOTE_SERVER variable.

    amp -k -s $REMOTE_SERVER user signup
    # followed by user verification and login

    amp -k -s $REMOTE_SERVER config create mongo_dbinit mongo_dbinit.sh
    amp -k -s $REMOTE_SERVER config create consul_kv consul.cloud.json

    amp -k -s $REMOTE_SERVER config stack deploy -c mbaas.stage1.yml
    # Check point: amp -k -s $REMOTE_SERVER service ps mbaas_mongo-primary

    amp -k -s $REMOTE_SERVER config stack deploy -c mbaas.stage2.yml
    # Check point: amp -k -s $REMOTE_SERVER service ps mbaas_mongo-init
    # Have a look at consul, and update the key value store with your domain name. There's also the global/nfs_server and global/nfs_server_ip that should be updated if you have one.

    amp -k -s $REMOTE_SERVER config stack deploy --with-registry-auth -c mbaas.stage3.yml

    curl -k -H Host:acs.mbaas.local.appcelerator.io https://$REMOTE_SERVER/v1/admins/ping.json

## Links

* Mongo: https://mongo.mbaas.local.appcelerator.io
* Consul: https://consul.mbaas.local.appcelerator.io
* ACS: https://acs.mbaas.local.appcelerator.io/v1/admins/ping.json and https://acs.mbaas.local.appcelerator.io/v1/admins/pingdb.json
* Dashboard: https://dashboard.mbaas.local.appcelerator.io

## MBaaS Dashboard

In order to login to **Dashboard**, use the following credentials:

* login: admin@dashboard.local
* password: changeme

## AMP Dashboard

For local deployment, connect to http://dashboard.local.appcelerator.io/

For cloud deployment, check the output of the cluster creation for the link to the dashboards.

An extra dashboard sample with panels specific to MBaaS can be imported in Grafana: `MBaaS_dashboard.json`.
