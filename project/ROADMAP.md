# Roadmap

## Overview

This document provides a high level overview of the direction of the AMP project.

The project follows a **time-based** release process. Until the project reaches v1.0,
we are currently releasing updates every two weeks.

Roadmap items are filed on the [issue tracker](https://github.com/appcelerator/amp/issues)
with the [`roadmap`](https://github.com/appcelerator/amp/labels/roadmap) label.
They are open to community participation:

* Take a look at the [contributing guide](CONTRIBUTING.md).
* Join the [conversation](../README.md#join-the-conversation) with other users and contributors.

Anyone can open an issue to suggest a feature or report a bug. If you see an issue
that you would like to work on, please make sure to leave a comment indicating your
intention to avoid duplication of effort. If a maintainer is already assigned to it,
you can still offer to help.

## Roadmap

| Milestone          | Target date | Release captain                                |
|--------------------|-------------|------------------------------------------------|
| [AMP 0.4.0][0.4.0] | dec 2016  | **[@ndegory](https://github.com/ndegory)**
| [AMP 0.3.0][0.3.0] | 11/07/2016  | **[@ndegory](https://github.com/ndegory)**
| [AMP 0.2.1][0.2.1] | 10/24/2016  | **[@subfuzion](https://github.com/subfuzion)**
| [AMP 0.1.1][0.1.1] | 10/10/2016  | **[@subfuzion](https://github.com/subfuzion)**
| [AMP 0.1.0][0.1.0] | 09/26/2016  | **[@subfuzion](https://github.com/subfuzion)**

### 0.1 Basic foundation

The goal for the `0.1` milestone is a basic foundation that can run a service or a
collection of related services (stack) defined in yaml file, and provide facilities
to monitor filtered streams of logs and stats as well as to query stored logs and
stats.

### 0.2 Networking and messaging improvements

The theme for the `0.2` milestone is to add stronger networking features and to move
away from Kafka to NATS for messaging.

### 0.3 Architectural improvements, refactoring, support for workers

The goal for the `0.3` milestone is to provide separation between stacks and
the infrastructure services that launch and manage their lifecycles, and to lay
the foundation for managing workers.

### 0.4 TBD

[0.1.0]: https://github.com/appcelerator/amp/milestone/1?closed=1
[0.1.1]: https://github.com/appcelerator/amp/milestone/2?closed=1
[0.2.0]: https://github.com/appcelerator/amp/milestone/3
[0.2.1]: https://github.com/appcelerator/amp/milestone/5
[0.3.0]: https://github.com/appcelerator/amp/milestone/4
[0.4.0]: https://github.com/appcelerator/amp/milestone/6
