AtSea Shop Demonstration Application
====================================

The AMP version of the [AtSea Shop Demonstration Application](https://github.com/dockersamples/atsea-sample-shop-app).

Run in this directory:

    $ docker secret create postgres_password postgres.password
    $ docker secret create staging_token staging.token
    $ amp stack up -c atsea.yml

The webapp will be available at [http://atsea.examples.local.appcelerator.io](http://atsea.examples.local.appcelerator.io).

## Differences with Original

- Reverse proxy service has been removed from the stack since AMP already provides this feature. The `appserver` service has been updated with the following attributes:
```
    environment:
      SERVICE_PORTS: 8080
      VIRTUAL_HOST: "https://atsea.examples.*,atsea.examples.*"
    networks:
      - default
```
- The default network `public` has been added in order to expose the `appserver` service behind the AMP reverse proxy.

- Docker secrets required for the reverse proxy have been removed accordingly.

- Placement constraints on all services have been removed allowing the stack to be deployed locally.

Credit: Docker ([LICENSE](https://github.com/dockersamples/atsea-sample-shop-app/blob/master/LICENSE))
