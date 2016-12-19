Instavote
=========

A voting application based on the canonical Docker swarm example.

Run in this directory:

    $ amp stack up -f stack.yml instavote

The voting app will be available at [http://vote.instavote.local.atomiq.io](http://vote.instavote.local.atomiq.io).

The results app will be available at [http://results.instavote.local.atomiq.io](http://results.instavote.local.atomiq.io).

Architecture
------------

![Architecture diagram](architecture.png)

* A Python webapp which lets you vote between two options
* A Redis queue which collects new votes
* A .NET worker which consumes votes and stores them in…
* A Postgres database backed by a Docker volume
* A Node.js webapp which shows the results of the voting in real time

Credit: Docker ([LICENSE](https://github.com/docker/example-voting-app/blob/master/LICENSE))

