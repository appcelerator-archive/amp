# appclerator/protoc is based on alpine and includes latest go and protoc
FROM appcelerator/protoc
RUN apk --no-cache add make bash git docker
COPY . /go/src/github.com/appcelerator/amp
WORKDIR /go/src/github.com/appcelerator/amp
RUN make install-host
EXPOSE 50051
ENTRYPOINT [ "amplifier" ]
