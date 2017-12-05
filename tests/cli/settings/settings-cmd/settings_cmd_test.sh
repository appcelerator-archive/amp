#!/bin/bash

# server address passed as command-line argument
test_cli_config() {
  echo $(amp -k settings) | grep -Eq "Settings file: none\s+[[:alpha:][:space:]]+: Server:\s+127.0.0.1"
}

# server address passed in local config
test_local_config() {
  mkdir -p $PWD/.amp
  echo "Server: google.com" > $PWD/.amp/amp.yml
  echo $(amp -k settings) | grep -Eq "Settings file: $PWD/.amp/amp.yml\s+[[:alpha:][:space:]]+: Server:\s+google.com"
}

test_local_cleanup() {
   rm -Rf $PWD/.amp
}

# server address passed in home config
test_home_config() {
  mkdir -p $HOME/.config/amp
  echo "Server: aws.com" > $HOME/.config/amp/amp.yml
  echo $(amp -k settings) | grep -Eq "Settings file: $HOME/.config/amp/amp.yml\s+[[:alpha:][:space:]]+: Server:\s+aws.com"
}

test_home_cleanup() {
   rm -Rf $HOME/.config/amp
}
