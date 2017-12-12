# User Management Commands

The `amp user` command is used to manage all user related operations for AMP.

Other user-related commands that are not managed by `amp user` include `login`, `logout` and `whoami`.

## Usage

```
$ amp user --help

Usage:	amp user [OPTIONS] COMMAND

User management operations

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
  -s, --server string   Specify server (host:port)

Commands:
  forgot-login              Retrieve account name
  get                       Get user information
  ls                        List users
  resend-verification-token Resend verification email to registered address
  rm                        Remove one or more users
  signup                    Signup for a new account
  verify                    Verify account

Run 'amp user COMMAND --help' for more information on a command.
```

> TIP: Use `-h` or `--help` option for any of the AMP commands or sub-commands to more information about the command's usage.

## Examples

> NOTE: For the purpose of illustration, we will use the local cluster (which is default) for running AMP commands. 

* To signup for a new user account:
```
$ amp user signup
[your.server.com:50101]
username: sample
email: sample@amp.com
password: [password]
[user sample @ your.server.com:50101]
Hi sample! Please check your email to complete the signup process.
```
> NOTE: If you are working on a cluster without email verification, such as a local cluster,
you will not need to verify your account as you will not be sent an email and you will be logged in automatically.
> ```
>  amp user signup
>  [127.0.0.1:50101]
>  username: sample
>  email: sample@amp.co,
>  password:
>  Verification is not necessary for this cluster.
>  Hi sample! You have been automatically logged in.
> ```

After signing up, you will then be sent an email to your registered address. In this email, you will
be sent a link to verify your account with or you can verify your account with the provided CLI command.

* To verify your account using the token in verification email:
```
$ amp user verify [token]
[your.server.com:50101]
Your account has now been activated.
```
> NOTE: If you are working on a cluster without email verification, such as a local cluster,
this command will be disabled. 
> ```
>  amp -k user verify <TOKEN>
>  [user sample @ 127.0.0.1:50101]
>  Error: `amp user verify` disabled. This cluster has no registration policy
> ```
> If you are using hosted AMP, you will need to verify your account.

* To login to your new account:
```
$ amp login
[127.0.0.1:50501]
username: sample
password: [password]
Welcome back sample!
```

* Once you have logged in, you can check who is currently logged in with `amp whoami`:
```
$ amp whoami
[user sample @ 127.0.0.1:50101]
Logged in as user: sample
```
In addition, every `amp` command will display who you are logged in as at the top of the command output.
```
$ amp
[user sample @ 127.0.0.1:50101]

Usage:  amp [OPTIONS] COMMAND
...
```

* To logout of your account:
```
$ amp logout
[user sample @ 127.0.0.1:50101]
You have been logged out!
```

* In the instance that you have forgotten the username associated with your email,
you can have the username sent to your registered email account:
```
$ amp user forgot-login sample@amp.com
Your login name has been sent to the address: sample@amp.com
```
> NOTE: If you are working on a cluster without email verification, this command will be disabled.

* To retrieve details of a specific user:
```
$ amp user get foo
```

* To retrieve a list of users:
```
$ amp user ls
```

* To remove a user:
```
$ amp user rm foo
```
This command only allows you to delete your own account. If you try to delete another user, you will see the following error:
```
$ amp user rm sample1
[user sample @ 127.0.0.1:50101]
Error: user not authorized
```
However, the `su` account has the privileges of removing other accounts in the cluster. 
