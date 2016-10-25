# Stack file reference

A stack is a collection of services that run in a specific environment on AMP under
[Docker swarm mode](https://docs.docker.com/engine/swarm/). A stack file is used to
define this environment in [YAML](http://yaml.org/) format. It is intentionally very
similar to a [Compose file](https://docs.docker.com/compose/compose-file/)
(`docker-compose.yml`) and to a Docker Cloud
[stack file](https://docs.docker.com/docker-cloud/apps/stack-yaml-reference/)
(`docker-cloud.yml`), but with a number of unique extensions.

## Example

```
version: '1'  # stack file version (default 1)
services:
  web: # service name
    image: appcelerator.io/amp-demo
    public: # publication specification
      - name: www
        protocol: tcp
        publish_port: 80
        internal_port: 3000
    replicas: 3
    environment:
      REDIS_PASSWORD: password
    labels: # service labels
      io.appcelerator.amp.label: amp-demo
    container_labels:
      io.appcelerator.amp.label: web
  redis:
    image: redis
    mode: global # default replicated
    environment:
      - PASSWORD=password
    labels:
      - "io.appcelerator.amp.label=redis"
    container_labels:
      - "io.appcelerator.amp.label=db"
    networks:
      app-net:
        aliases:
         - stack1-redis
networks:
  app-net:
    driver: bridge
    driver_opts:
      com.docker.network.enable_ipv6: "true"
    ipam:
      driver: default
      config:
      - subnet: 172.16.238.0/24
        gateway: 172.16.238.1
      - subnet: 2001:3984:3989::/64
        gateway: 2001:3984:3989::1  
```

## Reference


## version

The stack file version. The default is "1".

## services

The `services` key is a map of names for each service specification.

### image (required)

The image used to deploy this service. This is the only mandatory key.

### public

The public specification for exposing this service externally.

...

## networks

...
