# ampwrapper

This project generates a small binary (`amp`) for supported platforms to provide a
convenient wrapper for running the actual amp cli in a Docker container without
requiring a user to supply all the options required to make that work.

For example, the Docker command currently looks like this:

    $ docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock --name amp --network swarmnet appcelerator/amp:local OPTIONS COMMAND

This convenience wrapper insulates the user against an unwieldy command as well as
potential future updates to it.

Usage:

    $ amp OPTIONS COMMAND

The binary is made available with amp releases and is also currently distributed under
the `bin` directory for the appropriate host operating system and architecture. A symlink
or shortcut should be placed in the user's path.

