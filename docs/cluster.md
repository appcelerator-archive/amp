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
  rm          Destroy an amp cluster
  status      Retrieve details about an amp cluster

Run 'amp cluster COMMAND --help' for more information on a command.
```

### Examples

>NOTE: Currently a number of `amp cluster create` options are in a state of transition.
We recommend you stick to the outlined examples for deploying cluster environments.

All the cluster commands use both the `--tag` and `--provider` options, which specify the tag for
the cluster images as well as the target for the cluster you are performing operations on.

#### Creating a cluster on your local system

This is the default mode.

* To create a cluster locally:
```
$ amp cluster create --local-no-logs --local-no-metrics
...
2017/08/04 01:17:59 ampctl (version: 0.14.0-dev, build: 08772ef3)
...
{"SwarmID":"sdo9fm7ner6htnu56ww2plo0k","NodeID":"zxij8nozr9xay175jl78xdafu"}

```
This will deploy the `AMP` stack on your local docker engine.
Using the `--tag` option will allow you deploy a cluster with a specific image tag.
Otherwise, the image tag will default to be synchronized with the version of the CLI you are currently using.

The target for this cluster will be `localhost:50101`.

##### Secrets

`amp cluster create` uses a docker secret named `amplifier_yml` for amplifier configuration.

If the secret is not present before the invocation of `amp cluster create`, it will be automatically generated with sensible values for the following keys:
- `JWTSecretKey`: A secret key of 128 random characters will be generated.
- `SUPassword`: A super user password of 32 characters will be generated and displayed during the execution of the command.

If the secret is already created, it will be used as is without any modifications.

#### Creating a cluster on AWS

To target AWS, you should use the --provider aws option, refer to the help for details on aws options: `amp cluster create -h`.

More details on the AWS creation on the [AMP wiki](https://github.com/appcelerator/amp/wiki/AMP-Clusters-deployment-on-AWS)
