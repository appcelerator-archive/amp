# Example: Manual install on cloud provider

You can install AMP directly to servers you have on cloud providers.  This example shows how to create an <a href="https://aws.amazon.com/" target="_blank"> Amazon Web Services (AWS)</a> EC2 instance, and install AMP on it.

You can use this same general approach to create AMP-ready hosts on other cloud providers.

### Step 1. Sign up for AWS

1. If you are not already an AWS user, sign up for <a href="https://aws.amazon.com/" target="_blank"> AWS</a> to create an account and get root access to EC2 cloud computers. If you have an Amazon account, you can use it as your root user account.

2. Create an IAM (Identity and Access Management) administrator user, an admin group, and a key pair associated with a region.

    From the AWS menus, select **Services** > **IAM** to get started.

    See the AWS documentation on <a href="http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/get-set-up-for-amazon-ec2.html" target="_blank">Setting Up with Amazon EC2</a>. Follow the steps for "Create an IAM User" and "Create a Key Pair".

    If you are just getting started with AWS and EC2, you do not need to create a virtual private cloud (VPC) or specify a subnet. The newer EC2-VPC platform (accounts created after 2013-12-04) comes with a default VPC and subnet in each availability zone. When you launch an instance, it automatically uses the default VPC.

### Step 2. Configure and start an EC2 instance

Launch an instance to create a virtual machine (VM) with a specified operating system (OS) as follows.

  1. Log into AWS with your IAM credentials.

      On the AWS home page, click **EC2** to go to the dashboard, then click **Launch Instance**.

      ![EC2 dashboard](../images/ec2_launch_instance.png)

      AWS EC2 virtual servers are called *instances* in Amazon parlance. Once you set up an account, IAM user and key pair, you are ready to launch an instance. It is at this point that you select the OS for the VM.

  2. Choose an Amazon Machine Image (AMI) with the OS and applications you want. For this example, we select an Ubuntu server.

      ![Launch Ubuntu](../images/ec2-ubuntu.png)

  3. Choose an instance type.

      ![Choose a general purpose instance type](../images/ec2_instance_type.png)

  4. Configure the instance.

    You can select the default network and subnet, which are inherently linked to a region and availability zone.

      ![Configure the instance](../images/ec2_instance_details.png)

  5. Click **Review and Launch**.

  6. Select a key pair to use for this instance.

    When you choose to launch, you need to select a key pair to use. Save the `.pem` file to use in the next steps.

The instance is now up-and-running. The menu path to get back to your EC2 instance on AWS is: **EC2 (Virtual Servers in Cloud)** > **EC2 Dashboard** > **Resources** > **Running instances**.

To get help with your private key file, instance IP address, and how to log into the instance via SSH, click the **Connect** button at the top of the AWS instance dashboard.


### Step 3. Log in from a terminal and verify proper installation

1. Log in to the EC2 instance from a command line terminal.

    Change directories into the directory containing the SSH key and run this command (or give the path to it as part of the command):

        $ ssh -i "YourKey" ubuntu@xx.xxx.xxx.xxx

    For our example:

        $ cd ~/Desktop/keys/amazon_ec2
        $ ssh -i "my-key-pair.pem" ubuntu@xx.xxx.xxx.xxx

    We'll follow the instructions for installing AMP on Ubuntu. The next few steps reflect those instructions.

2. Check the kernel version to make sure it's 3.10 or higher.

        ubuntu@ip-xxx-xx-x-xxx:~$ uname -r
        3.13.0-48-generic


### Step 4. Install Docker Engine on the remote instance

1. Download the Docker install script and pipe to shell.

        ubuntu@ip-xxx-xx-x-xxx:~$ curl -sSL https://get.docker.com/ | sh


2. Start the Docker daemon.

        ubuntu@ip-xxx-xx-x-xxx:~$ sudo service docker start

3. Verify Docker Engine is installed correctly by running `docker run hello-world`.

**Note:** docker command will need to be run as sudo, unless user is added to docker group. eq.,
`sudo usermod -aG docker ubuntu `


        ubuntu@ip-xxx-xx-x-xxx:~$ sudo docker run hello-world
        ubuntu@ip-172-31-0-151:~$ sudo docker run hello-world
        Unable to find image 'hello-world:latest' locally
        latest: Pulling from library/hello-world
        b901d36b6f2f: Pull complete
        0a6ba66e537a: Pull complete
        Digest: sha256:8be990ef2aeb16dbcb9271ddfe2610fa6658d13f6dfb8bc72074cc1ca36966a7
        Status: Downloaded newer image for hello-world:latest

        Hello from Docker.
        This message shows that your installation appears to be working correctly.

        To generate this message, Docker took the following steps:
        1. The Docker client contacted the Docker daemon.
        2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
        3. The Docker daemon created a new container from that image which runs the executable that produces the output you are currently reading.
        4. The Docker daemon streamed that output to the Docker client, which sent it to your terminal.

        To try something more ambitious, you can run an Ubuntu container with:
        $ docker run -it ubuntu bash

        Share images, automate workflows, and more with a free Docker Hub account:
        https://hub.docker.com

        For more examples and ideas, visit:
        https://docs.docker.com/userguide/

### Step 5. Install AMP on the remote instance

1. Step 1

2. Step 2

## Where to go next
