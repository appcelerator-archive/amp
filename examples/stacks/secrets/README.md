# Secrets Example

## Overview

In this example, you will create a pair of x.509 certificates for the public and private keys
to encrypt and decrypt messages.

You will then create a compose file for a service that you will configure with the private key.
You will encrypt messages using the public key, then send the encrypted messages to the service.

You will store a private key in the swarm and make it available to the specific service to use.
This service will decrypt and log the message using the stored private key.

Another service will also attempt to access the secret, but will fail since it won't be
authorized.

## Setup

To create the private and public keys, we will use openssl.

If you don't have openssl alreay installed on your system, you can run it in container using
the following alias:

    $ alias openssl="docker run -it --rm -v $PWD:/root subfuzion/openssl"

### Create a private certificate

    $ openssl genrsa -out privatekey 1024
    Generating RSA private key, 1024 bit long modulus
    ..........++++++
    ...................++++++
    e is 65537 (0x10001)

### Create a public certificate

    $ openssl rsa -in privatekey -out publickey -outform PEM -pubout
    writing RSA key

### Use the public key to encrypt a message

    $ echo "hello world" > message.txt
    $ openssl rsautl -encrypt -inkey publickey -pubin -in message.txt -out message.dat

Verify the message is encrypted:

    $ cat message.dat
    ^PUIMqmf?+'WţHHl@F` kw<m8YgV   NrJ5p:           '5'\שu

Verify the private key can unencrypt it:

    $ openssl rsautl -decrypt -inkey privatekey -in message.dat

## Create a stack file for the service and secret

See the existing `stack.yml` file.

```
version: "3.1"

services:

  echo:
    image: "subfuzion/secure-echo"
    ports:
      - "8887:8887"

  secure_echo:
    image: "subfuzion/secure-echo"
    ports:
      - "8888:8888"
    secrets:
      - privatekey

secrets:
  privatekey:
    file: ./privatekey

```

## Deploy the stack

    $ amp stack deploy -c stack.yml demo

This will start two services: `demo_echo` listening on port `8887` and `demo_secure_echo` listening on `8888`.

## Test the services

Use netcat (`nc`) to send the encrypted message to each service and check the service logs for the result.
`demo_echo` is not configured to use the secret and therefore will fail to decrypt the message.
`demo_secure_echo` should succeed.

If you don't have netcat installed on your system, you can create an alias to use the version in the secure-echo image:

    $ alias nc="docker run -it --rm --entrypoint nc subfuzion/secure-echo"

Test `demo_echo`:

    $ cat message.dat | nc localhost 8887
    $ amp service logs demo_echo
    ...
    unable to load Private Key

Test `demo_secure_echo`:

    $ cat message.dat | nc localhost 8888
    $ amp service logs demo_secure_echo
    ...
    hello world

