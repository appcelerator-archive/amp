Counter
=======

Two sample web applications connected to a redis counter. Both display the combined total number of page views and are scaled with 3 replicas.

Run in this directory:

    $ amp stack up -f stack.yml counter

The app will be running at [http://go.counter.local.atomiq.io](http://go.counter.local.atomiq.io) and [http://python.counter.local.atomiq.io](http://python.counter.local.atomiq.io)

Test with

    $ curl http://go.counter.local.atomiq.io
    $ curl http://python.counter.local.atomiq.io

Or open urls in browser.

