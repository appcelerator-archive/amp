# banner-function

banner-function is a sample function showing serverless computing features.

There are only 3 things to consider when writing a function:

- Get your input parameter from the standard input
- Write your output to the standard output
- If needed, log to the standard error

# Preparing deployment

## Build the `banner-function` image

Run `make image`, it will create a docker image with the tag `appcelerator/banner-function:latest`

## Push your image to a public registry (e.g. DockerHub)

In order to use your function, you first need to push it to a public registry:

    $ docker tag appcelerator/banner-function [docker_id]/banner-function
    $ docker push [docker_id]/banner-function

# Local deployment

## Registering your function

In order to register your function, you need to run the following command:

    $ amp -s localhost fn create banner [docker_id]/banner-function

## Invoking your function via HTTP

In order to invoke a function, you can POST an HTTP request to `https://faas.local.atomiq.io/<function>`. Calls block until the function sends a response.
Invoke your test function like this:

    $ echo 'ATOMIQ rocks!' | curl -k https://faas.local.atomiq.io/banner --data-binary @-

The `@-` parameter tells `curl` to read from the standard input.

# Cloud deployment

## Registering your function

In order to register your function, you need to run the following command:

    $ amp fn create banner appcelerator/banner-function

## Invoking your function via HTTP

In order to invoke a function, you can POST an HTTP request to `https://faas.cloud.atomiq.io/<function>`. Calls block until the function sends a response.
Invoke your test function like this:

    $ echo 'ATOMIQ rocks!' | curl -k https://faas.cloud.atomiq.io/banner --data-binary @-

The `@-` parameter tells `curl` to read from the standard input.
