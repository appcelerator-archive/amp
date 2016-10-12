# AMP Background

Tony Pujals, Peter Svetlichny, and Henry Allen-Tilford began early incubation for the AMP project in 2015 as part of the Atomiq project, which was focused on providing a containerized environment for microservices with an emphasis on ease of deployment, orchestration, monitoring and metrics for what is now called "serverless computing."

In the Fall of 2015, he took the concept to Jeff Haynie, CEO of Appcelerator, and the project became the basis for developing a new platform at Appcelerator that would serve as the foundation for its next generation of services to provide scalability, high availability, and deployment options to multiple cluster backends, including public cloud and on-premises.

In early 2016, Appcelerator was acquired by Axway and AMP was able to leverage the resources of a distributed team that includes both Appcelerator and Axway employees from the U.S. and France.

At this time, the team began working on a new foundation in Node.js designed to fill the gaps around orchestration with the existing Docker APIs that were a vital prerequisite to achieve the project's overarching goals.

However, after working for a couple of months, Docker released version 1.12 with Swarm Mode that provided orchestration support features built into its engine. Leveraging these native facilities was an obvious decision, so the project was rebooted in July.

The team now consisted of Tony Pujals, Henry Allen-Tilford, and Chris Coy in the U.S., and Bertrand Quenin, Francois Reignat, Nicolas Degory, and Hadrien Gantzer in France. As part of the project reboot, the team was able to benefit from its initial prototyping experience and make the strategic decision to switch to Go, which is the language that Docker is implemented in.

The benefit to this choice was the ability to immediately leverage a large number of libraries being developed by Docker and other organizations that are a part of its ecosystem. Another strategic decision was to switch to Protocol Buffers and gRPC for very fast, efficient service communications and message serialization.
