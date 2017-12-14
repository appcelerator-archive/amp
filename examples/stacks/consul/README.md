# Consul cluster for Docker swarm

Consul can be deployed on swarm in two different manners: a single service with replicas, or multiple single-replica services.

The Multiple services approach has a few advantages:
 - it allows dedicated volumes (think non local volumes), which is better for HA,
 - does not require the dnsrr mode (no more VIPs).

## Single service

On an existing amp cluster, deploy the consul stack:

    $ amp stack deploy -c consul.single-service.yml examples
    [user ndegory @ REDACTED.us-west-2.elb.amazonaws.com:50101]
    Deploying stack examples using consul.yml
    Creating network examples_db
    Creating service examples_consul

It will create 3 replicas that will join a single cluster.

You can check that all replicas are up:

    $ amp service ps examples_consul
    [user ndegory @ REDACTED.us-west-2.elb.amazonaws.com:50101]
    ID                          NAME                IMAGE          DESIRED STATE   CURRENT STATE   NODE ID                     ERROR
    j7ry4rki29926uu9dhdbkeznj   examples_consul.1   consul:1.0.1   RUNNING         RUNNING         sujeodnpbcqkofm8plcl9ctes
    h9j2kj56lg0ywl5yokgr90ezx   examples_consul.2   consul:1.0.1   RUNNING         RUNNING         c7374cfesr6v7igc9fnuo3n8r
    laij1kcwe8u4h14j6j1ogqi6i   examples_consul.3   consul:1.0.1   RUNNING         RUNNING         m2y590q671tp5yfpip7k3o5xr

The UI is available on https://consul.examples.DOMAIN_NAME/, and you can also list the members of the Consul cluster, you should get a list of 3 task IPs:

    $ curl -k -H Host:consul.examples.local.appcelerator.io https://$REMOTE_SERVER/v1/status/peers
    ["10.0.3.4:8300","10.0.3.2:8300","10.0.3.3:8300"]

In case something goes wrong on a replica (container breaks or node breaks), a new replica will join the cluster. This example doesn't cleanup the failed nodes, so this is a procedure that should be done by another way:
To illustrate this, here's the member list on a cluster where a replica has been killed:

    $ docker exec -t examples_consul.1.j7ry4rki29926uu9dhdbkeznj consul members
    Node          Address        Status  Type    Build  Protocol  DC   Segment
    342de01bc01a  10.0.3.4:8301  alive   server  1.0.1  2         dc1  <all>
    db254fbc9cc8  10.0.3.4:8301  failed  server  1.0.1  2         dc1  <all>
    f35acf5b3bab  10.0.3.3:8301  alive   server  1.0.1  2         dc1  <all>
    f46c8aac2ec8  10.0.3.2:8301  alive   server  1.0.1  2         dc1  <all>

    $ docker exec -t examples_consul.1.j7ry4rki29926uu9dhdbkeznj consul force-leave db254fbc9cc8
    $ docker exec -t examples_consul.1.j7ry4rki29926uu9dhdbkeznj consul members
    Node          Address        Status  Type    Build  Protocol  DC   Segment
    342de01bc01a  10.0.3.4:8301  alive   server  1.0.1  2         dc1  <all>
    db254fbc9cc8  10.0.3.4:8301  left    server  1.0.1  2         dc1  <all>
    f35acf5b3bab  10.0.3.3:8301  alive   server  1.0.1  2         dc1  <all>
    f46c8aac2ec8  10.0.3.2:8301  alive   server  1.0.1  2         dc1  <all>

## Multiple services

Same procedure as above, but instead use the other example stack file.

    $ amp stack deploy -c consul.multiple-services.yml examples
