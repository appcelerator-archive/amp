FROM appcelerator/alpine:20160726

RUN apk update && apk --no-cache add gpgme openjdk8-jre

ENV KAFKA_VERSION="0.10.0.1"
ENV SCALA_VERSION="2.11"

LABEL name="kafka" version="$KAFKA_VERSION"

RUN wget "http://mirror.cc.columbia.edu/pub/software/apache/kafka/$KAFKA_VERSION/kafka_${SCALA_VERSION}-$KAFKA_VERSION.tgz" -O /tmp/kafka.tgz \
&& mkdir -p /opt \
&& tar -xvzf /tmp/kafka.tgz -C /opt \
&& mv /opt/kafka_${SCALA_VERSION}-$KAFKA_VERSION /opt/kafka && rm /tmp/kafka.tgz

WORKDIR /opt/kafka

COPY ./check.sh /opt/kafka/bin/check.sh
COPY ./server.properties /opt/kafka/config/server.properties

EXPOSE 9092 4242

HEALTHCHECK --interval=5s --timeout=10s --retries=12 CMD /opt/kafka/bin/check.sh

ENV GOPATH /go
ENV PATH $PATH:/go/bin

COPY ./ /go/src/github.com/appcelerator/amp-docker-kafka

RUN apk update && \
    apk --virtual build-deps add go git make && \
    apk add curl && \
    cd /go/src/github.com/appcelerator/amp-docker-kafka && \
    go get -u github.com/Masterminds/glide && \
    glide install && \
    make install && \
    mv /go/bin/amp-docker-kafka /amp-docker-kafka && \
    apk del build-deps && cd / && rm -rf /go /var/cache/apk/* /root/.cache /root/.glide /tmp

CMD ["/amp-docker-kafka"]
