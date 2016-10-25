Instavote
=========

Getting started
---------------

Download [Docker for Mac or Windows](https://www.docker.com).

Run in this directory:

    $ amp stack up -f stack.yml instavote

The app will be running at [.](http://vote.instavote.localhost.tv), and the results will be at [http://results.instavote.localhost.tv](http://results.instavote.localhost.tv).

Architecture
-----

![Architecture diagram](architecture.png)

* A Python webapp which lets you vote between two options
* A Redis queue which collects new votes
* A .NET worker which consumes votes and stores them in…
* A Postgres database backed by a Docker volume
* A Node.js webapp which shows the results of the voting in real time

