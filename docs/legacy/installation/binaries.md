# Installation from binaries

**This instruction set is meant for hackers who want to try out AMP
on a variety of environments.**

Before following these directions, you should really check if a packaged
version of AMP is already available for your distribution. We have
packages for many distributions, and more keep showing up all the time!

## Check runtime dependencies


## Check kernel dependencies

AMP kernel requirements are those of the Docker daemon. A 3.10 Linux kernel is the minimum version supported, it was released in June 2013 so it's a safe guess that any recent Linux distro should meet the requirement.

Note that AMP also has a client mode, which can run on virtually any
Linux kernel (it even builds on MacOS and MS Windows!).

## Enable AppArmor and SELinux when possible


## Get the AMP binaries

The links are available on the github release page.

Substitute X.Y.Z in the links below with the latest available release.

### Get the Linux binaries

[get.amp.appcelerator.io](https://get.amp.appcelerator.io/builds/Linux/x86_64/amp-vX.Y.Z.tgz)

#### Install the Linux binaries

    tar xzf amp-vX.Y.Z.tgz -C /usr/local/bin/

#### Run AMP on Linux

### Get the Mac OS X binary

[get.amp.appcelerator.io](https://get.amp.appcelerator.io/builds/Darwin/x86_64/amp-vX.Y.Z.tgz)

### Get the Windows binary

[get.amp.appcelerator.io](https://get.amp.appcelerator.io/builds/Windows/x86_64/amp-vX.Y.Z.tgz)


## Giving non-root access


## Upgrade AMP


## Next steps
