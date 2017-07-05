## Cluster Management Commands

The `amp cluster` command manages all cluster-related operations for AMP.

### Usage

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
  ls          List deployed amp clusters
  rm          Destroy an amp cluster
  status      Retrieve details about an amp cluster
  update      Update an amp cluster

Run 'amp cluster COMMAND --help' for more information on a command.
```

### Examples

>NOTE: Currently a number of `amp cluster create` options are in a state of transition.
We recommend you stick to the outlined examples for deploying cluster environments.

All the cluster commands use both the `--tag` and `--provider` options, which specify the tag for
the cluster images as well as the target for the cluster you are performing operations on.



#### Creating a cluster on your local system

* To create a cluster locally:
```
$ amp cluster create --provider local
...
Cluster status: healthy
...
```
This will deploy the `AMP` stack on your local docker engine. This AMP deployment only
uses one Manager node to deploy services on.

The target for this cluster will be `localhost:50101`.

##### Secrets

`amp cluster create` uses a docker secret named `amplifier_yml` for amplifier configuration.

If the secret is not present before the invocation of `amp cluster create`, it will be automatically generated with sensible values for the following keys:
- `JWTSecretKey`: A secret key of 128 random characters will be generated.
- `SUPassword`: A super user password of 32 characters will be generated and displayed during the execution of the command.

If the secret is already created, it will be used as is without any modifications.

#### Creating a cluster on AWS

>NOTE: Creating a cluster on AWS is currently in a state of transition

To create a cluster on AWS:
```
$ amp cluster create
```



* To update a cluster with new parameter values:
```
$ amp cluster update
```

* To retrieve the status of the cluster:
```
$ amp cluster status
```

* To remove a cluster:
```
$ amp cluster remove
```
