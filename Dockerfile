# appcelerator/protoc is based on alpine and includes latest go and protoc
FROM appcelerator/protoc:0.3.0
RUN echo "@community http://nl.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
RUN apk --no-cache add bash alpine-sdk
WORKDIR /go/src/github.com/appcelerator/amp
COPY . /go/src/github.com/appcelerator/amp
ARG BUILD=unknown
RUN make BUILD=$BUILD install-host
EXPOSE 50101
ENTRYPOINT []
CMD [ "/go/bin/amplifier"]
