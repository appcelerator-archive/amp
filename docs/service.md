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
  logs        Display filtered service logs
  ls          List services
  ps          List tasks of a service
  scale       Scale a replicated service

Run 'amp service COMMAND --help' for more information on a command.
```

>NOTE: To be able to perform any service related operations, you must be logged in to AMP using a verified account.

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
>Use the `service-id` or `service-name` as the argument.

* To list the available services with detailed status about their tasks:
```
$ amp service ls
ID                          NAME            MODE         REPLICAS   FAILED TASKS   STATUS    IMAGE                       TAG
0k1fzgc2pfuyowej0gedm6kmr   pinger_pinger   replicated   3/3        0              RUNNING   subfuzion/pinger            latest
```
> NOTE: this command only displays services which are part of user deployed stacks.

* To list the tasks of a service:
```
$ amp service ps pinger_pinger
ID                          IMAGE                     DESIRED STATE   CURRENT STATE   NODE ID                     ERROR
2oyhxm5eon3u40didkvglc4oa   subfuzion/pinger:latest   RUNNING         RUNNING         z8abovo2189upwpgpd063qs4d
9uy23btldzx9ww901djhaph4c   subfuzion/pinger:latest   RUNNING         RUNNING         z8abovo2189upwpgpd063qs4d
nwm34nsxexmhujixni7b1k6gh   subfuzion/pinger:latest   RUNNING         RUNNING         z8abovo2189upwpgpd063qs4d
```

* To retrieve logs of a service:
```
$ amp service logs pinger
2017/07/05 21:15:13 listening on :3000
2017/07/05 21:15:14 listening on :3000
2017/07/05 21:15:15 listening on :3000
2017/07/05 22:00:49 listening on :3000
...
```
>Use the `service-id` or `service-name` as the argument.

* To scale the number of replicas of the service:
```
$ amp service scale
service id or name: pinger_pinger
replicas: 2
Service pinger_pinger has been scaled to 2 replicas.
```
