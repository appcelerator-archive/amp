# =============================
FROM golang:1.8-alpine as build
# =============================

# build args
# docker build --build-arg package="foo/bar" ...
# package is mandatory
ARG package
RUN test -n "$package"

# system packages
RUN apk --no-cache add git

# go
RUN go get github.com/LK4D4/vndr
ENV dir /go/src/${package}
RUN echo ${dir}
COPY . ${dir}
WORKDIR ${dir}
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
