# Secret Management Commands

The `amp secret` command manages all operations related to AMP secrets.

## Usage

```
$ amp secret --help

Usage:	amp secret [OPTIONS] COMMAND

Secret management operations

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
  -s, --server string   Specify server (host:port)

Commands:
  create      Create a secret from a file or STDIN as content
  ls          List secrets
  remove      Remove one or more secrets

Run 'amp secret COMMAND --help' for more information on a command.
```

> TIP: Use `-h` or `--help` option for any of the AMP commands or sub-commands to more information about the command's usage.

> NOTE: To be able to perform any secret related operations, you must be logged in to AMP using a verified account.

## Examples

* To create a secret:
```
$ amp secret create test ./test
c60z5iv6kguxdk5p1n4p7lulm
```
On success, this command returns an alphanumeric ID which is the secret ID.

* To list all secrets:
```
$ amp secret ls
ID                          NAME
8lpnhb7k140s65ftdjwav1k35   certificate_amp
c60z5iv6kguxdk5p1n4p7lulm   test
huns47stt3ade3neufohl5pzg   alertmanager_yml
i2gjv0m8ilhl5aqixmp0oeoy0   amplifier_yml
```

* To remove one or more secrets:
```
$ amp secret rm test
test
```
