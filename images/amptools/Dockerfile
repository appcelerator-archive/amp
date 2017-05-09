FROM appcelerator/gotools:1.8.0
#RUN apk --no-cache add sudo iproute2
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
  man \
  pcregrep \
  iproute2 \
  sudo \
  && rm -rf /var/lib/apt/lists/*

# Add Docker client from:
# https://github.com/docker-library/docker/blob/master/17.05-rc/Dockerfile
ENV DOCKER_BUCKET test.docker.com
ENV DOCKER_VERSION 17.05.0-ce-rc2
ENV DOCKER_SHA256_x86_64 fc2521a30797212410f46faece24204b9528990dac823c145b530c05acd50e9d
ENV DOCKER_SHA256_armel c4427822445e5f8c970ded1c106780a32641a7344c92a8dc57a9be4e0ea9cba0

ENV dockerArch=x86_64

RUN set -ex; \
	curl -fSL "https://${DOCKER_BUCKET}/builds/Linux/${dockerArch}/docker-${DOCKER_VERSION}.tgz" -o docker.tgz; \
# /bin/sh doesn't support ${!...} :(
	sha256="DOCKER_SHA256_${dockerArch}"; sha256="$(eval "echo \$${sha256}")"; \
	echo "${sha256} *docker.tgz" | sha256sum -c -; \
	tar -xzvf docker.tgz; \
	mv docker/* /usr/local/bin/; \
	rmdir docker; \
	rm docker.tgz; \
	docker -v

COPY sudoers /etc/sudoers.d/amp
RUN chmod 0440 /etc/sudoers.d/amp
RUN adduser --shell /bin/sh sudoer

# pass commands through docker-entrypoint first for special handling
# it's fine to override entrypoint if not running a docker command
COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["sh"]
