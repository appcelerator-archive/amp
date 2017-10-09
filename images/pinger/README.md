# appcelerator/pinger

Simple service for demo purposes. Responds with `pong` to a `GET` request at: `http://<host>:<port>/ping`, and logs to `stdout`.

## Running pinger

    $ docker run -it --rm -p 3000:3000 --name pinger appcelerator/pinger

or

    $ docker service create -p 3000:3000 --name pinger appcelerator/pinger

## Scaling pinger

Try scaling scaling pinger up

    $ docker service scale pinger=3

Confirm all the replicas are running as the service scales to the desired amount

    $ docker service 
    $ docker service ls --filter name=pinger
    ID            NAME    REPLICAS  IMAGE                COMMAND
    bfjh7nervm3z  pinger  3/3       appcelerator/pinger

Then test

    $ curl localhost:3000/ping
    [f0afaf6c5ea5] pong

Each time you curl, you should see a different hostname in the brackets of the response
as Docker automatically load balances requests among all the service replicas.

### Note
Because HTTP/1.1 defaults to persistent connections, to assure the same result for
demo purposes in the browser, the `ping` service closes the HTTP connection with each response
by setting `Connection: close` in the response header.

## Development

`build.sh` - builds an alpine-based image: `appcelerator/pinger`

`run.sh` - convenience script for `docker run ...`

`service-start.sh` - convenience script for `docker service create ...`

`service-stop.sh` - convenience script for `docker service rm ...`
