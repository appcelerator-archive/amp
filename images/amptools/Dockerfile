FROM appcelerator/gotools:1.7.0
#RUN apk --no-cache add sudo iproute2

# Add Docker client from:
# https://github.com/docker-library/docker/blob/master/17.03/Dockerfile
ENV DOCKER_BUCKET get.docker.com
ENV DOCKER_VERSION 17.03.0-ce
ENV DOCKER_SHA256 4a9766d99c6818b2d54dc302db3c9f7b352ad0a80a2dc179ec164a3ba29c2d3e

RUN set -x \
	&& curl -sfSL "https://${DOCKER_BUCKET}/builds/Linux/x86_64/docker-${DOCKER_VERSION}.tgz" -o docker.tgz \
	&& echo "${DOCKER_SHA256} *docker.tgz" | sha256sum -c - \
	&& tar -xzvf docker.tgz \
	&& mv docker/* /usr/local/bin/ \
	&& rmdir docker \
	&& rm docker.tgz \
	&& docker -v

COPY docker-entrypoint.sh /usr/local/bin/
COPY sudoers /etc/sudoers.d/amp
RUN chmod 0440 /etc/sudoers.d/amp
RUN adduser --shell /bin/sh sudoer

# pass commands through docker-entrypoint first for special handling
# it's fine to override entrypoint if not running a docker command
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["sh"]

