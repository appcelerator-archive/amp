#!/bin/bash

# server address passed as command-line argument
test_cli_config() {
  echo $(amp -k settings) | pcregrep -q "Settings file: none\s+[[:alpha:][:space:]]+: Server:\s+localhost"
}

# server address passed in local config
test_local_config() {
  mkdir -p $PWD/.amp
  echo "Server: LOCAL" > $PWD/.amp/amp.yml
  echo $(amp -k settings) | pcregrep -q "Settings file: $PWD/.amp/amp.yml\s+[[:alpha:][:space:]]+: Server:\s+LOCAL"
}

test_local_cleanup() {
   rm -Rf $PWD/.amp
}

# server address passed in home config
test_home_config() {
  mkdir -p $HOME/.config/amp
  echo "Server: HOME" > $HOME/.config/amp/amp.yml
  echo $(amp -k settings) | pcregrep -q "Settings file: $HOME/.config/amp/amp.yml\s+[[:alpha:][:space:]]+: Server:\s+HOME"
}

test_home_cleanup() {
   rm -Rf $HOME/.config/amp
}
