# Docker forked sources

Because Docker 1.13 does not provide a remote stack api, we need to embed some of the
Docker sources to provide this facility.

All updates are prefixed by {AMP} in code and consist of the following:

- remove all cobra dependencies
- make some functions public
- add public constructor on some private structs

## NOTICE

The original source files used in this package are copyrighted by Docker, Inc.,
under the [Apache 2.0 license](https://github.com/docker/docker/blob/master/LICENSE)
and can be found here:

https://github.com/docker/docker/

