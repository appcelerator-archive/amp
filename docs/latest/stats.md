
### Stats

The `amp stats` command is used to query or stream statistics. It provides group, filtering and historic options to manage what is presented.

    $ amp stats --help

    Usage:	amp stats SERVICE

    Compute statistics for the service name SERVICE
    Without SERVICE name, display the statistics of all available services
    Options:

      type of statistics:
          --cpu                     display cpu statistics
          --mem                     display memory statistics
          --net                     display net tx/rx statistics
          --io                      display disk io statistics
          if no one of the options above is set then display all statistics, cpu, meme, net, io

      statistics group:
          --container               present a list all available containers with their statistics
          --node                    present a list all available nodes with their statistics
          --service                 present a list all available services with their statistics
          --stack                   present a list all available stacks with their statistics
          default is service

      filters:
           --container-id string     select only the container having their id starting by the given value
           --container-name string   select only the container having their name starting by the given value
           --container-state string  select only the container having their state starting by the given value
           --service-id string       select only the services having their id starting by the given value
           --stack-name string       select only the stacks having their name starting by the given value
           --task-id string          select only the task having their id starting by the given value
           --node-id string          select only the node having their name starting by the given value
           if several filters are used they are considered linked by a logical AND

      period of computing:
          --period string            define the period of time inside which statistics are computed     

          the given value should have the following format: now-xt
          where:
            - x is a number
            - t is a time unit amount: y=year, M=month, w=week, d=day, h=hour, m=minute, s=second
          for instance, now-10h means that the period starts from 10 hours in the past until now
          By default, the period is "now-10s"

      historic statistics:
          historic statistics display statistics by time period,
          where current statistic displays the last available statistics in the period one shot.

          --time-group string       define the period of time use to group statistics

          the given value should have the following format: xt
          where:
            - x is a number
            - t is a time unit amount: y=year, M=month, w=week, d=day, h=hour, m=minute, s=second
          for instance, 10s means 10 seconds, 3d means 3 days, ...
          there is no default for time-group and no default for period if time-group is set
          then for historic stats time-group and period should be set all together

     -other options
         -f --follow                if set the statistics are continually computed

         Current statistics are recomputed every 2 seconds
         Historic statistics add new lines depending of the time group.            


A few useful examples:

* To compute and follow stats by services:
```
  $ amp stats -f
```

* To compute and follow the stats for a specific service:
```
  $ amp stats -f etcd
```

* To compute only cpu and mem statistics for all stakcs
```
  $ amp stats --stack --cpu --mem
```

* To compute stats by node
```
  $ amp stats --node
```

* To compute stats for a stack on a specific node
```
  $ amp stats --stack-name monitoring --node-id aNodeId
```

* To compute historic stats for a service
```
  $ amp stats --period now-1m --time-group=10s elasticsearch
```

* To compute historic stats for a stack and follow to see a new line every 3 seconds
```
  $ amp stats --stack-name monitoring --period now-30s --time-group=3s -f
```
