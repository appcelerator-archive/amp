FROM appcelerator/alpine:3.7.0

RUN apk --update add sudo ansible py2-ansible-lint@testing py-boto py2-boto3@testing py2-futures@testing py2-s3transfer@testing py2-botocore@testing && \
    rm -rf /var/cache/apk/*

CMD [ "ansible-playbook", "--version" ]
