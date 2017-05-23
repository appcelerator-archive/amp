### AMP Configuration

The `amp config` command displays the AMP configuration. Currently, only the server address is part of the configuration.
It can be accessed by storing the value in a config file or passing it as a command-line argument.
The config file can be stored in one of the following paths:

- Create a new directory called `.amp` in the current working (local) directory and add the config file `amp.yml` in this directory.
- Create a new directory called `amp` in the `$HOME/.config` directory and add the config file `amp.yml` in this directory.

The default location to store the config file is in the `$HOME/.config/amp` directory


#### Usage

* As a command-line argument:
```
    $ amp -s localhost config
    [localhost:50101]
    Config:
     Server:       localhost:50101
```

* Reading from the local directory (from `./.amp/amp.yml`)

* Reading from the HOME directory (from `$HOME/.config/amp.yml`)
