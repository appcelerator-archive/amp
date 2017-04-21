Pinger
======

Simple HTTP service (Go [source](https://github.com/subfuzion/docker-pinger)) that responds with "pong" when you GET /ping

Run in this directory:

    $ amp --server=localhost:50101 stack deploy -c pinger.yml pinger

The app will be available at [https://pinger.local.atomiq.io/ping](https://pinger.local.atomiq.io/ping)

Test with

    $ curl -k https://pinger.local.atomiq.io/ping
    [9a58616b614d] pong
    
    $ curl -k https://pinger.local.atomiq.io/ping
    [81c3c0958600] pong

    $ curl -k https://pinger.local.atomiq.io/ping
    [ba7d86288a02] pong

Because 3 replicas was specified, you should see 3 different hosts in
the brackets as you repeat the curl command (or refresh your browser page).
    

