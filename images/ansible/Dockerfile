FROM appcelerator/alpine:3.7.1

#RUN apk --update add sudo ansible py2-ansible-lint@testing py-boto py2-boto3@testing py2-futures@testing py2-s3transfer@testing py2-botocore@testing && \
ENV ANSIBLE_VERSION v2.5.0
RUN apk --update add sudo ansible py-boto py2-boto3@testing py2-futures@testing py2-s3transfer@testing py2-botocore@testing && \
    apk --virtual build-deps --no-cache add py2-pip git gcc && \
    pip install git+https://github.com/ansible/ansible.git@${ANSIBLE_VERSION} && \
    apk del build-deps && rm -rf /var/cache/apk/*

CMD [ "ansible-playbook", "--version" ]
