#!/usr/bin/env bash

amp -k login --name su --password password
TOKEN=$(cat ~/.config/amp/localhost.credentials)
curl -k --header "Authorization: amp $TOKEN" https://gw.local.appcelerator.io/v1/stacks | grep "{}"
