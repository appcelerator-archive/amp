# =============================
FROM golang:1.8-alpine as build
# =============================
# system deps
RUN apk --no-cache add git

ENV PKG github.com/appcelerator/amp/cluster/plugin/local
ENV DIR /go/src/${PKG}
COPY . ${DIR}
WORKDIR ${DIR}

ARG LDFLAGS="-s"
# build
RUN go build -ldflags "${LDFLAGS}" -o /tmp/amp-local cmd/main.go

# =============================
FROM alpine
# =============================
RUN apk add --no-cache ca-certificates
COPY --from=build /tmp/amp-local /usr/local/bin/amp-local
ENTRYPOINT [ "amp-local" ]

