Pinger
======

Simple HTTP service (Go [source](https://github.com/subfuzion/docker-pinger)) that responds with "pong" when you GET /ping

Run in this directory:

    $ amp stack up -f stack.yml pinger

The app will be available at [http://www.pinger.local.atomiq.io/ping](http://www.pinger.local.atomiq.io/ping)

Test with

    $ curl http://www.pinger.local.atomiq.io/ping
    

