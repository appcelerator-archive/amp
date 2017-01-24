FROM docker:1.13-dind
RUN apk --no-cache add iproute2 bind-tools drill
ENTRYPOINT ["dockerd-entrypoint.sh", "--experimental"]
CMD []

