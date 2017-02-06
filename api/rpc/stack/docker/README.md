# Docker forked sources

Because in 1.13.0 docker does not provide remote stack api, we need to embbed these Docker sources in order to be able to manage stack remotely.

Few updates have been done on the sources to be able to build them with as few dependencies as possible.

All updates are prefixed by {AMP} in code, mainly:

- remove all cobra dependencies
- make some functions public
- add public constructor on some private structs
