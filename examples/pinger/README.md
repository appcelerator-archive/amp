Pinger
======

Simple HTTP service (Go [source](https://github.com/appcelerator/docker-pinger)) that responds with "pong" when you GET /ping

Run in this directory:

    $ amp stack up -f stack.yml pinger

Depending on your configuration, you may need to explicitly specify the server address as shown below:

    $ amp stack up -f stack.yml pinger --server localhost:8080

The app will be available at [http://www.pinger.local.atomiq.io](http://www.pinger.local.atomiq.io)

Test with

    $ curl http://www.pinger.local.atomiq.io/ping
    

