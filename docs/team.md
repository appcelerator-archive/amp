## Team Management Commands

The `amp team` command is used to manage all team related operations for AMP.

### Usage

```
$ amp team --help

Usage:  amp team [OPTIONS] COMMAND

Team management operations

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
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

>NOTE: To be able to perform any team related operations, you must be logged in to AMP using a verified account.

>`amp team` commands that require the `--team` option will remember the last used team using a local config file.
If you want to override this, specify the `--team` option in the command.

#### Examples

* To create a team in an organization:
```
$ amp team create team
organization name: org
Team has been created in the organization.
```

* To retrieve details of a specific team:
```
$ amp team get team
organization name: org
Team: team
Organization: org
Created On: 05 Jul 17 15:25
```

* To retrieve the list of teams:
```
$ amp team ls
organization name: org
TEAM   CREATED ON
team   05 Jul 17 15:25
```

* To remove a team:
```
$ amp team rm team
organization name: org
team
```

### Team Member Management Commands

The `amp team member` command is used to manage all team member related operations for AMP.

#### Usage

```
$ amp team member --help

Usage:  amp team member [OPTIONS] COMMAND

Team member management operations

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
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
$ amp team member add johndoe
organization name: org
team name: team
Member(s) have been added to team.
```
>NOTE: The member to be added to the team must be a existing and verified user account,
who is a member of the given organization.

* To list members in a team:
```
$ amp team member ls --org org --team team
organization name: org
team name: team
MEMBER
sample
johndoe
```

* To remove a member from a team:
```
$ amp team member rm johndoe
organization name: org
team name: team
johndoe
```

### Team Resource Management Commands

The `amp team resource` command is used to manage all team resource related operations for AMP.

#### Usage

```
$ amp team resource --help

Usage:  amp team resource [OPTIONS] COMMAND

Team resource management operations

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
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
$ amp team resource add 93fce7d3f...
organization name: org
team name: team
Resource(s) have been added to team.
```
>Resources are stacks. To add the stack to the team, the stack must be deployed in the context of the
organization using the organization account.

* To list resources available to a team:
```
$ amp team resource ls
organization name: org
team name: team
RESOURCE ID                                                        PERMISSION LEVEL
93fce7d3fc8ada786c7db6956849343bcc5700f65d8b1512523561166a2ec455   TEAM_READ
```

* To change the permission level of a resource:
```
$ amp team resource perm 93fce7d3f... write
organization name: org
team name: team
Permission level has been changed.
```
The default permission level of a resource is `read`. The permission level can be changed to `read`, `write` or `admin`.

* To remove a resource from a team:
```
$ amp team resource rm --org org --team team 93fce7d3f...
organization name: org
team name: team
93fce7d3fc8ada786c7db6956849343bcc5700f65d8b1512523561166a2ec455
```
