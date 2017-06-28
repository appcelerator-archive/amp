FROM alpine:3.6

RUN echo "@edge http://nl.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories && \
    echo "@edgetesting http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories

RUN apk --update add ca-certificates sudo ansible@edge py-boto@edge py-boto3@edgetesting && \
    rm -rf /var/cache/apk/*

CMD [ "ansible-playbook", "--version" ]
