Cassandra stack file
====================

## Cassandra on cattle nodes

`amp stack deploy ./cassandra.cattle.yml` will deploy a stack named `cassandra` with 3 replicas.
This is a custom image that adds an initialization of the Cassandra configuration for the seeds and the listening address compatible with Docker Swarm.

If deployed on a single node, the volume mount should be disabled (it shouldn't be shared by different containers).

## Cassandra on pet nodes

`amp stack deploy ./cassandra.pets.yml` will deploy a stack named `cassandra with 3 one-task services.

It's based on the official Cassandra image, and declare all 3 tasks to be seeds, they're all equal.

It uses the rexray/ebs volume driver to persist the data on an EBS instance.
