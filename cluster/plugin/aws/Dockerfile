# =============================
FROM golang:1.9-alpine as build
# =============================
# system deps
RUN apk --no-cache add git

# go deps
RUN go get github.com/LK4D4/vndr
ENV PKG github.com/appcelerator/amp/cluster/plugin/aws
ENV DIR /go/src/${PKG}
COPY . ${DIR}
WORKDIR ${DIR}
RUN vndr

ARG LDFLAGS="-s"
# build
RUN go build -ldflags "$LDFLAGS" -o /tmp/aws cmd/main.go

# =============================
FROM alpine
# =============================
RUN apk add --no-cache ca-certificates
COPY --from=build /tmp/aws /usr/local/bin/aws
ENTRYPOINT [ "aws" ]

