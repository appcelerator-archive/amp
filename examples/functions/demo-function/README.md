# amp-demo-function

amp-demo-function is a sample function showing AMP serverless computing features.

There are only 3 things to consider when writing a function:

- Get your input parameter from the standard input
- Write your output to the standard output
- If needed, log to the standard error

# Installing locally

Run `make`, it will install `amp-demo-function` in your `$GOPATH/bin`

# Building an image

Run `make image`, it will create a docker image with the tag `appcelerator/amp-demo-function:latest`

# Testing locally

After installing locally, you can test the function like this:

    $ cat Dockerfile | amp-demo-function
