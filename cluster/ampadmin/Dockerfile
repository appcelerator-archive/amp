# =============================
FROM golang:1.9-alpine as build
# =============================
RUN apk add --no-cache build-base
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache git

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
RUN CGO_ENABLED=0 go build -a -ldflags "-s" -o /tmp/ampadmin ${PKG}/cmd

# =============================
FROM scratch
# =============================
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /tmp/ampadmin /ampadmin
ENTRYPOINT [ "/ampadmin" ]

