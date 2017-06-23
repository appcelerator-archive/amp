#!/bin/bash

test_stack_up() {
  amp stack up -c platform/stacks/visualizer.stack.yml
  amp stack ls 2>/dev/null | grep -Eq "visualizer"
}
