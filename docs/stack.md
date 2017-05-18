### Stacks

The `amp stack` command is used to manage AMP stacks.

    $ amp stack --help

    Usage:	amp stack [OPTIONS] COMMAND

    Stack management operations

    Options:
      -h, --help            Print usage
      -s, --server string   Specify server (host:port)

    Commands:
      deploy      Deploy a stack with a docker compose v3 file
      logs        Get all logs of a given stack
      ls          List deployed stacks
      rm          Remove a deployed stack

    Run 'amp stack COMMAND --help' for more information on a command.

### Examples

To be able to perform any stack related operations, you must be logged in to AMP using a verified account.

* To deploy a stack using a compose file:
```
    $ amp stack deploy -c [path-to-stackfile] [stackname]
```
The Docker name of the stack will be `[stackName]-[id]`, where `[id]` is a unique id given by AMP.

* To update a stack with a new compose file:
```
    $ amp stack deploy -c [path-to-stackfile] [stackname]
```
Use the name of the already deployed stack as the argument.

* To remove a stack:
```
    $ amp stack rm [stackname]
```

* To list the deployed stacks:
```
    $ amp stack ls
```
Note that this command only displays stacks created by the user. No infrastructure stacks are displayed.
