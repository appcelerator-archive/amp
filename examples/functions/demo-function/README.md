# demo-function

demo-function is a sample function showing serverless computing features.

There are only 3 things to consider when writing a function:

- Get your input parameter from the standard input
- Write your output to the standard output
- If needed, log to the standard error

# Function development

## Installing the go binary

Run `make`, it will install `demo-function` in your `$GOPATH/bin`

## Testing the function without deploying it

After installing the binary, you can test the function like this:

    $ cat Dockerfile | demo-function

# Preparing deployment

## Build the `demo-function` image

Run `make image`, it will create a docker image with the tag `appcelerator/demo-function:latest`

## Push your image to a public registry (e.g. DockerHub)

In order to use your function, you first need to push it to a public registry:

    $ docker tag appcelerator/banner-function [docker_id]/demo-function
    $ docker push [docker_id]/demo-function

# Local deployment

## Registering your function

In order to register your function, you need to run the following command:

    $ amp -s localhost fn create demo appcelerator/demo-function

## Invoking your function via HTTP

In order to invoke a function, you can POST an HTTP request to `https://faas.local.atomiq.io/<function>`. Calls block until the function sends a response.
Invoke your test function like this:

    $ cat Makefile | curl -k https://faas.local.atomiq.io/demo --data-binary @-

The `@-` parameter tells `curl` to read from the standard input but you can also invoke your function like this:

    $ curl -k https://faas.local.atomiq.io/demo --data-binary @Makefile

The `--data-binary` parameter tells `curl` to POST the content of the file exactly as specified with no extra processing whatsoever.
Without this parameter, `curl` would pass the content of the file to the server using the content-type `application/x-www-form-urlencoded` which is not expected for amp functions.

# Cloud deployment

## Registering your function

In order to register your function, you need to run the following command:

    $ amp fn create demo appcelerator/demo-function

## Invoking your function via HTTP

In order to invoke a function, you can POST an HTTP request to `https://faas.cloud.atomiq.io/<function>`. Calls block until the function sends a response.
Invoke your test function like this:

    $ cat Makefile | curl -k https://faas.cloud.atomiq.io/demo --data-binary @-

The `@-` parameter tells `curl` to read from the standard input but you can also invoke your function like this:

    $ curl -k https://faas.cloud.atomiq.io/demo --data-binary @Makefile

The `--data-binary` parameter tells `curl` to POST the content of the file exactly as specified with no extra processing whatsoever.
Without this parameter, `curl` would pass the content of the file to the server using the content-type `application/x-www-form-urlencoded` which is not expected for amp functions.
