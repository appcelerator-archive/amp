# Cluster Management Commands

The `amp cluster` command manages all cluster related operations for AMP.

## Usage

```
$ amp cluster --help

Usage:  amp cluster [OPTIONS] COMMAND

Cluster management operations

Options:
  -h, --help                 Print usage
  -k, --insecure             Control whether amp verifies the server's certificate chain and host name
  -s, --server string        Specify server (host:port)
  -v, --volume stringSlice   Bind mount a volume

Management Commands:
  node        Cluster node management operations

Commands:
  create      Set up a cluster in swarm mode
  rm          Destroy an amp cluster
  status      Retrieve details about an amp cluster

Run 'amp cluster COMMAND --help' for more information on a command.
```

> NOTE: Currently a number of `amp cluster create` options are in a state of transition.
We recommend you stick to the outlined examples for deploying cluster environments.

For learning how to deploy a cluster locally on your machine, see the documentation on [local cluster creation](../localcluster.md)

For learning how to deploy a cluster on AWS, see the documentation on [AWS cluster creation](../awscluster.md)

> TIP: Use `-h` or `--help` option for any of the AMP commands or sub-commands to more information about the command's usage.

## Examples

### Local cluster 

* To create a cluster:
```
$ amp cluster create
```

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

* To know the status of a cluster:
```
$ amp cluster status
```

> NOTE: You must be logged in to the AMP cluster to run this command.

* To remove a cluster:
```
$ amp cluster rm
```

### AWS cluster 

* To create a cluster:

If you don't have the AWS CLI installed, enter the following:
```
$ export REGION=us-west-2
$ export KEY_NAME=user-keypair
$ export ACCESS_KEY_ID=xxxxx
$ export SECRET_ACCESS_KEY=xxxxx
$ export STACK_NAME=amp-test

$ amp cluster create --provider aws --aws-region $REGION --aws-parameter KeyName=$KEY_NAME --aws-access-key-id $ACCESS_KEY_ID --aws-secret-access-key $SECRET_ACCESS_KEY --aws-stackname $STACK_NAME --aws-sync
```

If you have the AWS CLI installed, enter the following:
```
$ export KEY_NAME=user-keypair
$ export STACK_NAME=amp-test

$ amp cluster create --provider aws --aws-parameter KeyName=$KEY_NAME --aws-stackname $STACK_NAME --aws-sync
```

* To remove a cluster:
```
$ amp cluster rm --provider aws --aws-stackname $STACK_NAME --aws-sync
```
