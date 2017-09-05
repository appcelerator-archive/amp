#!/bin/bash

amp -k team member ls | pcregrep -M -q '.*owner.*\n.*user1.*\n.*user2'
#amp -k team member ls | pcregrep -q "user1"
#amp -k team member ls | pcregrep -q "user2"
