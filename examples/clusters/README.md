# Cloudformation templates for Docker Swarm cluster creation on AWS

## How to use

Use the aws cli or the console to deploy a swarm cluster.
It is compatible with regions us-east-1, us-east-2, us-west-2, eu-west-1 and ap-southeast-2.

Alternatively it can be used through the aws plugin of AMP:

    docker run -it --rm -v ~/.aws:/root/.aws appcelerator/amp-aws:latest init --region us-west-2 --stackname STACKNAME --parameter KeyName=KEYNAME --template https://s3.amazonaws.com/io-amp-binaries/templates/latest/aws-swarm-asg.yml

## Content

The template will create the infrastructure (1 VPC, 3 subnets for HA on 3 datacenters, security groups, internet gateway, instance profile), 1 autoscaling group for the manager nodes, 1 autoscaling group for the amp core services worker nodes, 1 autoscaling group for the amp monitoring services worker nodes and 1 autoscaling group for the user services worker nodes.
Each autoscaling group run a userdata that initializes or join the swarm depending on the nature of the group and depending on the status of the swarm.
The engine API of all managers are available from all nodes in the VPC, which allow to set the labels on the nodes.

Once the nodes are up and running, it will run the appcelerator/ampadmin image to check the prerequisites, and setup AMP.

## Parameters

| Parameter | Description | Default Value | Example |
| --------- | ----------- | ------------- | ------- |
| KeyName   | Name of an existing EC2 KeyPair to enable SSH access to the instances | - | YOURNAME-REGION |
| ManagerSize | Number of manager nodes, should be 1, 3 or 5 | 3 | |
| CoreWorkerSize | Number of worker nodes for core services | 3 | |
| UserWorkerSize | Number of worker nodes for user services | 3 | |
| MetricsWorkerSize | Number of worker nodes for metrics services | 1 | |
| LinuxDistribution | AMI OS, Debian, Ubuntu or Default | Default | Ubuntu |
| ManagerInstanceType | Instance type for the manager nodes. Must be a valid EC2 HVM instance type | t2.small | m4.large |
| CoreInstanceType | Instance type for the core worker nodes. Must be a valid EC2 HVM instance type | m4.large | c4.large |
| UserInstanceType | Instance type for the user worker nodes. Must be a valid EC2 HVM instance type | t2.medium | m4.large |
| MetricsInstanceType | Instance type for the metrics worker nodes. Must be a valid EC2 HVM instance type | t2.large | m4.large |
| DrainManager | Should we drain services from the manager nodes? | no | yes |
| AufsVolumeSize | Size in GB of the EBS for the /var/lib/docker FS | 26 | 100 |
| OverlayNetworks | name of overlay networks that should be created once swarm is initialized | public core monit | public storage search mq |
| DockerChannel | channel for Docker installation | stable | edge |
| DockerPlugins | space separated list of plugins to install | | rexray/ebs |
| InstallApplication | install AMP | yes | no |
| NFSEndpoint | enable a NFSv4 service inside the VPC | no | yes |
| EnableSystemPrune | Enable Docker system prune | yes | no |
| MonitoringPort | Public port for the dashboard | 8080 | |
| EnableTLS | Docker socket secured with TLS | yes | |

## Output

The output of the stack lists the DNS name of the ELB in front of the manager nodes. It can be used for ssh access, https access to swarm services and configuration of the remote server in the CLI (--server option).

| Output | Description |
| --------- | ----------- |
| VpcId | VPC ID |
| DNSTarget | public facing endpoint for the cluster, It can be used for ssh access, https access to swarm services and configuration of the remote server in the CLI |
| InternalDockerHost | Docker host for services requiring access to the secured API |
| InternalRegistryTarget | internal endpoint for the registry service |
| MetricsURL | URL for cluster health dashboard |
| NFSEndpoint | NFSv4 Endpoint |
| InternalPKITarget | internal endpoint for the PKI service |


## Custom AMI

the default option for the AMI (Default) is a pre package AMI based on Ubuntu Xenial, with prerequisite packages already installed (in particular Docker).

To build a new version of this image, run the build-ami.sh script. It may take more than 15 min to build the AMI. You'll also need to create a variables.yml file with a content similar to:

```
---
ec2_key_name: "KEY_NAME"
```

Once done, copy the AMI Id in the cloudformation template (aws-swarm-asg.yml).

Docker will be installed with the latest version from the stable channel. If you want to build it with a specific version, you can add this line to your variables.yml:

```
docker_version: "17.09.1"
```

## Registry

An option of the template is the inclusion of a Docker registry.
It includes a S3 bucket as registry backend, and an autoscaling group of registry containers.
The registry is composed of non swarm nodes and is not part of the swarm.
The registry is only available from the VPC, all Docker swarm nodes are configured with the internal endpoint of the registry as mirror registry.

## Docker socket protected with TLS

When the EnableTLS option is set to yes, all swarm nodes are started with a certificate to protect the Docker socket. It is then available on port 2376 (instead of 2375 when not secured).
An autoscaling group with a single instance runs a cfssl docker container that is providing a CA, and serves an API for certificate generation. It is used by all the nodes to get the server certificate (for the Docker daemon) and a client certificate (for the Docker CLI).
The Manager external ELB has a listener on port 2376 that allows external access to the Docker engine API on the manager nodes (round robin). This is by default blocked by the security group, but can be open for a range of IP if a direct access is needed.
To be able to authenticate, you need the CA certificate, a key and a certificate. This is served by the cfssl container, but is not available from outside of the VPC (for security reason). This can be implemented as an AMP API, that would offer the interfaces with the PKS, generating and providing these 3 pem files to a client, with the added value of authorization. The API is not yet implemented.

#### Using the client certificate on a node of the swarm

Services can use the certificate available on the swarm. For instance, core services requiring access to the API on manager node can be scheduled on non manager nodes, and use the client certificate to connect to the API on the manager node. For that, the service has to mount the certificate from the host.

The certificate is available as well as its private key in /etc/docker: client.cert and client.key.

#### How to get the pem files for a Docker client

Identify the CA URL, it's the DNS name of the CA ELB. You can get it from the AWS console or from the status of the cluster (amp -s CLUSTER_URL cluster status). The URL should include the scheme.

The procedure below can be done only if you open the PKS service outside of the VPC, this is done by adding a rule to the PKS security group (look for CASecurityGroup).

From you client (usually a server from outside the swarm), do:
```
docker run --rm cfssl/cfssl:latest info -remote=CA_URL | jq -r .certificate
```

Paste the result on your machine in `~/.docker/ca.pem`.

Then, prepare a JSON file with the CSR.

```
{
    "CN": "USERNAME",
    "hosts": [
        "$(hostname)"
    ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "US",
            "L": "Santa Clara",
            "O": "Axway",
            "OU": "AMP",
            "ST": "California"
        }
    ]
}
```


Submit the CSR and save the response:

```
docker run --rm -v $PWD/csr.json:/csr.json cfssl/cfssl:latest gencert -remote=CA_URL -profile client /csr.json > response.json
```

Extract the key on your machine:
```
jq -r .key < response.json > ~/.docker/client.cert
```

Extract the certificate on your machine:
```
jq -r .cert < response.json > ~/.docker/client.cert
```

You can now use the Docker CLI by setting these variables:
```
export DOCKER_TLS_VERIFY=1
export DOCKER_HOST=MANAGER_EXTERNAL_ELB_DNS_NAME:2376
docker info
```
