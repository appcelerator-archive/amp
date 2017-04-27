#!/bin/bash

# server passed on commandline overrides both local and home configs
test_cli_local_home() {
  # create local config
  mkdir -p $PWD/.amp
  echo "Server: LOCAL" > $PWD/.amp/amp.yml

  # create home config
  mkdir -p $HOME/.config/amp
  echo "Server: HOME" > $HOME/.config/amp/amp.yml

  # server address passed as command line argument
  amp -k -s SERVER version | pcregrep -q "Server:[[:space:]]+SERVER"
  local ec=$?

  rm -R $PWD/.amp
  rm -R $HOME/.config/amp
  return $ec
}

# server passed in local config overrides home config
test_local_home() {
  # create local config
  mkdir -p $PWD/.amp
  echo "Server: LOCAL" > $PWD/.amp/amp.yml

  # create home config
  mkdir -p $HOME/.config/amp
  echo "Server: HOME" > $HOME/.config/amp/amp.yml

  amp -k version | pcregrep -q "Server:[[:space:]]+LOCAL"
  local ec=$?

  rm -R $PWD/.amp
  rm -R $HOME/.config/amp
  return $ec
}

# server passed in home config
test_home() {
  # create home config
  mkdir -p $HOME/.config/amp
  echo "Server: HOME" > $HOME/.config/amp/amp.yml

  amp -k version | pcregrep -q "Server:[[:space:]]+HOME"
  local ec=$?

  rm -R $HOME/.config/amp
  return $ec
}

# no server passed; read default server address
#test_cloud() {
#  amp version | pcregrep -q "Server:\s+cloud.appcelerator.io"
#}

test_teardown() {
  rm -Rf $PWD/.amp
  rm -Rf $HOME/.config/amp
}
