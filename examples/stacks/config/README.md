# Config Example

## Overview

In this example, you will store config data in the swarm that will be available to a service.

## Add data that represents a configuration file

    $ echo "foo=bar" | amp config create demo.config -

## Verify config was stored

    $ amp config ls
    # demo.config should appear in the list

## Create a stack file with a simple service to demonstrate it has access to the config

See the existing `stack.yml` file.

```
version: "3.3"

services:

  config:
    image: "alpine"
    command: "cat /demo.config"
    deploy:
      restart_policy:
        condition: none
    configs:
      - source: demo_config
        target: /demo.config

configs:
  demo_config:
    external: true

```

## Deploy the stack

TODO: this works with docker -- get unsupported compose version in amp, working to resolve right now

    $ amp stack deploy -c stack.yml demo
    Creating service demo_config

## Test the services

```
$ docker service logs demo_config
demo_config.1.tiwkpbjyyyg0@moby    | foo=bar
```


## Clean up

    $ amp stack rm demo
    $ amp config rm demo_config


