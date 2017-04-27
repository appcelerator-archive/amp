#!/bin/bash

test_stack_up() {
  amp -k stack up -c platform/stacks/visualizer.stack.yml
  amp -k stack ls 2>/dev/null | grep -Eq "visualizer"
}
