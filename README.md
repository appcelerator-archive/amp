# AMP

The open source unified CaaS/FaaS platform for Docker, batteries included.

 * Host your own high availability cluster or use `cloud.atomiq.io`
 * [Docker Infrakit](https://github.com/docker/infrakit) for self-healing infrastructure
 * Use Docker Compose v3 [stackfiles](https://docs.docker.com/compose/compose-file/) to deploy your stacks ([CaaS](https://blog.docker.com/2016/02/containers-as-a-service-caas/))
 * Support lambda style [serverless](https://en.wikipedia.org/wiki/Serverless_computing) tasks ([FaaS](https://martinfowler.com/articles/serverless.html))
 * Account management support for users, organizations and teams with role-based access controls
 * Logs and metrics realtime filtered feeds and historical query support
 * [Kibana dashboard](https://www.elastic.co/guide/en/kibana/current/dashboard.html) service included

The current version is `0.9`. While not recommended for production use quite yet, it's getting close
(about six weeks away). In the meantime, you can use the current playground hosted at `cloud.atomiq.io`,
and you can also host your own cluster. You can even create a full cluster on your own laptop
with `amp cluster create` using the CLI.

## Getting started

For getting started and more detailed information, see [docs](docs/).

## Contributing

If you want to hack on the project, we have a fully containerized toolchain.
All you need is Docker to build, test, and deploy! We would love for you to get involved,
so check out [CONTRIBUTING](project/CONTRIBUTING.md) and other docs under [project](project/).

## Community

If you want to chat with the developers and other members of the community, we've got an
IRC channel (#atomiq) and you can also join our Slack channel in the next day or so.

## License

AMP is licensed under the Apache License, Version 2.0.
See [LICENSE](https://github.com/appcelerator/amp/blob/master/LICENSE)
for the full license text.

