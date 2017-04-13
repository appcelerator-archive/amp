
### Logs

The `amp logs` command is used to query or stream logs. It provides useful filtering options to manage what is presented.

```
    $ amp logs --help

    Usage:	amp logs [OPTIONS] SERVICE

    Fetch log entries matching provided criteria

    Options:
          --container string   Filter by the given container
      -f, --follow             Follow log output
      -h, --help               Print usage
      -i, --infra              Include infrastructure logs
      -m, --meta               Display entry metadata
          --msg string         Filter the message content by the given pattern
          --node string        Filter by the given node
      -n, --number int         Number of results (default 1000)
      -s, --server string      Specify server (host:port)
          --stack string       Filter by the given stack
```

### Examples

* To fetch and follow all the logs:
```
  $ amp logs -f
```

* To fetch and follow the logs for a specific service:
```
  $ amp logs -f etcd
```

* To search for a specific pattern through all the logs:
```
  $ amp logs --msg error
```

* To fetch the logs for a service called `foobar`, using partial service name:
```
  $ amp logs foo
```

* To fetch the all the logs, including the infrastructure ones:
```
  $ amp logs -i
```

* To fetch the all the logs and display metadata associated with each entry:
```
  $ amp logs -m
```
