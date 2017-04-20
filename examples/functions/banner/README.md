# banner-function

banner-function is a sample function showing serverless computing features.

There are only 3 things to consider when writing a function:

- Get your input parameter from the standard input
- Write your output to the standard output
- If needed, log to the standard error

# Building the `banner-function` image

Run `make image`, it will create a docker image with the tag `appcelerator/banner-function:latest`

# Local deployment

## Push your image to the local registry

In order to use your function, you first need to push it to the local registry:

    $ docker tag appcelerator/banner-function localhost:5000/appcelerator/banner-function
    $ docker push localhost:5000/appcelerator/banner-function

## Registering your function

In order to register your function, you need to run the following command:

    $ amp -s localhost fn create banner appcelerator/banner-function

## Invoking your function via HTTP

In order to invoke a function, you can POST an HTTP request to `localhost:50102/<function>`. Calls block until the function sends a response.
Invoke your test function like this:

    $ echo 'Atomiq Rocks!!' | curl localhost:50102/banner --data-binary @-

The `@-` parameter tells `curl` to read from the standard input.

# Cloud deployment

## Push your image to the atomiq registry

In order to use your function, you first need to push it to `registry.cloud.atomiq.io`:

    $ docker tag appcelerator/banner-function registry.cloud.atomiq.io/appcelerator/banner-function
    $ docker push registry.cloud.atomiq.io/appcelerator/banner-function

## Registering your function

In order to register your function, you need to run the following command:

    $ amp fn create banner appcelerator/banner-function

## Invoking your function via HTTP

In order to invoke a function, you can POST an HTTP request to `localhost:50102/<function>`. Calls block until the function sends a response.
Invoke your test function like this:

    $ echo 'Atomiq Rocks!!' | curl localhost:50102/banner --data-binary @-

The `@-` parameter tells `curl` to read from the standard input.
