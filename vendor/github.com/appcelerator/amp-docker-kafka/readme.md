Kafka image.

### Maintainers

To trigger a build, run:

    curl -X POST https://registry.hub.docker.com/u/appcelerator/kafka/trigger/cc9bca73-f7b3-48db-9c2d-754bf430201b/

Finally, verify that the image was built successfully on the [Build Details page](https://hub.docker.com/r/appcelerator/kafka/builds/).

### Tags

- `0.9`, `0.9.0`, `0.9.0.1`
- `0.10`, `0.10.0`, `0.10.0.1`, `latest`

### Exposed ports

- `9092`


### Env. variables

  - ZOOKEEPER_CONNECT: Zookeeper nodes connection string, default localhost:2181
  - TOPIC_LIST: list of the topics needed to be created at Kafka statup format: "name1 name2 name3 ..."

### sample with Docker compose: zookeeper, kafka

    Kafka UI available at:    http://localhost


    version: '2'
    services:
      zookeeper:
        image: appcelerator/zookeeper:latest
        ports:
         - "2181:2181"
         - "2888:2888"
         - "3888:3888"
      kafka:
        image: appcelerator/kafka:latest
        ports:
         - "9092"
        environment:
         - ZOOKEEPER_CONNECT=zookeeper:2181
        depends_on:
         - zookeeper
      kafka-manager:
        image: sheepkiller/kafka-manager
        ports:
         - "80:9000"
        environment:
         - ZK_HOSTS=zookeeper:2181
        depends_on:
         - zookeeper


### sample with 3 Kafka nodes and 3 zookeeper nodes: zookeeper1, zookeeper2, zookeeper3, without containerPilot


    $ZOOKZEEPER_CONNECT="zookeper1:2181,zookeeper2:2181,zookeeper3:2181"


    docker run -d --name=kafka1 \
      -p 9092:9092 \
      -e "ZOOKEEPER_CONNECT=$ZOOKEEPER_CONNECT" \
      -e "BROKER_ID=1" \
      appcelerator/kafka:latest


    docker run -d --name=kafka2 \
      -p 9092:9092 \
      -e "ZOOKEEPER_CONNECT=$ZOOKEEPER_CONNECT" \
      -e "BROKER_ID=2" \
      appcelerator/kafka:latest


    docker run -d --name=kafka3 \
      -p 9092:9092 \
      -e "ZOOKEEPER_CONNECT=$ZOOKEEPER_CONNECT" \
      -e "BROKER_ID=3" \
      appcelerator/kafka:latest
