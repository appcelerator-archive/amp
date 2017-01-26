# Grafana Docker image

This project builds a Docker image with the latest master build of Grafana and customization for AMP.

## Running your Grafana container

Start your container binding the external port `3000`.

    docker run -d --name=grafana -p 3000:3000 appcelerator/grafana-amp

Try it out, default admin user is admin/changeme.

## Configuration (ENV, -e)

Same as [base image](https://github.com/appcelerator/docker-grafana)

## Tags

- ```1.0.1``` (Grafana 3.1.1)
- ```1.1.1``` (Grafana 4.0)
- ```1.1.2```, ```latest``` (Grafana 4.1)
