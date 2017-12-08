FROM appcelerator/gotools:1.15.0

RUN rm -rf /etc/apk/repositories
RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories
RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
RUN apk update
RUN apk add gosu docker make nodejs jq

## swagger-combine for API documentation
RUN npm install --save swagger-combine

## originally: -rwxr-xr-x 1 root root 1687016 Jan 24  2017 /usr/sbin/gosu
## adding the sticky bit to allow users to execute command as root
RUN adduser -D -g "" -s /bin/sh sudoer
RUN chgrp sudoer /usr/bin/gosu && chmod +s /usr/bin/gosu

# pass commands through docker-entrypoint first for special handling
# it's fine to override entrypoint if not running a docker command
COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["sh"]
