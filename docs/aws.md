# AWS deployment

The deployment of an AMP cluster (single cluster mode), or of an additional cluster dedicated to an organization are performed in a similar way.
Both involve the creation of a bootstrap consisting of N swarm manager nodes (N=1, 3 or 5), each one running an Infrakit container. Infrakit follows the swarm election process, the leader container is always on the leader manager node.
Once this seed of N manager nodes is ready, a configuration is pushed to infrakit to build the worker nodes in the same swarm. From there, the swarm cluster is self managed.

## new management cluster or single-cluster AMP

```
+-------+                +--------+                  +-----+     +----------+   +-----+               +-----+
| Admin |                | cloned |                  | aws |     | infrakit |   | aws |               | aws |
| user  |                | repo   |                  | cli |     | cli      |   | API |               |     |
+-------+                +--------+                  +-----+     +----------+   +-----+               +-----+

 +----+  +- deploy-aws ->  +---+                      +--+
 |    |                    |   |   +-new container -> |  |                        ++
 |    |                    |   |                      +--+  +---- new stack -->   || +----+------------+ +-+
 +----+                    |   | +----------------------------------------------+ || | +- + elb -----> | | |
                           |   | | +- get status ---> +--+                      | || |                 | | |
                           |   | |                    |  |  +- get status ----> | || | +- + asg -----> | | |
                           |   | |     loop           |  |                      | || |                 | | |
                           |   | |                    |  |  +- in progress ---+ | || | +- + 3 ec2 ---> | | |
                           |   | | ^+ in progress +-+ +--+                      | || +----+------------+ | |
                           |   | +----------------------------------------------+ ||                     | |
                           |   |                      +--+                        ||                     | |
                           |   |                      |  |  <- create complete-+  ||                     | |
                           |   |   <- complete -----+ +--+                        ||                     | |
                           |   |                                                  ||                     | |
                           |   |   +- get url ------> +--+                        ||                     | |
                           |   |                      |  |  +- get url -------->  ||                     | |
                           |   |                      |  |                        ||                     | |
                           |   |                      |  |  <---- manager url -+  ||                     | |
                           |   |   <--- manager url + +--+                        ++                     | |
                           |   |                                    +---+                                | |
                           |   |   +-- commit worker group ----->   |   |                                | |
                           |   |                                    |   |   +- watch worker group --->   | |
                           |   |                                    +---+                                | |
                           |   | +---------------------------------------------------------------------+ | |
                           |   | | +--- get status +------------->  +---+                              | | |
                           |   | |                                  |   |   +- get status -----------> | | |
                           |   | |      loop                        |   |                              | | |
                           |   | |                                  |   |   <-not converged ---------+ | | |
                           |   | | ^----+ not con^erged +--------+  +---+                              | | |
                           |   | +---------------------------------------------------------------------+ | |
                           |   |                                    +---+                                | |
                           |   |                                    |   |   <------------- converged +   | |
                           |   |  <------------------- converged +  |   |                                | |
                           |   |                                    +---+                                | |
                           |   |                                                                         | |
                           |   |  +---------------- establish ssh tunnel to engine api -------------->   | |
                           |   |                                                                         | |
                           |   |  +----------------- push docker secrets ---------------------------->   | |
                           |   |                                                                         | |
                           +---+  +------------------deploy core stacks ----------------------------->   +-+
```

## new organization cluster

```
+-------+                +----------+               +-----+      +----------+                 +-----+
| amp   |                | amplifier|               | aws |      | infrakit |                 | aws |
| cli   |                |          |               | API |      | cli      |                 |     |
+-------+                +----------+               +-----+      +----------+                 +-----+

 +----+  +- create ----->  +---+
 |    |                    |   |                       ++                                         +-+
 |    |                    |   | +- new stack ----->   || +-----------------+------------+        | |
 |    |  <- cluster id -+  |   | +-+---------------+   || |  +------------- + elb -----> |        | |
 +----+                    |   | | - get status->  |   || |                              |        | |
                           |   | |                 |   || |  +------------- + asg -----> |        | |
                           |   | |    loop         |   || |                              |        | |
                           |   | |                 |   || |  +------------- + 3 ec2 ---> |        | |
                           |   | | <-in progress + |   || +-----------------+------------+        | |
                           |   | +---------------+-+   ||                                         | |
                           |   |                       ||                                         | |
                           |   |                       ||                                         | |
                           |   |   <+ complete +----+  ||                                         | |
                           |   |                       ||                                         | |
                           |   |   ++ get url +----->  ||                                         | |
                           |   |                       ||                                         | |
                           |   |   <--+ manager url +  ||                                         | |
                           |   |                       ++                                         | |
                           |   |                                                                  | |
                           |   |                                    +---+                         | |
                           |   |   +-+ commit worker group +---->   |   |                         | |
                           |   |                                    |   |   +- watch group ->     | |
                           |   |                                    +---+                         | |
                           |   | +------------------------------------------------------------+   | |
                           |   | | +--- get status -------------->  +---+                     |   | |
                           |   | |                                  |   |   +- get status +-> |   | |
                           |   | |      loop                        |   |                     |   | |
                           |   | |                                  |   |   <-not converged + |   | |
                           |   | | <----- not con^erged ---------+  +---+                     |   | |
                           |   | +------------------------------------------------------------+   | |
                           |   |                                    +---+                         | |
                           |   |                                    |   |   <-converged ----+     | |
                           |   |  <------------------- converged +  |   |                         | |
                           |   |                                    +---+                         +-+
                           +---+
```
