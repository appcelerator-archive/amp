# Stack file reference

A stack is a collection of services that run in a specific environment on AMP under
[Docker swarm mode](https://docs.docker.com/engine/swarm/). 

A stack file is used to define this environment in [YAML](http://yaml.org/) format. 
The amp stack file is full compatible with docker compose file
 [Compose file](https://docs.docker.com/compose/compose-file/)
(`docker-compose.yml`)

However, it exists some setting to add specific amp behavior:


## add an automated mapping between a service and haproxy


an haproxy mapping allow to make public a service REST api. For instance you want this url: https://wwww.test.amp.appcelerator.io/myCommand, execute the REST api command "myCommand" using the port 80 of the service "myService" which is in the stack "test".
So you need a mapping between the logical name "www" used in the url with the internal port 80 of the service "myService"

to do so, add a io.amp.mapping with the value "www:80", in the definition of your service, in the yml stack file, as for instance:


version: "3"

services:
    pinger:
        image: appcelerator/pinger
        deploy:
            mode: replicated
            replicas: 2
            labels:
                io.amp.mapping: "www:3000"
        networks:
        - public

networks:
    public:
        external:
     

if you use this file to create a stack named "test": amp stack deploy test -c [thisFilePath]

then the url: https://www.test.appcelerator.io/myCommand, will execute the command "myCommand" using the port "3000" of the service "pinger" of the stack "test"

This label create a new entry in haproxy which know that all url starting by www.test. should be routed to the pinger service, port 3000.
the rest of the url "amp.amplifier.io" lead to haproxy ifself directly (in dev env.) or thtough cloud load-balancer.

So "amp.amplifier.io" should be an DNS entry in a public DNS domain as Route53.

to achieve the mapping, it's needed to declare the amp public network in the networks section of the compose file like:


networks:
    public:
        external:

and to attach your service to it in order to let him become visible by haproxy. All the other services not attached to public network stay hiden and private from outside.








