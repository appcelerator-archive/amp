From appcelerator/alpine:3.5.1
RUN echo "http://nl.alpinelinux.org/alpine/v3.4/community" >> /etc/apk/repositories 
RUN apk --update add docker
COPY ./test.sh /bin
CMD ["/bin/test.sh"]
