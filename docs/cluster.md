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

> NOTE: Currently a number of `amp cluster create` options are in a state of transition.
We recommend you stick to the outlined examples for deploying cluster environments.

For learning how to deploy a cluster locally on your machine, see the documentation on [local cluster creation](localcluster.md)

For learning how to deploy a cluster on AWS, see the documentation on [AWS cluster creation](awscluster.md)

> TIP: Use `-h` or `--help` option for any of the AMP commands or sub-commands to more information about the command's usage.
