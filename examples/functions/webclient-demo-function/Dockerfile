# golang:alpine provides an up to date go build environment
FROM golang:alpine
RUN apk --no-cache add make
WORKDIR /go/src/github.com/appcelerator/amp-demo-function
COPY . /go/src/github.com/appcelerator/amp-demo-function
RUN make
ENTRYPOINT []
CMD [ "/go/bin/amp-demo-function"]
