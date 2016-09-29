# AMP AWS Instructions

The following readme is a general summary on how to access and use AMP within an AWS hosted VPC.
It is assumed that a swarm manager has already been provisioned and is running on AWS.



----------


##Prerequisites

 - [AMP Cli](https://github.com/appcelerator/amp#prerequisites) installed on developer machine

 -	AWS Instance running at least the swarm manager (A medium sized EC2 instance is recommended)

## Ports

By default, AMP uses HAPROXY to route traffic from external sources into the swarm cluster. The accessible ports should be:

### Port 8080
This port is reserved for the `amplifier` service, which is responsible for all internal interaction with the cluster.

###Port 80 and 443
These ports are reserved for HTTP and HTTPS respectively. Application services that have been deployed can be reached via HTTP through these ports. In order to route traffic to internal containers, a public DNS entry should be created in order allow access to internal containers. As an example, a service named `pinger` could be reached via `pinger.engage.amp.appcelerator.io`, assuming a DNS entry exists for the subdomain `engage.amp.appcelerator.io`


----------


## Using the CLI
In order to monitor and/or manage the AMP cluster the *server* option should be selected. For example,

`amp --server engage.amp.appcelerator.io:8080 stats`

Full documentation on the CLI functionality can be found [here] (https://github.com/appcelerator/amp#cli)


----------


## TROUBLESHOOTING

### Working around docker

The current version of docker that is running in the swarm is `1.12.1`. In some cases, during the development cycle, it might be necessary to manually manage the Docker installation, on the swarm host machine. Below are some useful commands to keep handy, when the situation arises.


***NOTE***: You will need to be given access to the pem file containing the private SSH Key in order to execute the following commands. For example:
`ssh -i ~/.ssh/amp-engage ubuntu@54.183.106.39`

 - Delete All Services: `docker service rm $(docker service ls -q) `
 - Remove All  Containers: `docker rm -f $(docker ps -aq)`
 - Remove the Docker Networks: `docker network rm $(docker network ls -q)`
 - Restart the Docker Daemon: `sudo service docker restart`
 - Remove Docker Volumes: `docker volume rm $(docker volume ls -q)`


### AMP cli known issues

#### Adding amp to Kaspersky antivirus scan exclusion list
While working on an Axway workstation, you might face an issue where antivirus protection will prevent your amp client to successfully establish connection with external amp services (working locally will of course not cause any issue).

 - Find instructions to allow your amp client [here](https://axway.jiveon.com/docs/DOC-31691)

 ***NOTE***: Make sure to disable network scan within trusted application for your amp client or VirtualBox environment.
