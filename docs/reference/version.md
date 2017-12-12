# AMP Version Command

The `amp version` command displays the current version of AMP.

## Usage

```
$ amp version --help

Usage:	amp version [flags]

Show version information

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
  -s, --server string   Specify server (host:port)
```

## Examples

> TIP: It would be useful to set an alias for `amp` as `alias amp='amp -k'` in case the certificate on the server is not valid.

> NOTE: For the purpose of illustration, we will use the local cluster (which is default) for running AMP commands.

* To view the client and server version for AMP.
```
$ amp version
Client:
 Version:       v0.10.0-dev
 Build:         6f590348
 Server:        127.0.0.1:50101
 Go version:    go1.8
 OS/Arch:       darwin/amd64

Server:
 Version:       v0.10.0-dev
 Build:         04daab98
 Go version:    go1.8.1
 OS/Arch:       linux/amd64
```

> NOTE: The above example is just a sample of the output of `amp version` command. The output will vary according to the version of AMP used.

* Viewing the version when the target server doesn't exist.
```
amp version
Client:
 Version:       v0.12.0-dev
 Build:         4166f646
 Server:        127.0.0.1:50101
 Go version:    go1.8
 OS/Arch:       darwin/amd64

Server:         not connected
Error:          unable to establish grpc connection: context deadline exceeded
```
> NOTE: The above example is just a sample of the output of `amp version` command. The output will vary according to the version of AMP used.
