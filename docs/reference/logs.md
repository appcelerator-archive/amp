# Logs Command

The `amp logs` command is used to query or stream logs. It provides useful filtering options to manage what is presented.

## Usage

```
$ amp logs --help

Usage:	amp logs [OPTIONS] SERVICE [flags]

Display logs matching provided criteria

Options:
      --container string   Filter by the given Container
  -f, --follow             Follow log output
  -h, --help               Print usage
  -i, --include            Include AMP logs
  -k, --insecure           Control whether amp verifies the server's certificate chain and host name
  -m, --meta               Display entry metadata
      --msg string         Filter the message content by the given pattern
      --node string        Filter by the given node
  -n, --number int32       Number of results (default 1000)
  -r, --raw                Display raw logs (no prefix)
      --regexp             Treat '--msg' option as a regular expression
  -s, --server string      Specify server (host:port)
      --since int32        Number of days to include in the search (maximum 100) (default 2)
      --stack string       Filter by the given stack
```

> TIP: Use `-h` or `--help` option for any of the AMP commands or sub-commands to more information about the command's usage.

> NOTE: To be able to perform any logs related operations, you must be logged in to AMP using a verified account.

## Examples

* To fetch and follow all the logs:
```
$ amp logs -f
```

* To fetch and follow the logs for a specific service:
```
$ amp logs -f amp_etcd
```

* To search for a specific pattern through all the logs:
```
$ amp logs --msg error
```

* To fetch the logs for a service called `foobar`, using partial service name:
```
$ amp logs foo
```

* To fetch infrastructure logs:
```
$ amp logs -i
```

* To fetch all the logs and display metadata associated with each entry:
```
$ amp logs -m
```
