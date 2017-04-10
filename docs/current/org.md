### Organization Management

The `amp org` command manages all organization related operations for AMP.

```
    $ amp org

    Usage:	amp org [OPTIONS] COMMAND [ARGS...]

    Organization management operations

    Options:
          --help   Print usage

    Management Commands:
      member      Manage organization members

    Commands:
      create      Create organization
      get         Get organization
      ls          List organization
      rm          Remove organization

    Run 'amp org COMMAND --help' for more information on a command.
```

### Examples

To be able to perform any organization related operations, you must be logged in to AMP using a verified account.

* To create an organization:
```
    $ amp org create
```

* To retrieve a list of organizations:
```
    $ amp org ls
```

* To retrieve details of a specific organization:
```
    $ amp org get
```

* To remove an organization:
```
    $ amp org rm
    [or]
    $ amp org del
```

#### Organization Member Management Commands

```
    $ amp org member

    Usage:	amp org member [OPTIONS] COMMAND [ARGS...]

    Manage organization members

    Options:
          --help   Print usage

    Commands:
      add         Add member to organization
      ls          List members
      rm          Remove member from organization
      role        Change member role

    Run 'amp org member COMMAND --help' for more information on a command.
```

#### Examples

* To add a member to an organization:
```
    $ amp org member add
```

* To list members in an organization:
```
    $ amp org member ls
```

* To change role of a member in an organization:
```
    $ amp org member role
```

* To remove a member from an organization:
```
    $ amp org member rm
    [or]
    $ amp org member del
```
