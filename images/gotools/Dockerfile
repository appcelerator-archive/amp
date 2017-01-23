FROM golang:1.7-alpine
# This builds a convenience "all-in-one" image for go development.
# It intentionally does not remove any build prerequisites like most of our other
# images since this image is meant strictly for building things.

# The alpine-sdk adds a bunch of essential tools, such as
# ca-certificates, musl, musl-dev, openssl, curl, make, gcc, g++, libstdc++, libc-dev, tar, bzip2, git
# The rest are for building grpc or for any other convenience.
RUN apk --no-cache add alpine-sdk bash vim autoconf automake libtool

# Helper commands, like go-getter and install-protoc
ADD bin /usr/local/bin/

# Use go-getter for getting specific package versions (when possible).
# Add go packages to be installed to Gofile
ADD Gofile /go
RUN go-getter /go/Gofile

# Install glide for vendoring support
RUN curl https://glide.sh/get | sh
# Use /tmp/glide for the glide cache
VOLUME /tmp/glide

# Install protoc (download source, build and install -- required by protoc-gen-go and related tools)
RUN install-protoc

# Customize the shell in root/.bashrc, add other shell support here, etc.
ADD bash /root

# TODO: might want to use /go as the work directory...
WORKDIR /go/src
CMD [ "echo", "[gotools] specify the command to run" ]

