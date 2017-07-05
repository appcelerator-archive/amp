FROM appcelerator/gotools as build
RUN go get github.com/docker/swarmkit
WORKDIR /go/src/github.com/docker/swarmkit
# checkout the same version being used by github.com/docker/docker-ce v17.06 (api/vendor.conf)
RUN git checkout -b v17.06 79381d0840be27f8b3f5c667b348a4467d866eeb
RUN CGO_ENABLED=0 GO_LDFLAGS="-s" make binaries
RUN mv bin /tmp/bin

FROM scratch
COPY --from=build /tmp/bin/ /bin
CMD [ "/bin/swarmctl" ]
 
