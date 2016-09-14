package stack

import (
	"testing"
)

var examples = []string{
	`web:
  image: appcelerator.io/amp-demo
  ports:
    - 80:3000
  replicas: 3
  environment:
    REDIS_PASSWORD: password
redis:
  image: redis
  environment:
    - PASSWORD=password`,
	`tutum-cron:
  image: sillelien/tutum-cron
  environment:
    BACKUP_HOURLY_CRON_SCHEDULE: '0 * * * *'
    BACKUP_DAILY_CRON_SCHEDULE: '0 3 * * *'`,
	`influxdbData:
  image: busybox

influxdb:
  image: tutum/influxdb:0.9
  environment:
    - PRE_CREATE_DB=cadvisor
  ports:
    - "8083:8083"
    - "8086:8086"`,
	`lb:
  image: 'tutum/haproxy:latest'
  ports:
    - '80:80'

web-green:
  image: 'borja/bluegreen:v1'`,
	`btsync:
  image: "tutum/btsync"
  replicas: 3`,
	`lb:
  image: 'tutum/haproxy:0.2'
  ports:
    - '80:80'
mysql:
  image: 'mysql:5.6'
  environment:
    - MYSQL_ROOT_PASSWORD=**USE_A_STRONG_PASSWORD**
varnish:
  image: 'benhall/docker-varnish:latest'
  environment:
    - VARNISH_BACKEND_HOST=backend
    - VARNISH_BACKEND_PORT=80
    - 'VIRTUAL_HOST=example.com,www.example.com'
  ports:
    - '8080:80'
wordpress:
  image: 'wordpress:4.3-apache'
  environment:
    - WORDPRESS_DB_NAME=wp_example
    - WORDPRESS_TABLE_PREFIX=wp_`,
	`lb:
  image: dockercloud/haproxy
  ports:
    - "80:80"
web:
  image: dockercloud/quickstart-python
  replicas: 4
redis:
  image: redis`,
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
