# Changelog

## 0.5.0 (2017-01-03)

### Client

* Add AMP cluster management features [#519](https://github.com/appcelerator/amp/pull/519)
* Add Serverless features [#539](https://github.com/appcelerator/amp/pull/539)
* Enhanced `amp version` outputs [#612](https://github.com/appcelerator/amp/pull/612)
* Client refactoring [#613](https://github.com/appcelerator/amp/pull/613), [#614](https://github.com/appcelerator/amp/pull/614), [#615](https://github.com/appcelerator/amp/pull/615), [#616](https://github.com/appcelerator/amp/pull/616)
* Update config loading precedence [#617](https://github.com/appcelerator/amp/pull/617)
* Fix default value for AdminServerAddress [#627](https://github.com/appcelerator/amp/pull/627)
* Prevent amp logs from crashing when Verbose flag is set [#630](https://github.com/appcelerator/amp/pull/630)

### Documentation

* Deployment documentation fix [#609](https://github.com/appcelerator/amp/pull/609)
* Add 0.5.0 changelog [#628](https://github.com/appcelerator/amp/pull/628)

###  Platform

* Add AMP Bootstrap features [#528](https://github.com/appcelerator/amp/pull/528)
* Makefile enhancements [#622](https://github.com/appcelerator/amp/pull/622)
* Add `make rules` to print all the Makefile rules [#624](https://github.com/appcelerator/amp/pull/624)
* Fix Make check warnings [#625](https://github.com/appcelerator/amp/pull/625)
* Update haproxy to version 1.0.3 [#626](https://github.com/appcelerator/amp/pull/626)
* Update elasticsearch to version 5.1.1 [#619](https://github.com/appcelerator/amp/pull/619)
* Add GRPC rest gateway for browser-accessible API [#488](https://github.com/appcelerator/amp/pull/488)

### Examples

* Add GRPC rest gateway based UI [#527](https://github.com/appcelerator/amp/pull/527)
* Add `functions` to the UI [#605](https://github.com/appcelerator/amp/pull/605)
* Add `kv` to the UI [#606](https://github.com/appcelerator/amp/pull/606)

## 0.4.0 (2016-12-20)

### Client

* Configuration file moved to $HOME/.config/amp/amp.yaml [#592](https://github.com/appcelerator/amp/pull/592)
* Fix registry commands [#536](https://github.com/appcelerator/amp/pull/536)
* add amp platform command [#455](https://github.com/appcelerator/amp/pull/455)
* AMP version option [#462](https://github.com/appcelerator/amp/pull/462), [#589](https://github.com/appcelerator/amp/pull/589)
* AMP kv commands [#511](https://github.com/appcelerator/amp/pull/511)

###  Platform

* Grafana 4.0 [#518](https://github.com/appcelerator/amp/pull/518)
* Influxdata stack version 1.1 [#443](https://github.com/appcelerator/amp/pull/443), [#474](https://github.com/appcelerator/amp/pull/474)
* Robustness [#442](https://github.com/appcelerator/amp/pull/442), [#443](https://github.com/appcelerator/amp/pull/443), [#476](https://github.com/appcelerator/amp/pull/476)
* grpc code consistency [#586](https://github.com/appcelerator/amp/pull/586) [#585](https://github.com/appcelerator/amp/pull/585) [#584](https://github.com/appcelerator/amp/pull/584) [#583](https://github.com/appcelerator/amp/pull/583)
* Fix amp stats panic [#533](https://github.com/appcelerator/amp/pull/533) [#534](https://github.com/appcelerator/amp/pull/534)
* add cli color theme [#495](https://github.com/appcelerator/amp/pull/495)
* Fix stats issues [#404](https://github.com/appcelerator/amp/pull/404) [#403](https://github.com/appcelerator/amp/pull/403)

### Vendoring

* Global update of vendors, based on tags instead of commits [#501](https://github.com/appcelerator/amp/pull/501), [#515](https://github.com/appcelerator/amp/pull/515)

### Documentation

* Documentation Update [#601](https://github.com/appcelerator/amp/pull/601)
* amp platform commands doc [#473](https://github.com/appcelerator/amp/pull/473)
* add stacks user guide [#391](https://github.com/appcelerator/amp/pull/391)

### Tests

* Dockerized integration tests [#479](https://github.com/appcelerator/amp/pull/479)
* refactor rpc tests [#416](https://github.com/appcelerator/amp/pull/416)
* improve stats tests [#405](https://github.com/appcelerator/amp/pull/405)
* CLI tests
  * Synchronous setup and tearDown [#522](https://github.com/appcelerator/amp/pull/522)
  * Asynchronous execution [#454](https://github.com/appcelerator/amp/pull/454)
  * Templating regexes [#478](https://github.com/appcelerator/amp/pull/478)
  * Templating ports [#441](https://github.com/appcelerator/amp/pull/441) [#445](https://github.com/appcelerator/amp/pull/445)
* Updated Timeout and delay values [#451](https://github.com/appcelerator/amp/pull/451)
* AMP kv commands [#511](https://github.com/appcelerator/amp/pull/511)
* Improved Regex [#544](https://github.com/appcelerator/amp/pull/544)

## 0.3.0 (2016-11-08)

### Client

* AMP config get and set [#420](https://github.com/appcelerator/amp/pull/420)

### Platform

* InfluxData and Grafana images are now specialized for AMP [#394](https://github.com/appcelerator/amp/pull/394)

### Tests

* Refactored CLI tests [#386](https://github.com/appcelerator/amp/pull/386), [#398](https://github.com/appcelerator/amp/pull/398)
* Regexp for tests [#383](https://github.com/appcelerator/amp/pull/383), [#387](https://github.com/appcelerator/amp/pull/387), [#397](https://github.com/appcelerator/amp/pull/397)
* Timeout, retry and delay in CLI tests [#421](https://github.com/appcelerator/amp/pull/421)
* Test coverage [#396](https://github.com/appcelerator/amp/pull/396)

### Documentation

* Stack documentation [#390](https://github.com/appcelerator/amp/pull/390), [#391](https://github.com/appcelerator/amp/pull/391)

## 0.2.2 (2016-10-25)

### Platform

* Hotfix on v0.2.1 (locked versions of Docker images for the AMP swarm)


## 0.2.1 (2016-10-24)

https://github.com/appcelerator/amp/milestone/5?closed=1

### Client

* Fix update amp stats [#357](https://github.com/appcelerator/amp/pull/357)
* Stack create command [#362](https://github.com/appcelerator/amp/pull/362)
* Add amp topic commands [#366](https://github.com/appcelerator/amp/pull/366)
* Add more CLI tests [#368](https://github.com/appcelerator/amp/pull/368)
* Add voting app as example [#379](https://github.com/appcelerator/amp/pull/379)
* Add regexp for logs [#383](https://github.com/appcelerator/amp/pull/383)

### Platform

* Fix stack workaround upon docker 1.12.2 [#345](https://github.com/appcelerator/amp/pull/345)
* Add stack volume support [#355](https://github.com/appcelerator/amp/pull/355)
* Fix etcd version mismatch [#365](https://github.com/appcelerator/amp/pull/365)
* Remove reference to Kafka [#367](https://github.com/appcelerator/amp/pull/367)
* Swith to etcd 3.1 [#380](https://github.com/appcelerator/amp/pull/380)
* Slack output for Kapacitor based on environment variables [#347](https://github.com/appcelerator/amp/pull/347)
* Force swarm init on 127.0.0.1 [#349](https://github.com/appcelerator/amp/pull/349)
* Fix make test [#352](https://github.com/appcelerator/amp/pull/352)
* Update parse & networks tests [#356](https://github.com/appcelerator/amp/pull/356)
* Reduce Kapacitor alerts noise [#358](https://github.com/appcelerator/amp/pull/358)
* Use latest release of NATS streaming [#372](https://github.com/appcelerator/amp/pull/372)


## 0.2.0 (2016-10-20)

https://github.com/appcelerator/amp/milestone/3?closed=1

### Client

* Enrich CLI command tests [#305](https://github.com/appcelerator/amp/pull/305)
* Fix stack tests [#333](https://github.com/appcelerator/amp/pull/333)
* Fix Stack listing quite mode [#326](https://github.com/appcelerator/amp/pull/327)
* Add Stack listing options [#331](https://github.com/appcelerator/amp/pull/331)
* Improve error messages, verbose mode & consistant ids [#332](https://github.com/appcelerator/amp/pull/332)
* Fix log commands [#340](https://github.com/appcelerator/amp/pull/340)

### Documentation

* Documentation update [#300](https://github.com/appcelerator/amp/pull/300)

### Platform

* Add Network item for services in stack [#304](https://github.com/appcelerator/amp/pull/304)
* Use external network in stacks [#341](https://github.com/appcelerator/amp/pull/341)
* Add volumes/mount in stack file [#299](https://github.com/appcelerator/amp/issues/299)
* Replaces Kafka and zookeeper messaging by NATS [#325](https://github.com/appcelerator/amp/pull/325)
* Refactored state machine to use strings instead of integers [#335](https://github.com/appcelerator/amp/pull/335)
* Remove volumes after removing the services [#297](https://github.com/appcelerator/amp/pull/297)
* Fix Swarm monitor refresh [#308](https://github.com/appcelerator/amp/pull/308)
* ETCD - switching to stable branch 3.0 [#312](https://github.com/appcelerator/amp/pull/312)
* Add telegraf service and grafana dashboard for haproxy stats [#315](https://github.com/appcelerator/amp/pull/315)
* Use nats streaming as messaging system [#318](https://github.com/appcelerator/amp/issues/318)
* Fix explicit versions for images used by amp [#321](https://github.com/appcelerator/amp/pull/321)
* Faster docker build of amp image [#334](https://github.com/appcelerator/amp/pull/334)
* Moving amp-agent & amp-log-worker to amp [#337](https://github.com/appcelerator/amp/pull/337)

### Vendoring

* Updating vendor for nats streaming [#319](https://github.com/appcelerator/amp/pull/319)


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

### Platform

* AMP Swarm stop also removes user services [#234](https://github.com/appcelerator/amp/pull/234)
* Etcd ListRaw & Watch feature addition [#213](https://github.com/appcelerator/amp/pull/213)
* Fix HAproxy [#231](https://github.com/appcelerator/amp/pull/231)
* Add stack-id and stack-name as labels in containers [#237](https://github.com/appcelerator/amp/pull/237)
* Add service labels [#249](https://github.com/appcelerator/amp/pull/249)
* Add env support to service [#244](https://github.com/appcelerator/amp/pull/244)
* Add replicated/global mode support for Service [#256](https://github.com/appcelerator/amp/pull/256)
* Fix `~/registry/data` automated creation [#285](https://github.com/appcelerator/amp/pull/285)
* Networking basis enhancement - all services attached to amp-public by default [#204](https://github.com/appcelerator/amp/pull/204)
* Add Service network attachment [#266](https://github.com/appcelerator/amp/pull/266)

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

### Platform

* Swarm start/stop options [#16](https://github.com/appcelerator/amp/pull/16)
* Swarm monitor option [#33](https://github.com/appcelerator/amp/pull/33)
* HAproxy integration [#100](https://github.com/appcelerator/amp/pull/100)
* Etcd integration [#2](https://github.com/appcelerator/amp/pull/2)
* Kafka integration [#52](https://github.com/appcelerator/amp/pull/52)
* Registry integration [#155](https://github.com/appcelerator/amp/pull/155)
* ElasticSearch integrationo [#4](https://github.com/appcelerator/amp/pull/4)
* Telegraf/InfluxDB/Grafana/Kibana - TIGK stack integration for observability [#74](https://github.com/appcelerator/amp/pull/74)
* Zookeeper integration [#294](https://github.com/appcelerator/amp/pull/294)

### Vendoring

* Glide update [#37](https://github.com/appcelerator/amp/pull/37)
* Add fixed dependencies support [#47](https://github.com/appcelerator/amp/pull/47)
