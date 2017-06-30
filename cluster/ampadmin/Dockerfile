# =============================
FROM golang:1.8-alpine as build
# =============================
RUN apk --no-cache add git build-base

# go
RUN go get github.com/LK4D4/vndr
ENV PKG github.com/appcelerator/amp/cluster/ampadmin
ENV DIR /go/src/${PKG}
COPY . ${DIR}
WORKDIR ${DIR}
RUN vndr

# test
RUN go test -v -timeout 30m

# build
RUN go build -o /tmp/ampadmin ${PKG}

# =============================
FROM alpine
# =============================
RUN apk add --no-cache ca-certificates
COPY --from=build /tmp/ampadmin /usr/local/bin/ampadmin
ENTRYPOINT [ "ampadmin" ]

