# Ubuntu

AMP is supported on these Ubuntu operating systems:

- Ubuntu Xenial 16.04 (LTS)
- Ubuntu Wily 15.10
- Ubuntu Trusty 14.04 (LTS)

This page instructs you to install using AMP release packages and installation mechanisms. Using these packages ensures you get the latest release of Docker.

>**Note**: Ubuntu Utopic 14.10 and 15.04 exist in Docker's `APT` repository but are no longer officially supported.

## Prerequisites

AMP requires a 64-bit installation regardless of your Ubuntu version.
Additionally, your kernel must be 3.10 at minimum. The latest 3.10 minor version or a newer maintained version are also acceptable.

>**Note**: If you want to start fresh, you can also follow this [quick tutorial](./vbox-xenial-tuto.md) to initialize your AMP environment.

### Install Git & configure your GitHub account

You can use the apt package management tools to update your local package index. Afterwards, you can download and install the program:

1. Download and install using `apt`:

	    $ sudo apt-get update
		$ sudo apt-get install git

2. Create your [GitHub account](https://github.com/) - if not already done

3. Configure `Git` using your GitHub credentials:

		$ git config --global user.name "Your Name"
		$ git config --global user.email "youremail@domain.com"

4. Optionally, you can also setup SSH for `Git` - instructions can be found [here](https://help.github.com/articles/generating-an-ssh-key/)

### Install Go

In this early adopter phase, AMP project still requires you to build **Go** scripts in order to generate binaries. The Go binary distributions assume they will be installed in `/usr/local/go`, but it is possible to install the Go tools to a different location. In this case you must set the **GOROOT** environment variable to point to the directory in which it was installed.

> **Note:** Installation flow below is based on current recommendation which is to install Go tools under `/go`.

1. Download [latest Go release](https://golang.org/dl/) (1.7.1 or higher) and extract it into `/go`. For example:

	    sudo tar -C /go -xzf go$VERSION.$OS-$ARCH.tar.gz

2. You must now set the **GOROOT** environment variable to point to the directory in which it was installed. You can do this by adding this line to your `/etc/profile` (for a system-wide installation) or `$HOME/.profile` (or preferentially `$HOME/.bashrc`)

	    export GOROOT=/go

3. Add `$GOROOT/bin` to your **PATH** environment variable. You can do this by adding this line to your `/etc/profile` (for a system-wide installation) or `$HOME/.profile` (or preferentially `$HOME/.bashrc`)

	    export PATH=$PATH:$GOROOT/bin

4. Set the **GOPATH** environment variable to point to the working directory holding go sources for AMP project. You can do this by adding this line to your `/etc/profile` (for a system-wide installation) or `$HOME/.profile` (or preferentially `$HOME/.bashrc`)

	    export GOPATH=/go


### Install Docker

Kernels older than 3.10 lack some of the features required to run Docker containers. These older versions are known to have bugs which cause data loss and frequently panic under certain conditions. Find all information related to [Docker](https://docs.docker.com/engine/installation/linux/ubuntulinux/)

To check your current kernel version, open a terminal and use `uname -r` to display your kernel version:

    $ uname -r
    3.11.0-15-generic

>**Note**: If you previously installed Docker using `APT`, make sure you update your `APT` sources to the new Docker repository.

#### Update your apt sources

Docker's `APT` repository contains Docker 1.7.1 and higher. To set `APT` to use packages from the new repository:

1. Log into your machine as a user with `sudo` or `root` privileges.

2. Open a terminal window.

3. Update package information, ensure that APT works with the `https` method, and that CA certificates are installed.

	    sudo apt-get update
	    sudo apt-get install apt-transport-https ca-certificates

4. Add the new `GPG` key.

		sudo apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D`

5. Open the `/etc/apt/sources.list.d/docker.list` file in your favorite editor.

    If the file doesn't exist, create it.

6. Remove any existing entries.

7. Add an entry for your Ubuntu operating system.

    The possible entries are:

    - On Ubuntu Trusty 14.04 (LTS)

            deb https://apt.dockerproject.org/repo ubuntu-trusty main

    - Ubuntu Wily 15.10

            deb https://apt.dockerproject.org/repo ubuntu-wily main

    - Ubuntu Xenial 16.04 (LTS)

            deb https://apt.dockerproject.org/repo ubuntu-xenial main

    > **Note**: Docker does not provide packages for all architectures. You can find
  > nightly built binaries in https://master.dockerproject.org. To install docker on
    > a multi-architecture system, add an `[arch=...]` clause to the entry. Refer to the
    > [Debian Multiarch wiki](https://wiki.debian.org/Multiarch/HOWTO#Setting_up_apt_sources)
    > for details.

8. Save and close the `/etc/apt/sources.list.d/docker.list` file.

9. Update the `APT` package index.

```
    $ sudo apt-get update
```

10. Purge the old repo if it exists.

```
    $ sudo apt-get purge lxc-docker
```

11. Verify that `APT` is pulling from the right repository.

```
    $ apt-cache policy docker-engine
```

    From now on when you run `apt-get upgrade`, `APT` pulls from the new repository.

#### Prerequisites by Ubuntu Version

- Ubuntu Xenial 16.04 (LTS)
- Ubuntu Wily 15.10
- Ubuntu Trusty 14.04 (LTS)

For Ubuntu Trusty, Wily, and Xenial, it's recommended to install the `linux-image-extra-*` kernel packages. The `linux-image-extra-*` packages allows you use the `aufs` storage driver.

To install the `linux-image-extra-*` packages:

1. Open a terminal on your Ubuntu host.

2. Update your package manager.

	    $ sudo apt-get update

3. Install the recommended packages.

	    $ sudo apt-get install linux-image-extra-$(uname -r) linux-image-extra-virtual

4. Go ahead and install Docker.


#### Install

Make sure you have installed the prerequisites for your Ubuntu version.

Then, install Docker using the following:

1. Log into your Ubuntu installation as a user with `sudo` privileges.

2. Update your `APT` package index.

        $ sudo apt-get update

3. Install Docker.

        $ sudo apt-get install docker-engine

4. Start the `docker` daemon.

        $ sudo service docker start


#### Create a Docker group    

The `docker` daemon binds to a Unix socket instead of a TCP port. By default that Unix socket is owned by the user `root` and other users can access it with `sudo`. For this reason, `docker` daemon always runs as the `root` user.

To avoid having to use `sudo` when you use the `docker` command, create a Unix group called `docker` and add users to it. When the `docker` daemon starts, it makes the ownership of the Unix socket read/writable by the `docker` group.

>**Warning**: The `docker` group is equivalent to the `root` user; For details on how this impacts security in your system, see [*Docker Daemon Attack Surface*](https://docs.docker.com/engine/security/security/) for details.

To create the `docker` group and add your user:

1. Log into Ubuntu as a user with `sudo` privileges.

2. Create the `docker` group.

        $ sudo groupadd docker

3. Add your user to `docker` group.

        $ sudo usermod -aG docker $USER

4. Log out and log back in.

    This ensures your user is running with the correct permissions.


#### Configure a DNS server for use by Docker

Systems that run Ubuntu or an Ubuntu derivative on the desktop typically use `127.0.0.1` as the default `nameserver` in `/etc/resolv.conf` file. The NetworkManager also sets up `dnsmasq` to use the real DNS servers of the connection and sets up `nameserver 127.0.0.1` in /`etc/resolv.conf`.

When starting containers on desktop machines with these configurations, Docker users see this warning:

    WARNING: Local (127.0.0.1) DNS resolver found in resolv.conf and containers
    can't use it. Using default external servers : [8.8.8.8 8.8.4.4]

The warning occurs because Docker containers can't use the local DNS nameserver. Instead, Docker defaults to using an external nameserver.

To avoid this warning, you can specify a DNS server for use by Docker containers. Or, you can disable `dnsmasq` in NetworkManager. Though, disabling `dnsmasq` might make DNS resolution slower on some networks.

The instructions below describe how to configure the Docker daemon running on Ubuntu 14.10 or below. Ubuntu 15.04 and above use `systemd` as the boot and service manager. Refer to [control and configure Docker with systemd](https://docs.docker.com/engine/admin/systemd/) to configure a daemon controlled by `systemd`.

To specify a DNS server for use by Docker:

1. Log into Ubuntu as a user with `sudo` privileges.

2. Edit Docker service.

		$ sudo systemctl edit docker

		[Service]
		EnvironmentFile=-/etc/default/docker
		ExecStart=
		ExecStart=/usr/bin/docker daemon -H fd:// $DOCKER_OPTS

	> **Note:** To save in Nano: Ctrl+X, "Yes" you want to save buffer, then on next page showing file path leave default value and press ENTER
	> This create a file /etc/systemd/system/docker.service.d/override.conf
	> The empty line with `ExecStart=` is here to clear current value, as only one `ExecStart=` could be declared (type=oneshot)

3. Open the `/etc/default/docker` file for editing.

        $ sudo nano /etc/default/docker

4. Add a setting for Docker.

        DOCKER_OPTS="--dns 8.8.8.8"

    Replace `8.8.8.8` with a local DNS server such as `192.168.1.1`. You can also
    specify multiple DNS servers. Separated them with spaces, for example:

        --dns 8.8.8.8 --dns 192.168.1.1

    >**Warning**: If you're doing this on a laptop which connects to various networks, make sure to choose a public DNS server. If you're **working within Axway network** add the following DNS server:

        --dns=10.252.252.252 --dns=10.253.253.253

4. Save and close the file.

5. Restart the Docker daemon.

        $ sudo service docker restart


#### Configure Docker to start on boot

Ubuntu uses `systemd` as its boot and service manager `15.04` onwards and `upstart`
for versions `14.10` and below.

For `15.04` and up, to configure the `docker` daemon to start on boot, run

    $ sudo systemctl enable docker

For `14.10` and below the above installation method automatically configures `upstart`
to start the docker daemon on boot


#### Verify `docker` is installed correctly

Verify your work by running `docker` without `sudo`.

	$ docker run hello-world

  If this fails with a message similar to this:

    Cannot connect to the Docker daemon. Is 'docker daemon' running on this host?

  Check that the `DOCKER_HOST` environment variable is not set for your shell.
  If it is, unset it.

### Install Glide

In this early adopter phase, and as you'll be working from sources, you will need to retrieve all dependencies to properly run AMP. We currently use [Glide](https://glide.sh/) to manage all associated packages. Installing Glide is a pretty straight forward process detailed below:

1. Log into your machine as a user with `sudo` or `root` privileges.

2. Open a terminal window.

3. Get Glide

		$ sudo curl https://glide.sh/get | sh


## Install AMP

Hang in there, you are now ready to install AMP on your system, just a few more steps:


### Retrieve latest build on GitHub

1. Log into your machine as a user with `sudo` or `root` privileges.

2. Open a terminal window.

3. Create your workspace directory:

		$ sudo mkdir -p /go/src/github.com/appcelerator


