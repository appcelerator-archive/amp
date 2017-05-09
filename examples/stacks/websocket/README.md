Websocket Bash
==============

Serves up a websocket connected to a bash session

Run in this directory:

    $ ./deploy.sh

The app will be serving a websocket at [http://ws.websocket.local.appcelerator.io/](http://ws.websocket.local.appcelerator.io/)
and a web interface at [http://web.websocket.local.appcelerator.io/](http://web.websocket.local.appcelerator.io/)

Test using

    $ go run client.go

You can then check the logs with

    $ amp logs websocket

And stats with

    $ amp stats websocket

For extra fun try the following from within a bash session

    $ node client.js
