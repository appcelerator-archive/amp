### Password Management

The `amp password` command manages password-related operations for any AMP account.

```
    $ amp password

    Usage:	amp password [OPTIONS] COMMAND [ARGS...]

    Password management operations

    Options:
          --help   Print usage

    Commands:
      change      Change password
      reset       Reset password
      set         Set password

    Run 'amp password COMMAND --help' for more information on a command.
```

### Examples

* To update current password:
```
    $ amp password change
```

* To reset password:
```
    $ amp password reset
```
    [An email with instructions to reset password will be sent to the registered email address.]


* To set a new password:
```
    $ amp password set
```
