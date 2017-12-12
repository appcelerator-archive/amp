# Settings Command

The `amp settings` command displays settings for AMP CLI. It can be accessed by storing the value in a settings file or passing it as a command-line argument.

## Usage

```
amp settings --help

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

## Examples

* If the settings file is placed in the directory `$HOME/.config/amp` with `Server: local.appcelerator.io`, you can view the settings of AMP using:
```
$ amp settings
[local.appcelerator.io:50101]
Settings file: $HOME/.config/amp/amp.yml
AMP Settings:
  Server:        local.appcelerator.io:50101
```

* If the settings file is placed in the directory `$PWD/.amp` with `Server: local.appcelerator.io`, you can view the settings of AMP using:
```
$ amp settings
[local.appcelerator.io:50101]
Settings file: $PWD/.amp/amp.yml
AMP Settings:
  Server:        local.appcelerator.io:50101
```

* If you have not defined a `settings` file:
```
$ amp settings
[127.0.0.1:50101]
Settings file: none
AMP Settings:
 Server:        127.0.0.1:50101
```

> TIP: Use `-h` or `--help` option for any of the AMP commands or sub-commands to more information about the command's usage.
