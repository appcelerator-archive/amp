## Version Command

The `amp version` command displays the current version of AMP.

### Usage

```
$ amp version --help

Usage:	amp version

Show version information

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
  -s, --server string   Specify server (host:port)
```

### Examples

* To view the client and server version for AMP.
```
    $ amp version
    Client:
     Version:       v0.10.0-dev
     Build:         6f590348
     Server:        localhost:50101
     Go version:    go1.8
     OS/Arch:       darwin/amd64

    Server:
     Version:       v0.10.0-dev
     Build:         04daab98
     Go version:    go1.8.1
     OS/Arch:       linux/amd64
```

* Viewing the version when the target server doesn't exist.
```
    $ amp version
    Client:
      Version:       v0.12.0-dev
      Build:         4166f646
      Server:        localhost:50101
      Go version:    go1.8
      OS/Arch:       darwin/amd64

    Server:         not connected
```
