# =============================
FROM golang:1.8-alpine as build
# =============================
RUN apk --no-cache add git

# go
RUN go get github.com/LK4D4/vndr
ENV PKG github.com/appcelerator/amp/cluster/envoy
ENV DIR /go/src/${PKG}
COPY . ${DIR}
WORKDIR ${DIR}
RUN vndr

# test
RUN go test -v -timeout 30m

# build
RUN go build -o /tmp/envoy cmd/main.go

# =============================
FROM alpine
# =============================
RUN apk add --no-cache ca-certificates
COPY --from=build /tmp/envoy /usr/local/bin/envoy
ENTRYPOINT [ "envoy" ]

