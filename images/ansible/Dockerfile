FROM alpine:3.6

RUN echo "@edge http://nl.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories && \
    echo "@edgetesting http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories

RUN apk --update add ca-certificates sudo ansible@edge py2-ansible-lint@edgetesting py-boto@edge py2-boto3@edgetesting py2-futures@edgetesting py2-s3transfer@edgetesting py2-botocore@edgetesting && \
    rm -rf /var/cache/apk/*

CMD [ "ansible-playbook", "--version" ]
