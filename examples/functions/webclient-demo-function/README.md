# webclient-demo-function

webclient-demo-function is a sample function showing AMP serverless computing features.

> This example function will find how many Docker repositories a user or organisation has on the Docker Hub.

There are only 3 things to consider when writing a function:

- Get your input parameter from the standard input
- Write your output to the standard output
- If needed, log to the standard error

# Installing locally

Run `make`, it will install `webclient-demo-function` in your `$GOPATH/bin`

# Building an image

Run `make image`, it will create a docker image with the tag `appcelerator/webclient-demo-function:latest`

# Testing the Docker image:

```
# echo alexellis2 | docker run -i appcelerator/amp-webclientdemo-function:latest
The organisation or user alexellis2 has 97 repositories on the Docker hub.
```

# Testing locally

After installing locally, you can test the function like this:

    $ cat Dockerfile | amp-demo-function
