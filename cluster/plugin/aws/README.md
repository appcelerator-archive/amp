# AWS Plugin

This is the plugin for creating and initializing an AWS cluster.

`Dockerfile.compiler` builds the image for compiling the Go source files for
the plugin.

`Dockerfile` builds the image for the actual plugin that will be used by the
AMP cli.

For more details about the design and use, see the
[wiki](https://github.com/appcelerator/amp/wiki/AWS-Clusters).

### Options

The plugin allows you to provide all the parameters that are supported by the [Docker for AWS CloudFormation template](https://docs.docker.com/docker-for-aws/#configuration-options), including:

 * KeyName
 * EnableCloudWatchLogs
 * ManagerSize
 * ManagerInstanceType
 * ManagerDiskSize
 * ManagerDiskType
 * ClusterSize
 * InstanceType
 * WorkerDiskSize
 * WorkerDiskType

The format for adding a parameter (`-p | --parameter`) is more compact than the format used by the AWS CLI. For example, instead of:

    --parameters ParameterKey=ManagerInstanceType,ParameterValue=t2.medium ...

use this format instead:

    --parameter ManagerInstanceType=t2.medium ...

And for updates, instead of:

    --parameters ParameterKey=ManagerInstanceType,UsePreviousValue=true ...

just use this instead:

    --parameter ManagerInstanceType   # no assigned value assumes UsePreviousValue=true ...

## Prerequisites

* Ensure you have a keypair for the region you want to use to create a stack using the plugin. I recommend you upload your own public key and give it the same name for each region you want to test (see [keypair docs](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html#how-to-generate-your-own-key-and-import-it-to-aws)).

* Ensure your AWS credentials are configured (see [credentials docs](http://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html)).

### Sync Operations

You can block until a command completes using the `-s | --sync` option. For example:

    $ docker run -it --rm -v ~/.aws:/root/.aws appcelerator/amp-aws-plugin update --region us-west-2 --stackname tony-amp-10 -p KeyName=tony-amp-dev -p ClusterSize=4 --sync

## Trying it out

From the `cluster/plugin/aws` directory, run the following:

    $ make compiler
    $ make
    $ export REGION=us-west-2 # or any region where you have a valid keypair
    $ export KEYNAME=tony-amp-dev
    $ export STACKNAME=tony-amp-1
    $ export ONFAILURE=DO_NOTHING
    $ docker run -it --rm -v ~/.aws:/root/.aws appcelerator/amp-aws-plugin init \
        --keypair $KEYPAIR \ <- {FR}: seems that this line has to be deleted to make this cmd run
        --region $REGION \
        --stackname $STACKNAME \
        --onfailure $ONFAILURE \
        --parameter KeyName=$KEYNAME \
        --parameter ...

Verify output similar to the following:

```
$ docker run -it --rm -v ~/.aws:/root/.aws appcelerator/amp-aws-plugin init --region us-west-2 --name tony-amp-7
2017/06/14 20:31:23 {
  StackId: "arn:aws:cloudformation:us-west-2:654814900965:stack/tony-amp-7/6d88d410-5140-11e7-8af9-503acbd4dc29"
}
2017/06/14 20:31:23 2010-05-15
```

Also verify the stack has been created in the AWS console. The URL looks like this (updated the region value as appropriate):

https://us-west-2.console.aws.amazon.com/cloudformation/home?region=us-west-2#/stacks?filter=active&tab=resources

### SSH to an instance

On the page opened by the previous link, you can click on your stack, then click on "Outputs", then click on the "Managers" link to the manager instances that were created. From there, you can obtain the public IP and DNS names created for each manager. Using the ssh key associated with your AWS keypair, you can ssh into any of these instances.

For example, if your private key is called `id_rsa` and the IP address of a manager node is `34.210.103.37`, then you can ssh into it like this:

    $ ssh -i ~/.ssh/id_rsa docker@ 34.210.103.37

You can then run docker swarm commands, like `docker node ls`.

Once you are connected to a manager node you can then ssh to a worker node using the node name as reported by `docker node ls`. Make sure you have ssh agent forwarding set up from your host system to ensure your keys propagate so that your private key doesn't need to be installed on each node you ssh from (see [docs](https://docs.docker.com/docker-for-aws/deploy/#using-ssh-agent-forwarding)).

    # can only ssh to a worker from a manager node
    $ ssh <worker>

### SSH Tunnel

If you don't need to ssh to a work and just want to connect to a manager for running docker commands, it can be more convenient to open an ssh tunnel. You can do this as follows:

    $ ssh -i ~/.ssh/id_rsa -NL 127.0.0.1:2374:/var/run/docker.sock docker@34.210.103.37

Or run the command and put it in the background -- just don't forget to kill it when done to release the port.

    $ ssh -i ~/.ssh/id_rsa -NL 127.0.0.1:2374:/var/run/docker.sock docker@34.210.103.37 &

Then you can specify the docker host to connect to using that port for the tunnel to run docker commands. For example:

    $ docker -H localhost:2374 node ls

Finally, you can set `DOCKER_HOST` so you don't have to use the `-H` option each time:

    $ export DOCKER_HOST=localhost:2374
    $ docker node ls

### Update Stack

Update an existing stack with the `update` command:

    $ docker run -it --rm -v ~/.aws:/root/.aws appcelerator/amp-aws-plugin update \
        --region $REGION \
        --stackname $STACKNAME \
        -p KeyName=$KEYNAME \
        -p ClusterSize=4 # new value! \
        --sync # block until done!

### Destroy Stack

Execute the `destroy` command with the stackname and region used for the `init` command:

    $ docker run -it --rm -v ~/.aws:/root/.aws appcelerator/amp-aws-plugin destroy \
        --region $REGION \
        --stackname $STACKNAME

## Tests

    $ KEYNAME=<aws-keypair-name> REGION=<aws-region> make test

Tests now automatically destroy the stacks they create.
