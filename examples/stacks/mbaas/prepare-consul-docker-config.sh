#!/bin/bash

echo "checking prerequisite: jq"
which jq >/dev/null || exit 1
echo "checking prerequisite: base64"
which base64 >/dev/null || exit 1

CONFIG=$(dirname $0)/consul.cloud.json
TMPFILE=$(mktemp)
PROTOCOL=${PROTOCOL:-https}
DOMAIN=${DOMAIN:-mbaas.aws.appcelerator.io}
DASHBOARD_URL=${DASHBOARD_URL:-${PROTOCOL}://dashboard.$DOMAIN}
ENCODED_DOMAIN=$(echo $DOMAIN | base64)
ENCODED_DASHBOARD_URL=$(echo $DASHBOARD_URL | base64)

echo "$CONFIG will be updated with domain = $DOMAIN and dashboard_url = $DASHBOARD_URL."
echo " To set the new values, export the env variables DOMAIN and DASHBOARD_URL before running this script."
echo " To cancel, Ctrl-C, else press Enter."
read pause

DOMAIN_KEYS=$(jq -r '.[].key' $CONFIG | grep domain)
DASHBOARD_URL_KEYS=$(jq -r '.[].key' $CONFIG | grep dashboard_url$)
# TODO: same with the password fields

cp $CONFIG $TMPFILE
for k in $DOMAIN_KEYS; do
  org=$(eval jq -r "'.[] | select(.key == \"$k\").value'" $CONFIG | base64 -D)
  if [[ "x$org" != "x$DOMAIN" ]]; then
    echo "updating key $k ($org to $DOMAIN)"
    eval "jq '[.[] | select(.key == \"$k\").value=\"$ENCODED_DOMAIN\"]' $TMPFILE" > ${TMPFILE}.json || exit 1
    mv ${TMPFILE}.json $TMPFILE
  else
    echo "key $k was already set at the requested value"
  fi
done
for k in $DASHBOARD_URL_KEYS; do
  org=$(eval jq -r "'.[] | select(.key == \"$k\").value'" $CONFIG | base64 -D)
  if [[ "x$org" != "x$DASHBOARD_URL" ]]; then
    echo "updating key $k ($org to $DASHBOARD_URL)"
    eval "jq '[.[] | select(.key == \"$k\").value=\"$ENCODED_DASHBOARD_URL\"]' $TMPFILE" > ${TMPFILE}.json || exit 1
    mv ${TMPFILE}.json $TMPFILE
  else
    echo "key $k was already set at the requested value"
  fi
done
mv $TMPFILE $CONFIG
