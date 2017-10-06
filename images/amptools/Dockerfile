FROM appcelerator/gotools:1.11.0
#RUN apk --no-cache add gosu iproute2
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends -y \
  man \
  pcregrep \
  iproute2 \
  gosu \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

# Add Docker client from:
# https://github.com/docker-library/docker/blob/master/17.06/Dockerfile
ENV DOCKER_BUCKET download.docker.com
ENV DOCKER_VERSION 17.09.0-ce
ENV DOCKER_CHANNEL stable
ENV DOCKER_SHA256_x86_64 a9e90a73c3cdfbf238f148e1ec0eaff5eb181f92f35bdd938fd7dab18e1c4647

ENV dockerArch=x86_64

RUN set -ex; \
	curl -fSL "https://${DOCKER_BUCKET}/linux/static/${DOCKER_CHANNEL}/${dockerArch}/docker-${DOCKER_VERSION}.tgz" -o docker.tgz && \
	sha256="DOCKER_SHA256_${dockerArch}"; sha256="$(eval "echo \$${sha256}")" && \
	echo "${sha256} *docker.tgz" | sha256sum -c - && \
	tar -xzf docker.tgz && \
	mv docker/* /usr/local/bin/ && \
	rm -rf docker docker.tgz && \
	docker -v

# originally: -rwxr-xr-x 1 root root 1687016 Jan 24  2017 /usr/sbin/gosu
# adding the sticky bit to allow users to execute command as root
RUN adduser --disabled-password --gecos "" --shell /bin/sh sudoer
RUN chgrp sudoer /usr/sbin/gosu && chmod +s /usr/sbin/gosu

# pass commands through docker-entrypoint first for special handling
# it's fine to override entrypoint if not running a docker command
COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["sh"]
