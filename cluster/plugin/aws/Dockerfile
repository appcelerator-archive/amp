FROM alpine
RUN apk add --no-cache ca-certificates
COPY aws.alpine /usr/local/bin/aws
ENTRYPOINT [ "aws" ]
