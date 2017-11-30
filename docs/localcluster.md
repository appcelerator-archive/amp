# Local cluster creation

This is the default mode for cluster creation. 

## Why do you need a local cluster?

During development of your application, you will want to have a quick and easy way of running and tearing down a cluster, rather than relying on a cloud based solution. 

## Prerequisites

### Mac OS

To enable Docker engine metrics, you must add some configuration to the Docker daemon. 

Right click on the whale icon, go to `Preferences` -> `Daemon` and select the `Advanced` tab, and add this:
```
{
    "experimental" : true,
    "metrics-addr" : "0.0.0.0:9323"
}
```
Make sure you apply these changes for Docker to be configured with the new settings. 

<p align="center">
   <img width="492" height="683" src="images/DockerforMacDaemonConfig.png" alt="Docker for Mac Daemon settings">
 </p>
  
### Windows

To enable Docker engine metrics, you must add some configuration to the Docker daemon. 

You can configure settings by right clicking the whale icon in the Notifications area and clicking `Settings` -> `Daemon` -> `Advanced` tab.
```
{
    "experimental" : true,
    "metrics-addr" : "0.0.0.0:9323"
}
```
Make sure you apply these changes for Docker to be configured with the new settings. 

![Docker for Windows Daemon settings](images/DockerForWindowsDaemonConfig.png "Docker for Windows Daemon settings")

### Linux

To enable Docker engine metrics, you must add some configuration to the Docker daemon.  

Edit or create the `daemon.json` file in located `etc/docker` to include:
```
$ sudo nano /etc/docker/daemon.json
{
    "experimental" : true,
    "metrics-addr" : "0.0.0.0:9323"
}
```

You must perform this additional step to increase virtual memory needed for Elasticsearch.
```
$ sudo sysctl -w vm.max_map_count=262144
```

To make this change permanent, you can run the following and reboot:
```
$ echo "vm.max_map_count = 262144" | sudo tee -a /etc/sysctl.conf
```

## Cluster Deployment

To create a cluster locally:
```
$ amp cluster create
...
2017/08/04 01:17:59 ampctl (version: 0.18.0-dev, build: a51daf88)
...
{"Swarm Status":"active","Core Services":17,"User Services":0}

```
This will create a single node swarm cluster on your machine and deploy AMP services on top of it.

Once you have started a local cluster, you will be able to deploy stacks and monitor the associated services by signing up.
```
$ amp user signup --name user1 --email user1@amp.com --password [password]
Verification is not necessary for this cluster.
Hi user1! You have been automatically logged in.
```

> NOTE: See the [user](user.md) documentation for additional details about the user account related operations.

With the local cluster, you do not need to verify the account created and you will be logged in automatically after an account creation. 

While deploying a local cluster, the default certificate is self-signed. You need to use the `-k` or `--insecure` option when using any of the AMP commands. 

> TIP: Set an alias for `amp` as `alias amp='amp -k'`. 
  
## What's next?

You can now deploy a stackfile on your newly created local cluster. Please follow the instructions listed [here](stackdeploy.md).

## Cluster creation options

AMP comprises of 4 features: 

* core (mandatory)
* metrics (optional) 
* logs (optional)
* proxy (optional)

It is possible to disable the optional features using the following commands:

To create a local cluster without metrics: 
```
$ amp cluster create --local-no-metrics
``` 

To create a local cluster without logging:
```
$ amp cluster create --local-no-logs
```

To create a local cluster without proxy:
```
$ amp cluster create --local-no-proxy
```

> TIP: Use `-h` or `--help` option for any of the AMP commands or sub-commands to more information about the command's usage.
