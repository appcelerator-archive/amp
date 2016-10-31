Counter
======

Two replicated webapps connected to a redis counter

Run in this directory:

    $ amp stack up -f stack.yml counter

The app will be running at [http://go.counter.localhost.tv](http://go.counter.localhost.tv) and [http://python.counter.localhost.tv](http://python.counter.localhost.tv)

Test with

    $ curl http://go.counter.localhost.tv

Or open urls in browser
