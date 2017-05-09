Counter
=======

Two sample web applications connected to a redis counter. Both display the combined total number of page views and are scaled with 3 replicas.

Run in this directory:

    $ amp --server=localhost:50101 stack deploy -c stack.yml counter

The app will be running at [http://go.counter.local.appcelerator.io](http://go.counter.local.appcelerator.io) and [http://python.counter.local.appcelerator.io](http://python.counter.local.appcelerator.io)

Test with

    $ curl http://go.counter.local.appcelerator.io
    $ curl http://python.counter.local.appcelerator.io

Or open urls in ^browser^.
