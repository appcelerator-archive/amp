# Configuration Management Commands

The `amp config` command manages all configuration related operations for AMP.

## Usage

```
$ amp config --help

Usage:	amp config [OPTIONS] COMMAND

Configuration management operations

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
  -s, --server string   Specify server (host:port)

Commands:
  create      Create a config from a file or STDIN as content
  ls          List configs
  remove      Remove one or more configs

Run 'amp config COMMAND --help' for more information on a command.
```

> TIP: Use `-h` or `--help` option for any of the AMP commands or sub-commands to more information about the command's usage.

> NOTE: To be able to perform any secret related operations, you must be logged in to AMP using a verified account.

## Examples

* To create a config:
```
$ amp config create test ./test
66cxrj0wcn2ugqspyocmui2tb
```
On success, this command returns an alphanumeric ID which is the secret ID.

* To list all config:
```
$ amp config ls
ID                          NAME
66cxrj0wcn2ugqspyocmui2tb   test
c74xny5z0nv5mtyb8cl8k799y   prometheus_alerts_rules
```

* To remove one or more configs:
```
$ amp config rm test
test
```
