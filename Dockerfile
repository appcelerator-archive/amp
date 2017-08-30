FROM appcelerator/amptools:1.6.0

COPY . /go/src/github.com/appcelerator/amp
WORKDIR /go/src/github.com/appcelerator/amp

CMD ["go", "test", "github.com/appcelerator/amp/tests/integration/..."]
