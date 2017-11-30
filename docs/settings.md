## Settings Command

The `amp settings` command displays settings for AMP CLI. It can be accessed by storing the value in a settings file or passing it as a command-line argument.

### Usage

```
$ amp settings --help

Usage:	amp settings [flags]

Display AMP settings

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
  -s, --server string   Specify server (host:port)
```

In order to start using the CLI settings:
- You can either:
  - Create a new directory called `.amp` in the current working (local) directory and add the config file `amp.yml` in this directory, OR
  - Create a new directory called `amp` in the `$HOME/.config` directory and add the config file `amp.yml` in this directory (This is the default location).
- Add values to the file in format:
  - `Variable: value`

> NOTE: For the moment, the settings file only stores the `Server` parameter,
which is used to point the CLI to a target cluster. More will be added in future releases.

### Examples

* The settings file is placed in the directory `$HOME/.config/amp`:
```
Server: local.appcelerator.io
```

* You can view the settings of AMP using:
```
$ amp settings
[local.appcelerator.io:50101]
Settings file: /Users/Username/.config/amp/amp.yml
AMP Settings:
  Server:        local.appcelerator.io:50101
```
