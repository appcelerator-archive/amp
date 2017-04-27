#!/bin/bash

amp -k org switch org
amp -k stack up -c examples/stacks/pinger/pinger.yml
amp -k stack up -c examples/stacks/pinger/pinger.yml pi
