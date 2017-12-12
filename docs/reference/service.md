## Service Management Commands

The `amp service` command is used to manage all service related operations for AMP.

### Usage

```
$ amp service --help

Usage:	amp service [OPTIONS] COMMAND

Service management operations

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
  -s, --server string   Specify server (host:port)

Commands:
  inspect     Display detailed information of a service
  logs        Display logs of given service matching provided criteria
  ls          List services
  ps          List tasks of a service
  scale       Scale a replicated service

Run 'amp service COMMAND --help' for more information on a command.
```

> TIP: Use `-h` or `--help` option for any of the AMP commands or sub-commands to more information about the command's usage.

> NOTE: To be able to perform any service related operations, you must be logged in to AMP using a verified account.

### Examples

* To inspect a service:
```
$ amp service inspect pinger
{
        "ID": "0k1fzgc2pfuyowej0gedm6kmr",
        "Version": {
                "Index": 231
        },
        "CreatedAt": "2017-07-05T22:00:41.444720874Z",
        "UpdatedAt": "2017-07-05T22:00:41.449032363Z",
        "Spec": {
                "Name": "pinger_pinger",
                "Labels": {
                        "com.docker.stack.namespace": "pinger"
                },
                "TaskTemplate": {
                        "ContainerSpec": {
                                "Image": "subfuzion/pinger:latest",
                                "Labels": {
                                        "com.docker.stack.namespace": "pinger"
                                },
                                "Env": [
                                        "SERVICE_PORTS=3000",
                                        "VIRTUAL_HOST=pinger.examples.*,https://pinger.examples.*"
                                ],
                                "StopGracePeriod": 10000000000,
                                "DNSConfig": {}
                        },
...
```
> TIP: Use the `service-id` or `service-name` as the argument.

* To list the available services with detailed status about their tasks:
```
$ amp service ls
ID                          NAME            MODE         REPLICAS   STATUS    IMAGE                 TAG
0k1fzgc2pfuyowej0gedm6kmr   pinger_pinger   replicated   3/3        RUNNING   appcelerator/pinger      latest
```
> NOTE: this command only displays services which are part of user deployed stacks.

* To list the tasks of a service:
```
$ amp service ps pinger_pinger
ID                          NAME              IMAGE                        DESIRED STATE   CURRENT STATE   NODE ID                     ERROR
2oyhxm5eon3u40didkvglc4oa   pinger_pinger.1   appcelerator/pinger:latest   RUNNING         RUNNING         7z8aobvghzadvrb8n3zf22ake
9uy23btldzx9ww901djhaph4c   pinger_pinger.2   appcelerator/pinger:latest   RUNNING         RUNNING         7z8aobvghzadvrb8n3zf22ake
nwm34nsxexmhujixni7b1k6gh   pinger_pinger.3   appcelerator/pinger:latest   RUNNING         RUNNING         7z8aobvghzadvrb8n3zf22ake
```

* To retrieve logs of a service:
```
$ amp service logs pinger
         pinger_pinger.2 | listening on :3000
         pinger_pinger.1 | listening on :3000
         pinger_pinger.3 | listening on :3000
         pinger_pinger.2 | listening on :3000
         pinger_pinger.1 | listening on :3000
         pinger_pinger.3 | listening on :3000
...
```
> TIP: Use the `service-id` or `service-name` as the argument.

* To scale up or scale down the number of replicas of the service:
```
$ amp service scale
service id or name: pinger_pinger
replicas: 2
Service pinger_pinger has been scaled to 2 replicas.
```
