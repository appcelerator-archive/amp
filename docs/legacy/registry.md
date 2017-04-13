AMP registry
============

# Usage

AMP comes with an internal registry available through the AMP cli and providing images to the Swarm cluster. It is meant to host the Docker images used by the application services (defined in a stack).

## CLI

### localhost (development)

#### Push an image

```amp registry push appcelerator/pinger:latest```

### remote cluster

the cluster should have a FQDN with sub level aliases. Let's say the domain is amp.example.com, the registry is available on registry.amp.example.com.

If there's no legit certificate for this registry with this name (default use case), this url should be declared as insecure registry in your Docker configuration.

#### Configuration on Linux

```systemctl edit docker.service```

add the block (or adapt the existing file if you already have a customization)

```
[Service]
Environment="INSECURE_REGISTRY=registry.amp.example.com"
ExecStart=-
ExecStart=/usr/bin/dockerd $OPTIONS \
          $INSECURE_REGISTRY
```

#### Configuration on Mac OS

Go in Preferences, advanced tab, add an insecure registry.

#### Push an image

```amp registry push --domain amp.example.com  appcelerator/pinger:latest```

#### Check the registry catalog

```amp registry ls --domain amp.example.com```

## Stack definition

The internal registry images are available with the local alias local.appcelerator.io.
In the stack definition, use this alias in the hostname part of the image:

```
myservice1:
    image: local.appcelerator.io/pinger:latest
```
