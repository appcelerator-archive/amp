FROM alpine
RUN apk --no-cache add ca-certificates
COPY amp-aws.alpine /usr/local/bin/amp-aws
ENTRYPOINT [ "amp-aws" ]
