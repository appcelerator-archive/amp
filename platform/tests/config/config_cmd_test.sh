#!/bin/bash

# server address passed as command-line argument
test_cli_config() {
  echo $(amp -s localhost config) | grep -E "Server:[[:space:]]+localhost"
}

# server address passed as command-line argument
test_local_config() {
  mkdir -p $PWD/.amp
  echo "Server: LOCAL" > $PWD/.amp/amp.yml
  echo $(amp config) | grep -E "Server:[[:space:]]+LOCAL"
}

test_local_cleanup() {
   rm -Rf $PWD/.amp
}

# server address passed as command-line argument
test_home_config() {
  mkdir -p $HOME/.config/amp
  echo "Server: HOME" > $HOME/.config/amp/amp.yml
  echo $(amp config) | grep -E "Server:[[:space:]]+HOME"
}

test_home_cleanup() {
   rm -Rf $HOME/.config/amp
}
