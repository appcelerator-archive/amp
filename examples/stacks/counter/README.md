Counter
=======

A sample web application connected to a redis counter. It displays the total number of page views and is scaled with 3 replicas.

Run in this directory:

    $ amp -s localhost stack deploy -c counter.yml

The app will be running at [http://go.counter.local.appcelerator.io](http://go.counter.local.appcelerator.io)

Test with

    $ curl http://go.counter.local.appcelerator.io

Or open url in browser.
