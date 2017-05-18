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

### Examples

* To create a cluster:
```
    $ amp cluster create
```
    If no flags are passed to the command, a cluster with default number of worker and manager nodes is created. See `amp cluster create --help` for more options.

* To update a cluster with specific number of worker and manager nodes:
```
    $ amp cluster update
```
    See `amp cluster update --help` for options.

* To retrieve details about a cluster:
```
    $ amp cluster status
```

* To remove a cluster:
```
    $ amp cluster destroy
```
