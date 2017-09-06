#!/bin/bash

test_logs() {
  SECONDS=0
  local _timeout=35
  while true
  do
     amp -k logs | grep -q "pinger_pinger.*listening on :3000" && return 0
     sleep 1
     if [[ $SECONDS -gt $_timeout ]]; then
       echo "test_logs timed out, unable to find 'pinger_pinger.*listening on :3000' in logs:" >&2
       amp -k logs >&2
       amp -k service ps pinger_pinger >&2
       return 1
     fi
  done
}

test_logs_container() {
  amp -k logs --container "pinger_pinger." | grep -q "pinger_pinger.*listening on :3000"
}

test_logs_include() {
  local amplogs
  local ec
  amplogs=$(amp -k logs -i 2>/dev/null)
  if [[ -z "$amplogs" ]]; then
    echo "no logs available yet" >&2
    return 1
  fi
  echo $amplogs | grep -q "amp_"
  ec=$?
  [[ $? -ne 0 ]] && echo $amplogs
  return $ec
}

test_logs_metadata() {
  amp -k logs -m | grep -q ".*container_name:.*pinger_pinger.*container_state.*running.*"
}

test_logs_msg() {
  amp -k logs --msg "listening on :3000" | grep -q "pinger_pinger.*listening on :3000.*"
}

test_logs_node() {
  local code
  local amplogs
  nodeid=$(docker node inspect self --format "{{.ID}}")
  amplogs=$(amp -k logs --node $nodeid 2>/dev/null)
  echo $amplogs | grep -q "pinger_pinger.*listening on :3000.*"
  code=$?
  [[ $code -ne 0 ]] && echo $amplogs
  return $code
}

test_logs_number() {
  amp -k logs -n 2 | wc -l | grep -q "2"
}

test_logs_raw() {
  amp -k logs -r | grep -q ".*listening on :3000"
}

test_logs_regexp() {
  amp -k logs --regexp --msg ".*listening.*" | grep -q "pinger_pinger.*listening on :3000.*"
}

test_logs_since() {
  amp -k logs --since 1 | grep -q "pinger_pinger.*listening on :3000"
}

test_logs_stack() {
  amp -k logs --stack pinger | grep -q "pinger_pinger.*listening on :3000"
}
