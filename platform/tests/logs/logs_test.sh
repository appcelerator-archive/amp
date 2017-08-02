#!/bin/bash

test_logs() {
  while true
  do
     if amp -k logs | grep -q "pinger_pinger.*listening on :3000"
     then
             break
     fi
     sleep 1
  done
}

test_logs_container() {
  amp -k logs --container "pinger_pinger." | grep -q "pinger_pinger.*listening on :3000"
}

test_logs_include() {
  amp -k logs -i | grep -q "amp_"
}

test_logs_metadata() {
  amp -k logs -m | grep -q ".*container_name:.*pinger_pinger.*container_state.*running.*"
}

test_logs_msg() {
  amp -k logs --msg "listening on :3000" | grep -q "pinger_pinger.*listening on :3000.*"
}

test_logs_node() {
  nodeid=$(docker node inspect self --format "{{.ID}}")
  amp -k logs --node $nodeid | grep -q "pinger_pinger.*listening on :3000.*"
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
