#!/bin/bash

for id in $(amp stack ls -q)
do
  amp team resource perm $id write | grep -q "Permission level has been changed."
  amp team resource ls | grep -q "TEAM_WRITE"
done
