#!/usr/bin/env bash

amp="amp -s localhost"
$amp user signup --name user --password password --email email@user.amp
$amp login --name user --password password
TOKEN=$(cat ~/.config/amp/localhost.credentials)
curl -kvs -H "Authorization: amp $TOKEN" https://gw.local.appcelerator.io/v1/stacks | grep "{}"
