FROM golang:1.9-alpine AS build
RUN apk add --update git make gcc musl-dev wget ca-certificates openssl libvirt-dev libvirt-lxc libvirt-qemu git openssh
ENV GOPATH /go
ENV PATH /go/bin:$PATH
ARG INFRAKIT_VERSION
COPY build-infrakit /usr/local/bin/
RUN go get -u -v github.com/docker/infrakit/cmd/infrakit
RUN if [ -n "${INFRAKIT_VERSION}" ]; then cd /go/src/github.com/docker/infrakit/ && git checkout ${INFRAKIT_VERSION}; fi
RUN build-infrakit linux
RUN ls -l /build/infrakit

FROM appcelerator/alpine:3.6.0
RUN apk add --update ca-certificates openssl openssh
RUN mkdir -p /infrakit/plugins /infrakit/configs /infrakit/logs /infrakit/cli /infrakit/instance/terraform
VOLUME /infrakit
WORKDIR /infrakit
ENV INFRAKIT_HOME /infrakit
# Defined in pkg/discovery
ENV INFRAKIT_PLUGINS_DIR /infrakit/plugins
# Defined in pkg/cli
ENV INFRAKIT_CLI_DIR /infrakit/cli
# Defined in cmd/infrakit/playbook
ENV INFRAKIT_PLAYBOOKS_FILE /infrakit/playbooks
# When using the manager 'os' option
ENV INFRAKIT_LEADER_FILE /infrakit/leader
ENV INFRAKIT_STORE_DIR /infrakit/configs
# Default options for the manager plugin
ENV INFRAKIT_MANAGER_BACKEND swarm
ENV INFRAKIT_MANAGER_GROUP group-stateless
# Default options for the aws metadata plugin
ENV INFRAKIT_AWS_METADATA_TEMPLATE_URL file:///infrakit/aws-metadata-export.ikt
ENV INFRAKIT_AWS_METADATA_POLL_INTERVAL 300s
COPY --from=build /build/infrakit /usr/local/bin/infrakit
COPY aws-metadata-export.ikt ${INFRAKIT_HOME}/
CMD ["infrakit", "plugin", "start", "vars", "manager", "group", "swarm", "vanilla", "combo", "aws"]
#HEALTHCHECK --interval=5s --retries=3 --timeout=2s CMD infrakit var keys >/dev/null 2>&1
