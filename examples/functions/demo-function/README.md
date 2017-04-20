# demo-function

demo-function is a sample function showing serverless computing features.

There are only 3 things to consider when writing a function:

- Get your input parameter from the standard input
- Write your output to the standard output
- If needed, log to the standard error

# Installing the go binary

Run `make`, it will install `demo-function` in your `$GOPATH/bin`

# Building the `demo-function` image

Run `make image`, it will create a docker image with the tag `appcelerator/demo-function:latest`

# Testing the function without deploying it

After installing locally, you can test the function like this:

    $ cat Dockerfile | demo-function

# Local deployment

## Push your image to the local registry

In order to use your function, you first need to push it to the local registry:

    $ docker tag appcelerator/demo-function localhost:5000/appcelerator/demo-function
    $ docker push localhost:5000/appcelerator/demo-function

## Registering your function

In order to register your function, you need to run the following command:

    $ amp -s localhost fn create test appcelerator/demo-function

## Invoking your function via HTTP

In order to invoke a function, you can POST an HTTP request to `localhost:50102/<function>`. Calls block until the function sends a response.
Invoke your test function like this:

    $ cat Makefile | curl localhost:50102/test --data-binary @-

The `@-` parameter tells `curl` to read from the standard input but you can also invoke your function like this:

    $ curl localhost:50102/test --data-binary @Makefile

The `--data-binary` parameter tells `curl` to POST the content of the file exactly as specified with no extra processing whatsoever.
Without this parameter, `curl` would pass the content of the file to the server using the content-type `application/x-www-form-urlencoded` which is not expected for amp functions.

# Cloud deployment

## Push your image to the atomiq registry

In order to use your function, you first need to push it to `registry.cloud.atomiq.io`:

    $ docker tag appcelerator/demo-function registry.cloud.atomiq.io/appcelerator/demo-function
    $ docker push registry.cloud.atomiq.io/appcelerator/demo-function

## Registering your function

In order to register your function, you need to run the following command:

    $ amp fn create test appcelerator/demo-function

## Invoking your function via HTTP

In order to invoke a function, you can POST an HTTP request to `cloud.atomiq.io:50102/<function>`. Calls block until the function sends a response.
Invoke your test function like this:

    $ cat Makefile | curl cloud.atomiq.io:50102/test --data-binary @-

The `@-` parameter tells `curl` to read from the standard input but you can also invoke your function like this:

    $ curl cloud.atomiq.io:50102/test --data-binary @Makefile

The `--data-binary` parameter tells `curl` to POST the content of the file exactly as specified with no extra processing whatsoever.
Without this parameter, `curl` would pass the content of the file to the server using the content-type `application/x-www-form-urlencoded` which is not expected for amp functions.
