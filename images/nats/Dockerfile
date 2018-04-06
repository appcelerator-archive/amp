FROM golang:1.10.1-alpine AS builder

ENV VERSION v0.9.2
RUN apk update && apk add git
RUN go get -d github.com/nats-io/nats-streaming-server
WORKDIR /go/src/github.com/nats-io/nats-streaming-server
RUN git checkout ${VERSION}
RUN CGO_ENABLED=0 GOOS=linux   GOARCH=amd64         go build -v -a -tags netgo -installsuffix netgo -ldflags "-s -w -X github.com/nats-io/nats-streaming-server/version.GITCOMMIT=`git rev-parse --short HEAD`" -o /nats-streaming-server

FROM scratch
COPY --from=builder /nats-streaming-server /nats-streaming-server
# Expose client and management ports
EXPOSE 4222 6222 8222
# Run with default memory based store
ENTRYPOINT ["/nats-streaming-server", "-m", "8222"]
CMD []

