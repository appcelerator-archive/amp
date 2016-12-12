# Ubuntu

AMP is supported on these Ubuntu operating systems:

- Ubuntu Yakkety 16.10
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

1. Download [latest Go release](https://golang.org/dl/) (1.7.3 or higher) and extract it into `/usr/local/go`. For example:

	    sudo tar -C /usr/local xzf go$VERSION.$OS-$ARCH.tar.gz

2. You must now set the **GOROOT** environment variable to point to the directory in which it was installed. You can do this by adding this line to your `/etc/profile` (for a system-wide installation) or `$HOME/.profile` (or preferentially `$HOME/.bashrc`)

	    export GOROOT=/usr/local/go

3. Add `$GOROOT/bin` to your **PATH** environment variable. You can do this by adding this line to your `/etc/profile` (for a system-wide installation) or `$HOME/.profile` (or preferentially `$HOME/.bashrc`)

	    export PATH=$PATH:$GOROOT/bin

4. Set the **GOPATH** environment variable to point to the working directory holding go sources for AMP project. You can do this by adding this line to your `/etc/profile` (for a system-wide installation) or `$HOME/.profile` (or preferentially `$HOME/.bashrc`)

	    export GOPATH=/go
	    export PATH=$PATH:$GOPATH/bin


### Install Docker

You can use the upstream script from Docker, it works with a large range of Linux distributions: https://get.docker.com/

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

## Install AMP

Hang in there, you are now ready to install AMP on your system, just a few more steps:


### Retrieve latest build on GitHub

1. Log into your machine.

2. Open a terminal window.

3. Create your workspace directory:

        $ mkdir -p $GOROOT/src/github.com/appcelerator
        $ cd $GOROOT/src/github.com/appcelerator 

4. Clone AMP git repository.

    > To use the master (unstable), type the following command:  

        $ git clone git@github.com:appcelerator/amp.git

    > To use a specific release (stable), type the following command:  

        $ git checkout v0.4.0
	# you can check the available releases here: https://github.com/appcelerator/amp/releases

4. Build AMP from source
 
         $ make install
    
You now have AMP installed locally! Please check the [getting started](../../README.md#configuration) section to start playing with AMP.
