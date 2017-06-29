# =============================
FROM golang:1.8-alpine as build
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

# build
RUN go build -o /tmp/aws cmd/main.go

# =============================
FROM alpine
# =============================
RUN apk add --no-cache ca-certificates
COPY --from=build /tmp/aws /usr/local/bin/aws
ENTRYPOINT [ "aws" ]

