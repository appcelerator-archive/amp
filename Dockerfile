# appclerator/protoc is based on alpine and includes latest go and protoc
FROM appcelerator/protoc
RUN apk --no-cache add make bash git docker curl
RUN curl https://glide.sh/get | sh
WORKDIR /go/src/github.com/appcelerator/amp
COPY glide.lock /go/src/github.com/appcelerator/amp/
COPY glide.yaml /go/src/github.com/appcelerator/amp/
RUN glide install
COPY . /go/src/github.com/appcelerator/amp
RUN make install-host
EXPOSE 50101
ENTRYPOINT []
CMD [ "amplifier" ]
