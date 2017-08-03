#!/bin/bash

amp -k team get team | pcregrep -q "Team: team"
#amp -k team get team | grep -q "Organization: org"
#amp -k team get team | grep -q "Organization: default"
