
### Stack

The `amp stack` command is used to manage amp stacks

    $ `amp stack --help`  : display help

## amp stack deploy -c [stackFile] [stackName]

Deploy a stack using the compose file [stackFile] which is mandatory.

The docker name of the stack is going to be [stackName]-[id], where [id] is an unique id given by amp.

To update a stack with a new compose file use the id or the full name of the already deployed stack as with these commands:

`amp stack deploy -c [newStackFile] [id]`
or
`amp stack deploy -c [newStackFile] [stackName]-[id]`

## amp stack rm [id]

Remove the stack having the id [id]

## amp stack rm [stackName]-[id]

Remove the stack having the full name [stackName]-[id]

## amp stack ls

List the deployed stacks and give the following information:

- id: unique id given by amp
- name: short name given by user
- number of service
- owner: the owner name
- owner type: the owner type

optionally -q can be used to display only the stack id

Note that this command display only the amp stacks, all other stacks created out of amp are not displayed, infrastructure stacks are not displayed either.

# amp dedicated labels

## io.amp.mapping

When this label is added as a service label in a stack file, amp create a mapping to reverse proxy the service.

the label format is:

io.amp.mapping = " account=[myAccount] name=[aName] port=[internalServicePort]"

then the following urls can be used to request the service

http://[aName].[stackName].[myAccount].local.appcelerator.io
or
http://[aName].[stackName].[myAccount].cluster.atomiq.io
depending on your environment

without account=..., the url become:
http://[aName].[stackName].local.appcelerator.io
or
http://[aName].[stackName].cluster.atomiq.io

if name=... doesn't exist then name is set as the service name

[port] is mandatory
