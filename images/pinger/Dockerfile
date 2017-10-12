FROM golang:alpine AS builder
RUN echo "@edgecommunity http://nl.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
RUN apk update
RUN apk add git upx@edgecommunity
RUN go get -d github.com/prometheus/client_golang/prometheus
WORKDIR /go/src/github.com/appcelerator/pinger
COPY . .
RUN go build -ldflags="-s -w"
RUN upx pinger

FROM alpine
COPY --from=builder /go/src/github.com/appcelerator/pinger/pinger /usr/local/bin/
CMD [ "pinger" ]
