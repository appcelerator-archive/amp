### Organization Management

The `amp org` command manages all organization related operations for AMP.

```
    $ amp org

    Usage:	amp org COMMAND

    Organization management operations

    Options:
      -h, --help            Print usage
      -s, --server string   Specify server (host:port)

    Management Commands:
      member      Manage organization members

    Commands:
      create      Create organization
      get         Get organization
      ls          List organization
      rm          Remove organization
      switch      Switch account

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

* To switch between accounts (user or org):
```
    $ amp org switch
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

    Usage:	amp org member COMMAND

    Manage organization members

    Options:
      -h, --help            Print usage
      -s, --server string   Specify server (host:port)

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
