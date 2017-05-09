Javascript Microservice
=======================

Single function (serverless style) application.

Run in this directory:

    $ ./deploy.sh

The app will be running at [http://www.micro.local.appcelerator.io/](http://www.micro.local.appcelerator.io/).
It simply responds with a hello message on any route using any HTTP verb.

Example:

    $ curl www.micro.local.appcelerator.io
    [GET /] hello.

    $ curl www.micro.local.appcelerator.io/hi
    [GET /hi] hello.

    $ curl -X POST www.micro.local.appcelerator.io
    [POST /] hello.
