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
    $ amp team create --org org --team team
```

* To retrieve the list of teams:
```
    $ amp team ls --org org
```

* To retrieve details of a specific team:
```
    $ amp team get --org org --team team
```

* To remove a team:
```
    $ amp team rm --org org team
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
    $ amp team member add --org org --team team --member johndoe
```
Note that the member to be added to the team must be a existing and verified user account.

* To list members in a team:
```
    $ amp team member ls --org org --team team
```

* To remove a member from a team:
```
    $ amp team member rm --org org --team team
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
    $ amp team resource add --org org --team team --res RESOURCEID
```

* To list resources available to a team:
```
    $ amp team resource ls --org org --team team
```

* To remove a resource from a team:
```
    $ amp team resource rm --org org --team team RESOURCEID
```

* To change the permission level of a resource:
```
    $ amp team resource perm RESOURCEID write
```
The default permission level of a resource is `read`. The permission level can be changed to `read`, `write` or `admin`.
