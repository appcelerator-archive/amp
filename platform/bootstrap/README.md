# Deploy a Docker Swarm cluster with InfraKit

## Bootstrap

First, the cluster should be bootstrapped to get InfraKit running and ready to deploy the Swarm cluster.

Available bootstraps:

#### AWS

##### cloudformation bootstrap / AWS instance plugin

The Cloudformation template ```bootstrap.yml``` creates a VPC, subnet, internet gateway and the minimum required to build EC2 instances.
Select the instance type and the name of the EC2 Key pair name.
One EC2 instance will be created and will run InfraKit, it's public IP is revealed in the Cloudformation outputs.

The EC2 boostrap instance will run ```bootstrap -p aws```.

##### terraform bootstrap / Terraform instance plugin

The terraform-aws/bootstrap.tf creates the VPC, subnet, internet gateway and the minimum required to build EC2 instances.
3 managers are created, and infrakit then creates and monitor the worker nodes.

create a terraform-aws/terraform.tfvars file with the following content:
```
aws_name = "TAG"
aws_region = "us-west-2"
bootstrap_key_name = "YOUR-AWS-KEY-us-west-2"
aws_profile = "default"
infrakit_config_base_url = "https://raw.githubusercontent.com/appcelerator/amp/master/bootstrap"
```

and run
```
terraform apply $PWD/bootstrap/terraform-aws/
```

#### DigitalOcean

Coming soon

#### Docker in Docker

You can use the ```bootstrap``` script on you development machine to build a Swarm cluster composed of Docker containers.

## Deploy

The InfraKit instance renders the ```config.tpl``` template, and watches the groups defined in the transformed file (config.json).
This results in a full Swarm cluster.

## Security

The Swarm cluster is by default secured with self managed certificates.
The manager node remote API is enabled to allow nodes to get the join token, certificates are signed by a CA hosted on the bootstrap instance to enable mutual authentication.
The certificate management is only for demonstration purpose, not for production.
