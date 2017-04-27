#!/bin/bash

for id in $(amp -k stack ls -q)
do
  amp -k team resource perm $id write | grep -q "Permission level has been changed."
  amp -k team resource ls | grep -q "TEAM_WRITE"
done
