# golang:alpine provides an up to date go build environment
FROM golang:alpine
RUN echo "@community http://nl.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
RUN apk --no-cache add bash alpine-sdk
WORKDIR /go/src/github.com/appcelerator/amp
COPY . /go/src/github.com/appcelerator/amp
ARG BUILD=unknown
RUN make BUILD=$BUILD install
EXPOSE 50101
ENTRYPOINT []
CMD [ "/go/bin/amplifier"]
