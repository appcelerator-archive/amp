#!/bin/bash
#
docker build -t appcelerator/pinger . \
	&& docker tag appcelerator/pinger appcelerator/pinger:$(cat VERSION | sed 's/[:space:]*$//')
