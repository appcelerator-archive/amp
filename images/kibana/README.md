# appcelerator/kibana

Docker Image for [Kibana](https://www.elastic.co/products/kibana).

Based on Alpine Linux (appcelerator/alpine).

The container will wait for the availability of Elasticsearch, import the index pattern ampbeat-\* and the save objects.

## Run

Most basic form:
```
docker run -t -p 5601:5601 -e "ELASTICSEARCH_URL=http://myElasticSearchHost:9200" appcelerator/kibana
```
With SSL:
```
docker run -t -p 443:443 -v /etc/kibana/ssl:/etc/kibana/ssl -e SERVER_SSL_CERTIFICATE=/etc/kibana/ssl/kibana.crt -e SERVER_SSL_KEY=/etc/kibana/ssl/kibana.key -e "ELASTICSEARCH_URL=http://myElasticSearchHost:9200" appcelerator/kibana
```

# Configuration (ENV, -e)

- `ELASTICSEARCH_URL`: URL of ElasticSearch. Default value: `http://elasticsearch:9200`
- `KIBANA_BASE_PATH`: Value of 'server.basePath' inside kibana.yml. Default value: *Empty*
- `SERVER_SSL_CERTIFICATE`: full container path for a TLS certificate
- `SERVER_SSL_KEY`: full container path for a ssl key
