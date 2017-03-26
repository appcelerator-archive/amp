FROM docker:17.03-dind
RUN apk --no-cache add iproute2 bind-tools drill
RUN addgroup -S docker
ENTRYPOINT ["dockerd-entrypoint.sh", "--experimental"]
CMD []

