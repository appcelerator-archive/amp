### Service Management

The `amp service` command is used to manage all service-related operations for AMP.

    Usage:  amp service [OPTIONS] COMMAND

    Service management operations

    Options:
      -h, --help            Print usage
      -s, --server string   Specify server (host:port)

    Commands:
      inspect     Display detailed information of a service
      logs        Get all logs of a given service
      ls          List services
      ps          List tasks of a service

    Run 'amp service COMMAND --help' for more information on a command.

### Examples

To be able to perform any service related operations, you must be logged in to AMP using a verified account.

* To inspect a service:
```
    $ amp service inspect [service]
```
Use the `[service-id]` or `[service-name]` as the argument.

* To list the available services with detailed status about their tasks:
```
    $ amp service ls
```
Note that this command only displays services which are part of user deployed stacks.

* To retrieve logs of a service:
```
    $ amp service logs [service]
```
Use the `[service-id]` or `[service-name]` as the argument.

* To list the tasks of a service:
```
    $ amp service ps [service-id]
```
