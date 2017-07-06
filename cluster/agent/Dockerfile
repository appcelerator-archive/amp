# =============================
FROM golang:1.8-alpine as build
# =============================
RUN apk add --no-cache build-base
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache git

# go
# ensure any binaries are placed in /tmp/bin
ENV TMPBIN /tmp/bin
ENV TARGET ${TMPBIN}/ampctl
ENV PKG github.com/appcelerator/amp/cluster/agent
ENV PKGCMD ${PKG}/cmd
ENV PKGROOT /go/src/${PKG}
ENV PKGCMD ${PKG}/cmd
ENV LDFLAGS "-s"

RUN go get github.com/LK4D4/vndr

COPY . ${PKGROOT}
WORKDIR ${PKGROOT}

# get vendor packages
# TODO: uncomment after development (convenient to just existing vendor for now)
#RUN vndr

# test
RUN go test -v -timeout 30m ${PKG}/admin

# build
RUN CGO_ENABLED=0 go build -a -ldflags ${LDFLAGS} -o ${TARGET} ${PKGCMD}

# =============================
FROM scratch
# =============================
ENV SWARM_SOCKET "/var/run/docker/swarm/control.sock"

COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /tmp/bin/ /bin

ENTRYPOINT [ "/bin/ampctl" ]

