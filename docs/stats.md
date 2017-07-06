## Stats Command

The `amp stats` command is used to query or stream statistics. It provides group, filtering and historic options to manage what is presented.

### Usage

```
$ amp stats --help

Usage:	amp stats [OPTIONS] SERVICE

Display statistics

Options:
      --container                Display stats on containers
      --container-id string      Filter on container id
      --container-name string    Filter on container name
      --container-state string   Filter on container state
      --cpu                      Display cpu stats
  -f, --follow                   Follow stats output
  -h, --help                     Print usage
  -i, --include                  Include AMP stats
  -k, --insecure                 Control whether amp verifies the server's certificate chain and host name
      --io                       Display disk io stats
      --mem                      Display memory stats
      --net                      Display net rx/tx stats
      --node                     Display stats on nodes
      --node-id string           Filter on node id
      --period string            Historic period of metrics extraction, for instance: "now-1d", "now-10h", with y=year, M=month, w=week, d=day, h=hour, m=minute, s=second (default "now-10m")
  -s, --server string            Specify server (host:port)
      --service                  Display stats on services (default true)
      --service-id string        Filter on service id
      --stack                    Display stats on stacks
      --stack-name string        Filter on stack name
      --time-group string        Historic extraction by time group, for instance: "1d", "3h", , with y=year, M=month, w=week, d=day, h=hour, m=minute, s=second      
```

### Examples

* To compute and follow stats by services:
```
$ amp stats -f
```

* To compute and follow the stats for a specific service:
```
$ amp stats -f etcd
```

* To compute only cpu and mem statistics for all stacks:
```
$ amp stats --stack --cpu --mem
```

* To compute stats by node:
```
$ amp stats --node
```

* To compute stats for a stack on a specific node:
```
$ amp stats --stack-name monitoring --node-id qckrksy19cmh...
```

* To compute historic stats for a service:
```
$ amp stats --period now-1m --time-group=10s elasticsearch
```

* To compute historic stats for a stack and follow to see a new line every 3 seconds:
```
$ amp stats --stack-name monitoring --period now-30s --time-group=3s -f
```
