#!/bin/bash

#amp -k team member add user | grep -q "Member(s) have been added to team."
amp -k team member add user1 user2 | pcregrep -q "Member\(s\) have been added to team."
