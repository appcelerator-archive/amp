# webclient-demo

webclient-demo is a sample function showing ATOMIQ serverless computing features.

> This example function will find how many Docker repositories a user or organisation has on the Docker Hub.

There are only 3 things to consider when writing a function:

- Get your input parameter from the standard input
- Write your output to the standard output
- If needed, log to the standard error

# Installing locally

Run `make`, it will install `webclient-demo` in your `$GOPATH/bin`

# Building an image

Run `make image`, it will create a docker image with the tag `appcelerator/webclient-demo:latest`

# Testing the Docker image:

```
# echo alexellis2 | docker run -i appcelerator/webclient-demo:latest
The organisation or user alexellis2 has 154 repositories on the Docker hub.
```

# Local deployment

## Push your image to the local registry

In order to use your function, you first need to push it to the local registry:

    $ docker tag appcelerator/webclient-demo localhost:5000/appcelerator/webclient-demo
    $ docker push localhost:5000/appcelerator/webclient-demo

## Registering your function

In order to register your function, you need to run the following command:

    $ amp -s localhost fn create webclient-demo appcelerator/webclient-demo

## Invoking your function via HTTP

In order to invoke a function, you can POST an HTTP request to `https://faas.local.atomiq.io/<function>`. Calls block until the function sends a response.
Invoke your test function like this:

    $ echo alexellis2 | curl -k https://faas.local.atomiq.io/webclient-demo --data-binary @-

The `@-` parameter tells `curl` to read from the standard input but you can also invoke your function like this:

# Cloud deployment

## Push your image to the ATOMIQ registry

In order to use your function, you first need to push it to `registry.cloud.atomiq.io`:

    $ docker tag appcelerator/webclient-demo registry.cloud.atomiq.io/appcelerator/webclient-demo
    $ docker push registry.cloud.atomiq.io/appcelerator/webclient-demo

## Registering your function

In order to register your function, you need to run the following command:

    $ amp fn create webclient-demo appcelerator/webclient-demo

## Invoking your function via HTTP

In order to invoke a function, you can POST an HTTP request to `https://faas.cloud.atomiq.io/<function>`. Calls block until the function sends a response.
Invoke your test function like this:

    $ echo alexellis2 | curl -k https://faas.cloud.atomiq.io/webclient-demo --data-binary @-
