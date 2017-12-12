# Team Management Commands

The `amp team` command is used to manage all team related operations for AMP.

## Usage

```
$ amp team --help

Usage:	amp team [OPTIONS] COMMAND

Team management operations

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
  -s, --server string   Specify server (host:port)

Management Commands:
  member      Manage team members
  resource    Manage team resources

Commands:
  create      Create team
  get         Get team information
  ls          List teams
  rm          Remove one or more teams

Run 'amp team COMMAND --help' for more information on a command.
```

> TIP: Use `-h` or `--help` option for any of the AMP commands or sub-commands to more information about the command's usage.

> NOTE: To be able to perform any team related operations, you must be logged in to AMP using a verified account.

## Examples

> NOTE: `amp team` commands that require the `--team` option will remember the last used team using a local preferences file.
If you want to override this, specify the `--team` option in the command.

* To create a team:
```
$ amp team create team
Team has been created.
```

* To retrieve details of a specific team:
```
$ amp team get team
Team: team
Created On: 05 Dec 17 15:25
```

* To retrieve the list of teams:
```
$ amp team ls
TEAM   CREATED ON
team   05 Dec 17 15:25
```

* To remove a team:
```
$ amp team rm team
team
```

## Team Member Management Commands

The `amp team member` command is used to manage all team member related operations for AMP.

### Usage

```
$ amp team member --help

Usage:	amp team member [OPTIONS] COMMAND

Manage team members

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
  -s, --server string   Specify server (host:port)

Commands:
  add         Add one or more members
  ls          List members
  rm          Remove one or more members

Run 'amp team member COMMAND --help' for more information on a command.
```

### Examples

* To add a member to a team:
```
$ amp team member add johndoe
team name: team
Member(s) have been added to team.
```
> NOTE: The member to be added to the team must be a existing and verified user account.

* To list members in a team:
```
$ amp team member ls --team team
team name: team
MEMBER
sample
johndoe
...
```

* To remove a member from a team:
```
$ amp team member rm johndoe
team name: team
johndoe
```

## Team Resource Management Commands

The `amp team resource` command is used to manage all team resource related operations for AMP.

### Usage

```
$ amp team resource --help

Usage:	amp team resource [OPTIONS] COMMAND

Manage team resources

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
  -s, --server string   Specify server (host:port)

Commands:
  add         Add one or more resources
  ls          List resources
  perm        Change permission level over a resource
  rm          Remove one or more resources

Run 'amp team resource COMMAND --help' for more information on a command.
```

### Examples

* To add a resource to a team:
```
$ amp team resource add 93fce7d3f...
team name: team
Resource(s) have been added to team.
```
> NOTE: Resources can be stacks.

* To list resources available to a team:
```
$ amp team resource ls
team name: team
RESOURCE ID                                                        PERMISSION LEVEL
93fce7d3fc8ada786c7db6956849343bcc5700f65d8b1512523561166a2ec455   TEAM_READ
```

* To change the permission level of a resource:
```
$ amp team resource perm 93fce7d3f... write
team name: team
Permission level has been changed.
```
The default permission level of a resource is `read`. The permission level can be changed to `read`, `write` or `admin`.

* To remove a resource from a team:
```
$ amp team resource rm --team team 93fce7d3f...
team name: team
93fce7d3fc8ada786c7db6956849343bcc5700f65d8b1512523561166a2ec455
```
