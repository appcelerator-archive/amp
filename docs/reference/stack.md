## Stack Management Commands

The `amp stack` command is used to manage all stack related operations for AMP.

### Usage

```
$ amp stack --help

Usage:	amp stack [OPTIONS] COMMAND 

Stack management operations

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
  -s, --server string   Specify server (host:port)

Commands:
  deploy      Deploy a stack with a Docker Compose v3 file
  logs        Display filtered logs for stack
  ls          List deployed stacks
  rm          Remove one or more deployed stacks
  services    List services of a stack

Run 'amp stack COMMAND --help' for more information on a command.
``` 

> TIP: Use `-h` or `--help` option for any of the AMP commands or sub-commands to more information about the command's usage.   

>NOTE: To be able to perform any stack related operations, you must be logged in to AMP using a verified account.

### Examples

* To deploy a stack using a compose file:
```
$ amp stack deploy -c examples/stacks/counter/counter.yml counter
Deploying stack counter using examples/stacks/counter/counter.yml
Creating service counter_go
Creating service counter_redis
```
>NOTE: If no name is specified for the stack, the name of the stack will be the filename.

* To list the deployed stacks, with detailed status about their services:
```
$ amp stack ls
ID                                                                 NAME      RUNNING   COMPLETE   PREPARING   TOTAL   SERVICES   STATUS    OWNER
95508f3ca3ad3877e8c33e69a92a9e3490eb60395bd1b26f0c6f80f1f5521976   counter   2         0          0           2       2/2        RUNNING   su
```
> NOTE: this command only displays stacks created by the user. No infrastructure stacks are displayed.

* To list the services within a stack:
```
$ amp stack services counter
ID            NAME           MODE        REPLICAS  IMAGE
k5fzzzryjpda  counter_redis  replicated  1/1       redis
njsru7ka1gek  counter_go     replicated  3/3       htilford/go-redis-counter
```

* To view the logs of the entire stack:
```
$ amp stack logs counter
...
1:M 05 Jul 21:21:36.050 # WARNING: The TCP backlog setting of 511 cannot be enforced because /proc/sys/net/core/somaxconn is set to the lower value of 128.
1:M 05 Jul 21:21:36.056 # Server started, Redis version 3.2.9
1:M 05 Jul 21:21:36.056 # WARNING you have Transparent Huge Pages (THP) support enabled in your kernel. This will create latency and memory usage issues with Redis. To fix this issue run the command 'echo never > /sys/kernel/mm/transparent_hugepage/enabled' as root, and add it to your /etc/rc.local in order to retain the setting after a reboot. Redis must be restarted after THP is disabled.
...
```

* To update a stack with a new compose file:
```
$ amp stack deploy -c examples/stacks/counter/counter-2.yml counter
Deploying stack counter using examples/stacks/counter/counter-2.yml
Updating service counter_go (id: njsru7ka1gek1xzdt6z4b8wez)
Updating service counter_redis (id: k5fzzzryjpdaanlvqqu5b5qr7)
```

* To remove a stack:
```
$ amp stack rm counter
Removing service counter_redis
Removing service counter_go
```
