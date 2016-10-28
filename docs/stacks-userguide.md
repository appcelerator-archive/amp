# AMP - Stacks documentation



**Summary:**
- 1 Stack commands
- 2 Stack file description
- 3 ETCD storage
- 4 Networks
- 5 HAProxy



# 1 Stack Commands



## 1.1 Create a new stack


Command: **amp stack create `stack-name` -f `stack-file`**


- Creates a new stack having the name  `stack-name` using the `stack-file`  definition (see chapter stack file description). 
- Stores stack information in ETCD (see chapter ETCD storage). 


## 1.2 Start a stack


Command: **amp stack start `stack-name`**


- Retrieve the stack `stack-name` information from ETCD
- Creates or reuses stack private network
- If it exists at least one service having a public name, creates or reuses stack public network (see chapter Networks)
- Creates or reuses stack custom networks, if any are defined for this stack
- Stores stack networks information in ETCD
- Starts all defined services
- If it exists at least one service having public name definition, starts also a service HAProxy dedicated to the stack usage. (see chapter HAProxy)
- Stores stack services information in ETCD


## 1.3 stack up


Command: **amp stack up `stack-name` -f `stack-file`**


- Perform a stack create and a stack start in the same command.


## 1.4 Stop an existing stack


Command: **amp stack stop `stack-name`**


- Stops the stack `stack-name` services, including the stack HAProxy if exists
- Remove stack services information stored in ETCD


## 1.5 Remove a stack


Command: **amp stack rm `stack-name` [-f]**


- With -f, stops the stack `stack-name` if it is still running.
- Removes private stack network and public network if exists
- Removes custom stack network if not used by another running stack
- Removes all stack information stored in ETCD
- Stack name become free and can be reused in create or up command using any stack file.
 



# 2 Stack file description (version 1)




## 2.1 Global view


A stack file is a file respecting the yaml syntax. It contains all the needed information to create astacks: services, networks, volumes, …
All the stack information are stored in ETCD during the `amp stack create/up` command and then this file is not needed anymore during the rest of the stack life cycle.

The same stack file description can be used in several `amp stack create/up` commands using different stack names. Until the `amp stack rm` commands, a stack name can’t be reused, but a stack file description can.




## 2.2 version item 

`(not implemented yet, version is 1 for now, but this item should be set anyway now to distinct future versions)`


This item set the format version of the amp stack file. The stack file parsing, executed during amp stack create/up command, depends on this version.
By default, the amp stack file version is 1.




## 2.3 Global services item


services is a mandatory item and the most important one of the stack file. It describes one by one all the services used in the stack, like in the following one where two services **myservice1** and **myservice2** are declared:

```
version: 1
services:
        myservice1:
         …
        myservice2:
        ….
```


## 2.4 Service description items


Under each service item, as the previous **myservice1** and **myservice2** ones, we can add specific service items:



### 2.4.1 image, service item


set the docker image used by the service. This item is mandatory

```
Sample:
myservice1:
    image: appcelerator/pinger:latest
```



### 2.4.2 mode, service item


Sets the mode used by the service to start. 
Two modes are possible:
- global
- replicated


##### Global mode:


When the service is started, one container is created on each swarm nodes


Sample:
```
myservice1:
    image: appcelerator/pinger:latest
    mode: global
```

##### Replicated mode:


When the service is started a given number of containers will be created and spread on swarm nodes. The item replicas is used to defined the number of replicas.

Sample:
```
myservice1:
    image: appcelerator/pinger:latest
    mode: replicated
    replicas: 2
```

The default mode is replicated with one replicas. We can also use only **replicas** with a value >0, it sets the mode **replicated**


### 2.4.3 environment, service item


Adds environment variables to all the containers created by the service
By default, no environment variable is added.


Two syntaxes are supported map or list:


Sample with map
```
myservice1:
    image: appcelerator/pinger:latest
    environment:
      myvar1:  myvalue1
      myvar2: myvalue2
```

Sample with list
```
myservice1:
    Image: appcelerator/pinger:latest
    environment:
      - “myvar1= myvalue1”
      - “myvar2=myvalue2”
```

### 2.4.4 labels, service item 


Adds labels to the service (visible using docker service inspect)
By default no labels are added


Two syntaxes are supported map or list:


Sample with map
```
myservice1:
    image: appcelerator/pinger:latest
    labels:
      mylabel1:  myvalue1
      mylabel2: myvalue2
```

Sample with list
```
myservice1:
    image: appcelerator/pinger:latest
    labels:
      - “mylabel1= myvalue1”
      - “mylabel2=myvalue2”
```

### 2.4.5 containers labels, service item


Adds labels to all the containers created for the service (visible using docker inspect)
By default no container labels are added


Two syntaxes are supported map or list:


Sample with map
```
myservice1:
    image: appcelerator/pinger:latest
    container_labels:
      mylabel1:  myvalue1
      mylabel2: myvalue2
```

Sample with list
```
myservice1:
    image: appcelerator/pinger:latest
    container_labels:
      - “mylabel1= myvalue1”
      - “mylabel2=myvalue2”
```

### 2.4.6 public, service item


Define how the service can be requested from outside the stack.
There are two ways:
- Using name
- Using publish port

For each, a protocol can be defined, tcp by default




##### Using name:


Several public names can be set for a service. Each of them defines a relation between an internal service port and a logical name that the stack HAProxy can use to request the service (see HAProxy chapter for more details)


Sample:
```
myservice1:
    image: appcelerator/pinger:latest
    public:
      - name: myexternalname
        protocol: tcp
        internal_port: 3000
```

Then: `http://myexternalname.[domain]/xxx`
will be load-balanced to: `[myService containers addresses]:3000/xxx`




##### Using publish_port:


Several publish ports can be set for a service. Each of them defined a relation between the publish port and an internal service port. On all swarm node any http request on the publish port will be rerouted to the service internal port (see Network chapter for more details)


Sample:
```
myservice1:
    image: appcelerator/pinger:latest
    public:
       - publish_port: 3001
         internal_port: 3000
```

Then: `http://[any swarm node address]:3001/xxx`
will be load-balanced to: `[myService container addresses]:3000/xxx`


##### Using both name and publish_port


It’s possible to use name and publish port on the same internal port. They should use the same protocol.


Sample:
```
myservice1:
    image: appcelerator/pinger:latest
    public:
      -name: myexternalname
       publish_port: 3001
       internal_port: 3000
```

##### With several internal ports:

```
myservice1:
    image: appcelerator/pinger:latest
    public:
      - name: myexternalname1
        publish_port: 3001
        internal_port: 3000
      - name: myexternalname2
        protocol: tcp
        publish_port: 6001
        internal_port: 6000
```


### 2.4.7 networks, service item


Attaches the service containers to a custom or an external network (see Network chapter for more details). Several services DNS aliases can be defined to find the services containers addresses on the attached network.
By default, only one alias is defined: The service name itself.


Sample:
```
myservice1:
    image: appcelerator/pinger:latest
    networks:
       mynetworkname:
          aliases:
              -myalias1
              -myallias2
```

**mynetworkname** should be a custom or an external network defined in networks global item 



### 2.4.8 volumes, service item


Creates volume(s) in each service containers, Several kind of volumes are possible:
- Anonymous: 
- Host linked
- Named

With optional read-write or read-only specifications.
By default a volume is read-write


##### Anonymous


Volumes are created in each service containers at each startup.


Sample:
```
myservice1:
    image: appcelerator/pinger:latest
    volumes:
       - myvolume1folder
       - myvolume2folder
```

All paths should be absolute.


##### Host linked


The volumes are a link between a host folder and a container folder

Sample:
```
myservice1:
    image: appcelerator/pinger:latest
    volumes:
     - hostFolder2:myvolume1folder
     - hostFolder2:myvolume2folder
```

All paths should be absolute.


##### Named


The volumes have name and are re-used if exists at service startup.

Sample:
```
myservice1:
    image: appcelerator/pinger:latest
    volumes:
     - volumename1:myvolume1folder
     - volumename2:myvolume2folder
```

All paths should be absolute.


##### Options

Optional parameters can be added to volumes definitions. See docker documentation for more details.
Parameters can be added following the same syntaxe than docker compose.


Sample:
```
myservice1:
    image: appcelerator/pinger:latest
    volumes:
      - volumename1:myvolume1folder:ro
      - volumename2:myvolume2folder
```

Then, the volumename1 will be in read-only mode.



## 2.5 Global networks item


networks is an optional item. It describes one by one all the custom or external networks used by the stack services in their own networks item, like in the following where two networks **mynetwork1** and **mynetwork2** are declared:

```
version: 1
services:
        ….
networks:
        mynetwork1:
         …
        mynetwork2:
        ….
```



Two kinds of networks can be defined:
- Custom network: It’s  a network created when a stack starts and removed when a stack stops (not necessarily by the same stack, see network chapter)
- External network: It’s a network that pre-exist to the stack. They should be started before the stacks using it and a stack never removes it.




## 2.6 Custom networks description items


Under each network item, as the previous **mynetwork1** and **mynetwork2** ones, we can add specific custom networks description items:


### 2.6.1 driver, custom network item


driver is a mandatory information to define a custom network. It specify which driver should be used for this network. It could take several values (see docker documentation), but on amp contexte most of the time it’ll be “overlay” or “bridge”



Sample:
```
mynetwork1:
   driver: overlay
```

That the minimum information needed to define a custom network. All the other parameters can be taken by default.


### 2.6.2 driver_opts, custom network item


driver_opts is a set of parameters used to define the driver options


Sample:
```
mynetwork1:
   driver: overlay
   driver_opts:
      myopt1: value1
      myopt2: value2
```

### 2.6.3 ipam, custom network item


Specify custom IPAM config having several properties, each of which is optional:
- driver: custom IPAM driver, instead of the default
- options: list of driver options
- config: specific IPAM parameters:
- subnet: Subnet in CIDR format that represents a network segment
- ip_range: Range of IPs from which to allocate container IPS
- gateway: IPv4 or IPv6 gateway for the master subnet
- aux_addresses: Auxiliary IPv4 or IPv6 addresses used by Network driver as a mapping from hostname to IP.


Sample:
```
mynetwork1:
   driver: overlay
   ipam:
     driver: default
     config:
       - subnet: 172.28.0.0/16
         ip_range: 172.28.5.0/24
         gateway: 172.28.5.254
         aux_addresses:
          host1: 172.28.1.5
          host2: 172.28.1.6
          host3: 172.28.1.7
```


### 2.6.4 labels, custom network item


Adds labels to the network (visible using docker network inspect)
By default no labels are added


Sample:
```
mynetwork1:
   driver: overlay
   labels:
     mylabel1: value1
    mylabel2: value2
```


### 2.6.5 external network item


An external network is defined only by its name, but the service that is attached to it can choose to use another name inside the stack file. This name makes sense only inside the stack file



##### Sample using real name:

```
myservice1:
    image: appcelerator/pinger:latest
    networks:
       anexternalnetwork:
          aliases:
              -myalias
networks:
   anexternalnetwork:
      external: true
```

In this example, the service **myservice1** is going to be attached to an external network named **anexternalnetwork**. **Externanetwork** is the real name of the network.




##### Sample using network alias name

```
myservice1:
    image: appcelerator/pinger:latest
    networks:
       mynetworkname:
          aliases:
              -myalias
networks:
   mynetworkname:
      external: 
         name: anexternalnetwork
```

In this example, the service **myservice1** is going to be attached to an external network named **anexternalnetwork** too, but the service reference the network **mynetworkname** which itself reference the real external network **anexternalnetwork**



# 3 ETCD Storage


The infrastructure service ETCD is used to store stack informations in dedicated keys/values. 


The information is spreaded in three main root keys:
- amp/stacks  
- amp/services
- amp/networks


```
Keys                            Content                             Structure
_______________________________________________________________________________________
amp/stacks/[stackId]            Stack whole definition              stack.Stack
amp/stacks/[stackId]/networks   Stack network ids list              stack.IdList
amp/stacks/[stackId]/services   Stack service ids list              stack.IdList
amp/stacks/names/[stackName]    Stack Id                            stack.StackId
amp/services/[serviceId]        Service whole definition            service.ServiceSpec 
amp/networks/[networkId]        Custom network whole definition     stack.customNetwork
```



## 3.1 amp/stacks/`[stackId]`


This key is created by the commands `amp stack create` and `amp stack up` to store the parsing result of the stack file. This key is removed by the command amp `stack rm`



## 3.2 amp/stacks/`[stackId]`/networks


This key is created by the command `amp stack start` to store the list of the created network ids. This key is removed by the command amp `stack rm`



## 3.3 amp/stacks/`[stackId]`/services


This key is created by the command `amp stack start` to store the list of the created services ids. This key is removed by the command `amp stack stop`.



## 3.4 amp/stacks/names/`[stackName]`


This key is created by the commands `amp stack create` and `amp stack up` and used by all stack commands to get stack id by its name or to verify if a stack name already exist. This key is removed by `amp stack rm`.



## 3.5 amp/services/`[serviceId]`


This key is created by the command `amp stack start` to store the whole service definition and removed by the command `amp stack stop`.
This key is watched by amp HAProxy controller to update the stack HAProxy configuration on any stack service change.


## 3.6 amp/networks/`[networkId]`


This key is created by the command `amp stack start` if the custom network doesn’t already exist. If the network already exists, the command `amp stack start` just increments the number of the network owners in the key.
The owner number of this key is decremented  by the command `amp stack rm`.  if the number of owner becomes 0, the command remove the network and the key.



# 4 Networks


## 4.1 Default amp networks


Out of stack, amp manage two networks:
- Amp-infra:  It’s a dedicated network for all infrastructure services
- Amp-public: it’s a dedicated network for infrastructure services that needs to be public, mainly HAProxy and UIs as grafana, amp-iu, …


## 4.2 Default stack networks


When a stack is started, two default networks are created:
- A stack network, named `[stackName]`-private
- A stack network, named `[stackName]`-public



### 4.2.1 private stack network


The private stack stack network is created at the first stack start and removed only when the stack itself is removed.
Each service of the stack are attached on it by default with their short name as alias.
(The short name of a stack service is the name declared in the stack file, the long name of a stack service is form with the stack name dash its short name `[stackName]-[shortName]`)
This private network allows services to request one to each other using just their short name as DNS address. 
The amp networks use virtual IPs by default meaning that using the service alias as DNS name, requests are naturally load-balanced to all the targeted service containers.




### 4.2.2 public stack network


The public stack network is created at the first stack start, only if at least one service have a public name declared in this configuration. This network is removed when the stack itself is removed.
All the stack services having at least one public name are attached on this stack public network
In addition a stack HAProxy service is created and also attached to the stack public network.
This network allows to isolate the public services. They are requestable by HAProxy which are on the same network.
Then the stack HAProxy won’t be able to request a service only attached on the stack private network. (see HAProxy chapter for more details)



### 4.2.3 custom networks


A custom network is designed to allow inter-stacks communication and to be manageable by the stack themselves.
Meaning that several stacks can declare the same custom network and this network is created only once.
A stack service can be allowed to be attached to the custom network or not, if so its alias on this custom network is its long name ([stackName]-[shortServiceName])
Then two services belonging to two distinct stacks and allowed to be attached to the custom network can request each other using these aliases as DNS name.
A custom network is created by the first started stack declaring it and deleted by the last removed stack declared it.



### 4.2.3 external networks


An external network allows stack inter-communication the same way than the custom networks. The only difference with custom network is that an external network is never created or deleted by a stack. It should pre-exist to the stacks


### 4.2.4 Docker ingress network


When a service declare a publish port, it is attached to the default docker network ingress. Then each swarm node listen the publish port. When a request arrive it is load-balanced on the service containers internal port, no matter on which node they are.
This works even if the request reaches a node on which no container of the targeted service is running. The request will be rerouted to a node having such container.


# 5 HAProxy


## 5.1 global view


Amp HAProxy services uses host name to load-balance requests on the targeted service. They expect that the request host name has one of the form:
- `[servicePublicName]`.`[StackName]`.`[domain]`
- `[servicePublicName]`.`[domain]`

These two forms are evaluated in two kind of HAProxy:
- Amp Infrastructure HAProxy
- Stack HAProxy


## 5.2 Amp infrastructure HAProxy


This service exist once in the amp platform. It’s the main entry for services requests. It’s attached to the amp default network amp-public and listen on the ports 80, 8080 (443)
Its role is to either :
Load-balance requests to public infrastructure services
load-balance user service requests to the right stack HAProxy service.


The [domain] used to request services should be defined in order that any *.[domain] DNS names is resolved to the infrastructure HAProxy containers addresses.


### 5.2.1 Load-balance requests to public infrastructure services


Amp infrastructure HAProxy expects that Infrastructure service request hosts have the form
`[infrastructureServiceName]`.`[domain]`


If so, it load-balances the request to the targeted infrastructure service containers using the DNS name [infactructureServiceName]
It means that the infrastructure services which are allowed to be public are also attached to the network amp-public and have their own name as DNS alias.
The possible infrastructure services able to be requested this way are mainly UI (grafana, amp-ui).
The service amplifier uses GRPC protocol. It is specifically reversed-proxy using the entry port 8080. So each GRPC request arriving on the HAPROXY port 8080 is load-balanced on the amplifier containers listening on 50101 (their internal port).


### 5.2.2 Load-balance user service requests to stack HAProxy services.


Amp HAProxy infrastructure expects that user service request hosts have the form:
`[servicePublicName]`.`[stackName]`.`[domain]`


If so, it load-balances the requests to the corresponding stack HAProxy containers using the DNS name HAProxy-[stackName] where [stackName] has been extracted from request hosts.


## 5.3 Stack HAProxy


A stack HAproxy service exists in each stack having public services (declaring at least one public name in its configuration)
A stack HAProxy is attached to amp-public network with an alias HAProxy-{stackName] and is attached also to the stack public network. This way it is able to received requests from infrastructure HAProxy and to address them the stack public services.


It expects to receive request hosts having either the form:
- `[servicePublicName]`.`[domain]`
- `[servicePublicName]`.`[stackName]`.`[domain]`

If so, it load-balances the requests to the targeted service containers using the DNS name [servicePublicName] where [servicePublicName] is extracted from request hosts.
Then, only public services attached on the stack public network are able to be requested.



