# docker-influxdb-amp

[InfluxDB](https://influxdata.com/time-series-platform/influxdb/) image based on Alpine linux, customized for AMP.

## Usage

To create the image `appcelerator/influxdb-amp`, execute the following command in this folder:

    docker build -t appcelerator/influxdb-amp .

You can now push new image to the registry:

    docker push appcelerator/influxdb-amp


## Running your InfluxDB image

Start your image binding the external ports `8083` and `8086` in all interfaces to your container.

    docker run -d -p 8083:8083 -p 8086:8086 appcelerator/influxdb-amp


## Configuration (ENV, -e)

Same as [base image](https://github.com/appcelerator/docker-influxdb)

## Tags

- `1.1.3`, `latest` (influxdb 1.1.1)
