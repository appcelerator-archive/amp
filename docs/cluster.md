### Cluster Management

The `amp cluster` command manages all cluster-related operations for AMP.

```
$ amp cluster

Usage:  amp cluster [OPTIONS] COMMAND [ARGS...]

Cluster management operations

Options:
  -h, --help            Print usage
  -s, --server string   Specify server (host:port)

Commands:
  create      Create an amp cluster
  ls          List deployed amp clusters
  rm          Destroy an amp cluster
  status      Retrieve details about an amp cluster
  update      Update an amp cluster

Run 'amp cluster COMMAND --help' for more information on a command.
```

### Bring Your Own Swarm examples

AMP can be deployed on an existing Swarm cluster. It can be a single node or a full cluster.
If the host is not yet a swarm manager, the swarm will be initialized.

* To deploy amp on a manager node:
```
    $ amp -s localhost cluster create
```

* To remove amp from a local manager node:
```
    $ amp -s localhost cluster remove
```

### AWS example

This feature is in preview mode, and may not work in your particular setup. It will be made available on next stable version of AMP.

* To deploy amp on AWS:
```
    $ amp -s localhost cluster create -p aws
```

If no flags are passed to the command, a cluster with default number of worker and manager nodes is created. See `amp cluster create --help` for more options.

* To update a cluster with specific number of workers:
```
    $ amp -s localhost cluster update -p aws -i cluster_id -w 5
```
See `amp cluster update --help` for options.

* To retrieve details about a cluster:
```
    $ amp -s localhost cluster status -p aws -i cluster_id
```

* To remove a cluster:
```
    $ amp -s localhost cluster destroy -p aws -i cluster_id
```

## Secrets

`amp cluster create` uses a docker secret named `amplifier_yml` for amplifier configuration.

If the secret is not present before the invocation of `amp cluster create`, it will be automatically generated with sensible values for the following keys:
- `JWTSecretKey`: A secret key of 128 random characters will be generated.
- `SUPassword`: A super user password of 32 characters will be generated and displayed during the execution of the command.

If the secret is already created, it will be used as is without any modifications.
