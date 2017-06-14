#!/bin/bash

test_setup() {
  amp="amp -s localhost"
  $amp user signup --name user1 --password password --email email@user1.amp
}

test_stackname() {
  $amp stack up -c platform/stacks/visualizer.stack.yml
  $amp stack ls 2>/dev/null | grep -Eq "\svisualizer\s"
}

test_teardown() {
  $amp stack rm visualizer
  $amp user rm user1
}
