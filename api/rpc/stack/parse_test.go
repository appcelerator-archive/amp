package stack

import (
	"golang.org/x/net/context"
	"testing"
)

var examples = []string{
	`web:
  image: appcelerator.io/amp-demo
  public:
    - name: www
      protocol: tcp
      publish_port: 90
      internal_port: 3000
  replicas: 3
  environment:
    REDIS_PASSWORD: password
redis:
  image: redis
  environment:
    - PASSWORD=password`,
}

func TestParseStackYaml(t *testing.T) {
	for _, example := range examples {
		out, err := NewStackfromYaml(context.Background(), example)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Print("Out: ", out)
		t.Log(out)
	}
}
