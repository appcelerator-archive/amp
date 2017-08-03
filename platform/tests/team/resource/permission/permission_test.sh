#!/bin/bash

for id in $(amp -k stack ls -q)
do
  amp -k team resource perm $id write | pcregrep -q "Permission level has been changed."
  amp -k team resource ls | pcregrep -q "TEAM_WRITE"
  amp -k team resource perm $id READ | pcregrep -q "Permission level has been changed."
  amp -k team resource ls | pcregrep -q "TEAM_READ"
done
