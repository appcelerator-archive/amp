## User Management Commands

The `amp user` command is used to manage all user related operations for AMP.

Other user-related commands that aren't managed by `amp user` include `login`, `logout` and `whoami`.

### Usage

```
 $ amp user --help

Usage:	amp user [OPTIONS] COMMAND 

User management operations

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
  -s, --server string   Specify server (host:port)

Commands:
  forgot-login Retrieve account name
  get          Get user information
  ls           List users
  rm           Remove one or more users
  signup       Signup for a new account
  verify       Verify account

Run 'amp user COMMAND --help' for more information on a command.
```

### Examples

#### Signing up and Logging in

* To signup for a new user account:
```
$ amp user signup
username: sample
email: sample@user.com
password: [password]
Hi sample! Please check your email to complete the signup process.
```
>NOTE: If you are working on a cluster without email verification, such as a local cluster,
you will not need to verify your account as you will not be sent an email and you will be logged in automatically.

After signing up, you will then be sent an email to your registered address. In this email, you will
be sent a link to verify your account with or you can verify your account with the provided CLI command.

* To verify your account using the token in verification email.
```
$ amp user verify [token]
Your account has now been activated.
```
>NOTE: If you are working on a cluster without email verification, such as a local cluster,
this command will be disabled. If you are using hosted AMP, you will need to verify your account.

* To login to your new account.
```
$ amp login
username: sample
password: [password]
Welcome back sample!
```

* Once you have logged in, you can check you are logged with `amp whoami`
```
$ amp whoami
[user sample @ localhost:50101]
Logged in as user: sample
```
In addition, every `amp` command will display who you are logged in as at
the top of the command output:
```
$ amp
[user sample @ localhost:50101]

Usage:  amp [OPTIONS] COMMAND
...
```

* To logout of your account
```
$ amp logout
You have been logged out!
```

#### Forgotten your username

* In the instance that you have forgotten the username associated with your email,
you can have the username sent to your registered email account:
```
$ amp user forgot-login sample@user.com
Your login name has been sent to the address: sample@user.com
```
>NOTE: If you are working on a cluster without email verification, this command will be disabled.

#### User information

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
