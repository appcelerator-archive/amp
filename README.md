[![Build Status](https://travis-ci.org/appcelerator/amp.svg?branch=master)](https://travis-ci.org/appcelerator/amp)

# AMP

An open source [CaaS](https://blog.docker.com/2016/02/containers-as-a-service-caas/) for Docker, batteries included.

 * Use Docker Compose v3 [stackfiles](https://docs.docker.com/compose/compose-file/) to deploy your stacks
 * Account management support for users and teams with role-based access controls
 * Logs and metrics realtime filtered feeds and historical query support

While not recommended for production use quite yet, it's getting close (anticipated shortly after v0.17).
You'll be able to create your own HA cluster on the cloud, or use our playground.
In the meantime, you can get started on your own laptop with `amp cluster create` using the CLI.

## Getting started

For getting started and more detailed information, see [docs](docs/).

## Contributing

If you're already familiar with the project and want to hack on the project, we have a fully containerized toolchain.
All you need is Docker to build, test, and deploy! We would love for you to get involved,
so check out [CONTRIBUTING](project/CONTRIBUTING.md) and other docs under [project](project/).

## License

AMP is an open source project sponsored by [Axway](https://www.axway.com), and available under the Apache License, Version 2.0.
See [LICENSE](https://github.com/appcelerator/amp/blob/master/LICENSE)
for the full license text.
