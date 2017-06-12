FROM appcelerator/amptools

COPY . /go/src/github.com/appcelerator/amp
WORKDIR /go/src/github.com/appcelerator/amp

CMD ["go", "test", "github.com/appcelerator/amp/tests/integration/..."]
