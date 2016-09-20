# What is AMP?

AMP is the codename for a project that incubated at Atomiq, then became the Appcelerator Microservices Platform. Ultimately, the project's goal is to provide developers with facilities to deploy containerized microservices for what is now called "serverless computing" that can be activated

* on demand (for example, in response to a message or an API request)
* on schedule (think "cloud-based Cron")
* on event (by monitoring realtime message streams)

Microservices have access to system services that can help accelerate development. These services are automatically configured for high availablility and include:

* distributed key/value store
* high throughput, durable, ordered message queuing

What is referred to as AMP today is actually just the foundation for serverless computing. This foundation provides a Container-as-a-Service (CaaS), which at a high level has three important aspects. It provides:

* services for setting up and managing cluster infrastructure to host Docker containers;
* services for registering, building, and deploying Docker images;
* services for monitoring and querying multiplexed logs and metrics pertaining to cluster and application performance.

These services are accessible via a CLI, a web UI, and client libraries for supported languages.

This foundation allows developers to deploy complete containerized application stacks and manage, scale, and monitor their services running in a cluster.

## What distinguishes the foundation platform from any other CaaS solution?

AMP provides CaaS features, but it is also part of a bigger picture to support our enterprise customers in the API space. Axway and Appcelerator customers will reap the benefits of modern container technology while also enjoying the advantages of deep integration into our existing solutions for API gateways, federated security and policy enforcement, and analytics.

AMP also provides developers with excellent insight into their application's behavior. Developers can create filtered queries for logs and metrics that can return results for a specific period or can be streamed for realtime monitoring. These filtered queries allow developers to focus on specific aspects of their application's behavior, and multiple queries can execute and stream simultaneously without impacting their application's performance.

