# appcelerator/kibana

Docker Image for [Kibana](https://www.elastic.co/products/kibana).

Based on Alpine Linux (appcelerator/alpine).

## Run

Most basic form:
```
docker run -t -p 5601:5601 -e "ELASTICSEARCH_URL=http://myElasticSearchHost:9200" appcelerator/kibana
```

# Configuration (ENV, -e)

- `ELASTICSEARCH_URL`: URL of ElasticSearch. Default value: `http://elasticsearch:9200`
- `KIBANA_BASE_PATH`: Value of 'server.basePath' inside kibana.yml. Default value: *Empty*

## Tags

- 5.3.0, 5.3, latest
