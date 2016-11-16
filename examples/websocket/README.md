Websocket Bash
==============

Serves up a websocket connected to a bash session

Run in this directory:

    $ ./deploy.sh

The app will be serving a websocket at [http://ws.websocket.local.atomiq.io/](http://ws.websocket.local.atomiq.io/)
and a web interface at [http://web.websocket.local.atomiq.io/](http://web.websocket.local.atomiq.io/)

Test using

    $ go run client.go

You can then check the logs with

    $ amp logs websocket

And stats with

    $ amp stats websocket

For extra fun try the following from within a bash session

    $ node client.js
