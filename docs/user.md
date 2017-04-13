### User Management

The `amp user` command manages all user related operations for AMP.

```
    $ amp user

    Usage:	amp user COMMAND

    User management operations

    Options:
      -h, --help            Print usage
      -s, --server string   Specify server (host:port)

    Commands:
      forgot-login Retrieve account name
      get          Get user
      ls           List users
      rm           Remove user
      signup       Signup for a new account
      verify       Verify account

    Run 'amp user COMMAND --help' for more information on a command.
```

### Examples

* To create a user:
```
    $ amp user signup
```
    [An email with a verification token will be sent to the given email address.]

* To verify newly created account:
```
    $ amp user verify
```

* To retrieve account name:
```
    $ amp user forgot-login
```
    [An email with the username will be sent to the registered email address.]

* To retrieve a list of users:
```
    $ amp user ls
```

* To retrieve details of a specific user:
```
    $ amp user get
```

* To remove of a user:
```
    $ amp user rm
    [or]
    $ amp user del
```

* To login to AMP:
```
    $ amp login
```

* To see who's currently logged in (user or org):
```
    $ amp whoami
```

* To logout of an account:
```
    $ amp logout
```
