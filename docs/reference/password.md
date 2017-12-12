## Password Management Commands

The `amp password` command is used to manage all password related operations for AMP.

### Usage

```
$ amp password --help

Usage:	amp password [OPTIONS] COMMAND 

Password management operations

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
  -s, --server string   Specify server (host:port)

Commands:
  change      Change password
  reset       Reset password
  set         Set password

Run 'amp password COMMAND --help' for more information on a command.
```
> TIP: Use `-h` or `--help` option for any of the AMP commands or sub-commands to more information about the command's usage.

### Examples

* To update your current password:
```
$ amp password change
current password: [oldpassword]
new password: [newpassword]
Your password change has been successful.
```

* If you have forgotten your password and need it reset:
```
$ amp password reset [username]
```
> NOTE: If you are working on a cluster without email verification, such as a local cluster, this command will be disabled.

An email with instructions to reset password will be sent to the registered email address. In this email,
you will be sent a link to reset your password with or you can reset it with the provided CLI command.

* To set a new password:
```
$ amp password set --token [token]
```
> NOTE: If you are working on a cluster without email verification, such as a local cluster, this command will be disabled
