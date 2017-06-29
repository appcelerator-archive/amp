#!/bin/bash

amp team member add user | grep -q "Member(s) have been added to team."
amp team member add user1 user2 | grep -q "Member(s) have been added to team."
