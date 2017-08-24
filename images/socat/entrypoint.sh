#!/bin/sh
# if IN is already a host:port addr, we'll use it directly
# else, get the gateway, which is where the docker metrics are available
echo "$IN" | grep -q ":"
if [[ $? -ne 0 ]]; then
  gw=$(ip route | grep "default via" | awk '{print $3}')
  [[ -z "$gw" ]] && exit 1
  IN="${gw}:$IN"
fi

echo "forwarding $IN to $OUT" >&2
exec socat -d -d TCP-L:$OUT,fork TCP:$IN
