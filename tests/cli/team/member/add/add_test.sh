#!/bin/bash

amp -k team member add user1 user2 | grep -q "Member(s) have been added to team."
