Pinger
======

Simple HTTP service (Go [source](https://github.com/subfuzion/docker-pinger)) that responds with "pong" when you GET /ping

Run in this directory:

    $ amp -s localhost stack deploy -c pinger.yml

The app will be available at [http://pinger.examples.local.appcelerator.io/ping](http://pinger.examples.local.appcelerator.io/ping)

Test with

    $ curl http://pinger.examples.local.appcelerator.io/ping
    [9a58616b614d] pong

    $ curl http://pinger.examples.local.appcelerator.io/ping
    [81c3c0958600] pong

    $ curl http://pinger.examples.local.appcelerator.io/ping
    [ba7d86288a02] pong

Because 3 replicas was specified, you should see 3 different hosts in
the brackets as you repeat the curl command (or refresh your browser page).
