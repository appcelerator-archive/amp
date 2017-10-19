FROM golang:alpine AS builder
RUN apk --update add git gcc musl-dev linux-headers bash
# master as of Oct 19 2017
ENV VERSION 90bb0524feb1bb81901bac88fb25a3c3f100b6a7
RUN go get -d github.com/google/cadvisor
WORKDIR /go/src/github.com/google/cadvisor
RUN git checkout ${VERSION}
RUN build/assets.sh
RUN build/build.sh

FROM alpine:3.6
RUN apk --no-cache add ca-certificates wget device-mapper findutils
RUN apk --no-cache add zfs --repository http://dl-3.alpinelinux.org/alpine/edge/main/
RUN apk --no-cache add thin-provisioning-tools --repository http://dl-3.alpinelinux.org/alpine/edge/main/
RUN echo 'hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4' >> /etc/nsswitch.conf
COPY --from=builder /go/src/github.com/google/cadvisor/cadvisor /usr/bin/
EXPOSE 8080
ENTRYPOINT ["/usr/bin/cadvisor", "-logtostderr"]
