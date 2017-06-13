#!/bin/bash

amp="amp -s localhost"

test_setup() {
  $amp user signup --name user2 --password password --email user2@xmail
}

test_name() {
  find $HOME/.config/amp -name 'localhost*'
}

test_teardown() {
  $amp user rm user2
}
