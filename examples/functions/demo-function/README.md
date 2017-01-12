# demo-function

demo-function is a sample function showing AMP serverless computing features.

There are only 3 things to consider when writing a function:

- Get your input parameter from the standard input
- Write your output to the standard output
- If needed, log to the standard error

# Installing locally

Run `make`, it will install `demo-function` in your `$GOPATH/bin`

# Building an image

Run `make image`, it will create a docker image with the tag `appcelerator/amp-demo-function:latest`

# Testing locally

After installing locally, you can test the function like this:

    $ cat Dockerfile | demo-function

# Registering your function
In order to register your function against amp, you need to run the following command:

    $ amp fn create test appcelerator/amp-demo-function

# Invoking your function via HTTP
In order to invoke a function, you can POST an HTTP request to `localhost:4242/<function>` (to be changed). Calls block until the function sends a response.
Invoke your test function like this:

    $ cat Makefile | curl localhost:4242/test --data-binary @-

The `@-` parameter tells `curl` to read from the standard input but you can also invoke your function like this:

    $ curl localhost:4242/test --data-binary @Makefile

