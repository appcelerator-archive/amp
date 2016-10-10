# Changelog

## 0.1.1 (2016-10-10)

### Client

* Add Stack list support + Stack management by Id & Name [#206](https://github.com/appcelerator/amp/pull/206)
* Add Stack restart by Id & Name command [#209](https://github.com/appcelerator/amp/pull/209)
* Client template standardization [#210](https://github.com/appcelerator/amp/issues/210)
* Add `amp stats [serviceName/serviceId]` support [#217](https://github.com/appcelerator/amp/pull/217)
* Add `amp logs [serviceName]` support [#222](https://github.com/appcelerator/amp/pull/222)
* Add stack state to `amp stack ls` [#223](https://github.com/appcelerator/amp/pull/223)
* Fix stack rollback [#226](https://github.com/appcelerator/amp/pull/226)
* Add quiet mode for stack ls `amp stack ls -q` [#230](https://github.com/appcelerator/amp/pull/230)
* Add Logs by stack [#242](https://github.com/appcelerator/amp/pull/242)
* Add service rm support [#253](https://github.com/appcelerator/amp/pull/253)
* Add registry ls and auto tag [#292](https://github.com/appcelerator/amp/pull/292)

### Documentation

* Documentation Update [#207](https://github.com/appcelerator/amp/pull/207)

### Runtime

* Etcd ListRaw & Watch feature addition [#213](https://github.com/appcelerator/amp/pull/213)
* Fix HAproxy [#231](https://github.com/appcelerator/amp/pull/231)
* Add stack-id and stack-name as labels in containers [#237](https://github.com/appcelerator/amp/pull/237)
* Add service labels [#249](https://github.com/appcelerator/amp/pull/249)
* Add env support to service [#244](https://github.com/appcelerator/amp/pull/244)
* Add replicated/global mode support for Service [#256](https://github.com/appcelerator/amp/pull/256)
* Fix `~/registry/data` automated creation [#285](https://github.com/appcelerator/amp/pull/285)

### Networking

* Networking basis enhancement - all services attached to amp-public by default [#204](https://github.com/appcelerator/amp/pull/204)
* Add Service network attachment [#266](https://github.com/appcelerator/amp/pull/266)

### Swarm

* AMP Swarm stop also removes user services [#234](https://github.com/appcelerator/amp/pull/234)

### Vendoring

* Fix broken glide install #296 (https://github.com/appcelerator/amp/pull/296)


## 0.1.0 (2016-09-23)

Alpha release (limited Preview)

### Build

* Add a shrink script to reduce image size [#156](https://github.com/appcelerator/amp/pull/156)

### Client

* Add Log support [#11](https://github.com/appcelerator/amp/issues/11)
  * Log Streaming support [#66](https://github.com/appcelerator/amp/pull/66)
  * Log Filtering support [#67](https://github.com/appcelerator/amp/pull/67)
* Add Stats support [#68](https://github.com/appcelerator/amp/pull/68)
  * Streaming & Filtering support [#89](https://github.com/appcelerator/amp/pull/89)
* Add Registry Management features
  * Push images [#155](https://github.com/appcelerator/amp/pull/155)
* Add Stack support [#160](https://github.com/appcelerator/amp/pull/160)
  * Yaml parser support [#163](https://github.com/appcelerator/amp/pull/163)
  * Stop/Remove stack [#196](https://github.com/appcelerator/amp/pull/196)
* Add Service support [#177](https://github.com/appcelerator/amp/pull/177)
  * Service publication options [#197](https://github.com/appcelerator/amp/pull/197)
* Helpers integration [#200](https://github.com/appcelerator/amp/pull/200)

### Runtime

* HAproxy integration [#100](https://github.com/appcelerator/amp/pull/100)
* Etcd integration [#2](https://github.com/appcelerator/amp/pull/2)
* Kafka integration [#52](https://github.com/appcelerator/amp/pull/52)
* Registry integration [#155](https://github.com/appcelerator/amp/pull/155)
* ElasticSearch integrationo [#4](https://github.com/appcelerator/amp/pull/4)
* Telegraf/InfluxDB/Grafana/Kibana - TIGK stack integration for observability [#74](https://github.com/appcelerator/amp/pull/74)
* Zookeeper integration [#294](https://github.com/appcelerator/amp/pull/294)

### Swarm

* Swarm start/stop options [#16](https://github.com/appcelerator/amp/pull/16)
* Swarm monitor option [#33](https://github.com/appcelerator/amp/pull/33)

### Vendoring

* Glide update [#37](https://github.com/appcelerator/amp/pull/37)
* Add fixed dependencies support [#47](https://github.com/appcelerator/amp/pull/47)
