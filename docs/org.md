## Organization Management Commands

The `amp org` command is used to manage all organization related operations for AMP.

### Usage

```
$ amp org --help

Usage:	amp org [OPTIONS] COMMAND 

Organization management operations

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
  -s, --server string   Specify server (host:port)

Management Commands:
  member      Manage organization members

Commands:
  create      Create organization
  get         Get organization information
  ls          List organization
  rm          Remove one or more organizations
  switch      Switch account

Run 'amp org COMMAND --help' for more information on a command.
```
>NOTE: To be able to perform any organization related operations, you must be logged in to AMP using a verified account.

>`amp org` commands that require the `--org` flag will remember the last used organization using a local config file.
If you want to override this, specify the `--org` option in the command.

### Examples

* To create an organization:
```
$ amp org create
organization name: org
email: email@org.com
Organization has been created.
```

* To retrieve details of a specific organization:
```
$ amp org get org
Organization: org
Email: email@org.com
Created On: 28 Jun 17 11:28
```

* To retrieve a list of organizations:
```
$ amp org list
ORGANIZATION   EMAIL                    CREATED ON
org            email@org.com            28 Jun 17 11:28
so             super@organization.amp   27 Jun 17 14:00
```

* To switch between accounts (user or org):
```
$ amp org switch org
You are now logged in as: org
```
Logs the user in on behalf of `org`.

* To remove an organization:
```
$ amp org remove org
org
```

## Organization Member Management Commands

The `amp org member` command is used to manage all organization member operations for AMP.

### Usage

```
$ amp org member --help

Usage:	amp org member [OPTIONS] COMMAND 

Manage organization members

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
  -s, --server string   Specify server (host:port)

Commands:
  add         Add one or more members
  ls          List members
  rm          Remove one or more members
  role        Change member role

Run 'amp org member COMMAND --help' for more information on a command.
```

### Examples

* To add a member to an organization:
```
$ amp org member add johndoe
organization name: org
Member(s) have been added to organization.
```
> NOTE: The member to be added to the organization must be a existing and verified user account.

* To list members in an organization:
```
$ amp org member ls
organization name: org
MEMBER         ROLE
sample         ORGANIZATION_OWNER
johndoe        ORGANIZATION_MEMBER
```

* To change role of a member in an organization:
```
$ amp org member role
organization name: org
member: johndoe
role: owner
Member role has been changed.
```
When a member is added to the organization, the default role is `member`. The role can be changed to `member` and `owner`.

* To remove a member from an organization:
```
$ amp org member rm johndoe
organization name: org
johndoe
```
