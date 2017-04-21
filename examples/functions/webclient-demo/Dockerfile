# golang:alpine provides an up to date go build environment
FROM golang:alpine
RUN apk --no-cache add make
WORKDIR /go/src/github.com/appcelerator/webclient-demo
COPY . /go/src/github.com/appcelerator/webclient-demo
RUN make
ENTRYPOINT []
CMD [ "/go/bin/webclient-demo"]
