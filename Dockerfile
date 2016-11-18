# appclerator/protoc is based on alpine and includes latest go and protoc
FROM appcelerator/protoc:0.2.0
RUN echo "@community http://nl.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
RUN apk --no-cache add make bash git docker@community
WORKDIR /go/src/github.com/appcelerator/amp
COPY . /go/src/github.com/appcelerator/amp
RUN make install-host
EXPOSE 50101
ENTRYPOINT []
CMD [ "/go/bin/amplifier", "--service"]
