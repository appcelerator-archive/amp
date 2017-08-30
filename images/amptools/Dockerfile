FROM appcelerator/gotools:1.9.0
#RUN apk --no-cache add sudo iproute2
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends -y \
  man \
  pcregrep \
  iproute2 \
  sudo \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

# Add Docker client from:
# https://github.com/docker-library/docker/blob/master/17.06/Dockerfile
ENV DOCKER_BUCKET download.docker.com
ENV DOCKER_VERSION 17.06.1-ce
ENV DOCKER_CHANNEL stable
ENV DOCKER_SHA256_x86_64 e35fe12806eadbb7eb8aa63e3dfb531bda5f901cd2c14ac9cdcd54df6caed697
ENV DOCKER_SHA256_armel c3d7fcb09338a5a01320ef5e32fa9dc9161b1a4179700bb46509ff0e0d995054

ENV dockerArch=x86_64

RUN set -ex; \
	curl -fSL "https://${DOCKER_BUCKET}/linux/static/${DOCKER_CHANNEL}/${dockerArch}/docker-${DOCKER_VERSION}.tgz" -o docker.tgz && \
	sha256="DOCKER_SHA256_${dockerArch}"; sha256="$(eval "echo \$${sha256}")" && \
	echo "${sha256} *docker.tgz" | sha256sum -c - && \
	tar -xzf docker.tgz && \
	mv docker/* /usr/local/bin/ && \
	rm -rf docker docker.tgz && \
	docker -v

COPY sudoers /etc/sudoers.d/amp
RUN chmod 0440 /etc/sudoers.d/amp
RUN adduser --disabled-password --gecos "" --shell /bin/sh sudoer

# pass commands through docker-entrypoint first for special handling
# it's fine to override entrypoint if not running a docker command
COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["sh"]
