#!/bin/bash

amp -k stack up -c examples/stacks/counter/counter.yml

amp -k service ls --stack pinger | grep -vq "counter"
