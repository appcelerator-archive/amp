# Developer portal stack

## Deploy configs and secrets

Change the sensitive data in these files, and use them to create the config and secret.

    amp -s $CLUSTER_URL secret create devportal-mysql-root ./devportal-mysql-root.txt
    amp -s $CLUSTER_URL config create devportal-mysql-init-sql ./devportal-mysql-init.sql

## Deploy the stack

    amp -s $CLUSTER_URL stack deploy -c ./devportal.yml devportal

## Access the UI

the UI should be available on:
- https://adminer.devportal.$CLUSTER_URL
- https://frontend.devportal.$CLUSTER_URL
