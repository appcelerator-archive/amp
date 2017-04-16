### Team Management

The `amp team` command manages all team related operations for AMP.

```
    $ amp team

    Usage:	amp team [OPTIONS] COMMAND [ARGS...]

    Team management operations

    Options:
      -h, --help            Print usage
      -s, --server string   Specify server (host:port)

    Management Commands:
      member      Team member management operations
      resource    Team resource management operations

    Commands:
      create      Create team
      get         Get team
      ls          List team
      rm          Remove team

    Run 'amp team COMMAND --help' for more information on a command.
```

#### Examples

To be able to perform any team related operations, you must be logged in to AMP using a verified account.

* To create a team in an organization:
```
    $ amp team create
```

* To retrieve the list of teams:
```
    $ amp team ls
```

* To retrieve details of a specific team:
```
    $ amp team get
```

* To remove a team:
```
    $ amp team rm
    [or]
    $ amp team del
```

### Team Member Management Commands

```
    $ amp team member

    Usage:	amp team member [OPTIONS] COMMAND [ARGS...]

    Team member management operations

    Options:
      -h, --help            Print usage
      -s, --server string   Specify server (host:port)

    Commands:
      add         Add member to team
      ls          List members
      rm          Remove member from team

    Run 'amp team member COMMAND --help' for more information on a command.
```

#### Examples

* To add a member to a team:
```
    $ amp team member add
```

* To list members in a team:
```
    $ amp team member ls
```

* To remove a member from a team:
```
    $ amp team member rm
    [or]
    $ amp team member del
```

### Team Resource Management Commands

```
    $ amp team resource

    Usage:	amp team resource [OPTIONS] COMMAND [ARGS...]

    Team resource management operations

    Options:
      -h, --help            Print usage
      -s, --server string   Specify server (host:port)

    Commands:
      add         Add resource to team
      ls          List resources
      perm        Change permission level over a resource
      rm          Remove resource from team

    Run 'amp team resource COMMAND --help' for more information on a command.
```

#### Examples

* To add a resource to a team:
```
    $ amp team resource add
```

* To list resources available to a team:
```
    $ amp team resource ls
```

* To remove a resource from a team:
```
    $ amp team resource rm
    [or]
    $ amp team resource del
```

* To change the permission level over a resource:
```
    $ amp team resource perm RESOURCEID read|write|admin
```
