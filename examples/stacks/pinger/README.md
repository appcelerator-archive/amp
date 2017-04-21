Pinger
======

Simple HTTP service (Go [source](https://github.com/subfuzion/docker-pinger)) that responds with "pong" when you GET /ping

Run in this directory:

    $ amp --server=localhost:50101 stack up -f stack.yml pinger

The app will be available at [https://pinger.local.atomiq.io/ping](https://pinger.local.atomiq.io/ping)

Test with

    $ curl -k https://pinger.local.atomiq.io/ping
    

