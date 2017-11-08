MBaaS
=====

## Prerequisite
Login to Portus with your Portus credentials using docker login:

    $ docker login services-registry.cloudapp-enterprise-preprod.appctest.com:5000


## Local Deployment

The deployment consists in multiple stages:

* Stage 1: Starts mongodb, consul & redis
* Stage 2: Initializes the replica set in mongodb and the KV store in consul
* Stage 3: Starts ACS & Dashboard

The mongodb configuration and the consul configuration have to be created as Docker config beforehand.

    amp cluster create

    amp -k user signup

    amp -k config create mongo_dbinit mongo_dbinit.sh
    amp -k config create consul_kv consul.local.json

    amp -k stack deploy --with-registry-auth -c mbaas.stage1.yml
    # Check point: amp -k service ps mbaas_mongo-primary

    amp -k stack deploy --with-registry-auth -c mbaas.stage2.yml
    # Check point: amp -k service ps mbaas_mongo-init

    amp -k stack deploy --with-registry-auth -c mbaas.stage3.yml

## Links

* Mongo: https://mongo.mbaas.local.appcelerator.io
* Consul: https://consul.mbaas.local.appcelerator.io
* ACS: https://acs.mbaas.local.appcelerator.io
* Dashboard: https://dashboard.mbaas.local.appcelerator.io

## Dashboard

In order to login to **Dashboard**, use the following credentials:

* login: admin@dashboard.local
* password: changeme
