package stack

import (
	"testing"
)

var examples = []string{
	`web:
  image: appcelerator.io/amp-demo
  public:
    - name: www
      protocol: tcp
      publish_port: 80
      internalPort: 3000
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
		out, err := parseStackYaml(example)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(out)
	}
}
