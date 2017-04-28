# Elasticsearch

[Elasticsearch](https://www.elastic.co/products/elasticsearch) Docker image based on Alpine Linux.

### Exposed ports

- `9200`, `9300`


### Env. variables

Variable | Description | Default value | Example
 ------- | ----------- | ------------- | -------
JAVA_HEAP_SIZE | Java heap size in MB | | 1024
java_max_direct_mem_size | Java max direct memory size | |
java_opts | Java options | |
NETWORK_HOST | IP on which ES will be listening | 0.0.0.0 | \_eth0\_
UNICAST_HOSTS | DNS names for cluster discovery | | tasks.elasticsearch
MIN_MASTER_NODES | Minimum number of master nodes (quorum) | 1 | 2
MEMORY_LOCK | Bootstrap memory lock | false | true

When JAVA_HEAP_SIZE is empty, the value is set depending on system max memory (256 to 10% of max memory).

### System prerequisites

Elastic highly recommends to set the VM mmap count to 262144 on the host: https://www.elastic.co/guide/en/elasticsearch/reference/5.0/vm-max-map-count.html

    sudo sysctl -w vm.max_map_count=262144

The hard limit for file descriptors should be at least 65535. You can check it with `ulimit -Hn`.
