FROM golang:alpine
RUN echo "@community http://nl.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
RUN apk --no-cache add make bash git alpine-sdk
WORKDIR /go/src/github.com/appcelerator/amp
COPY . /go/src/github.com/appcelerator/amp
