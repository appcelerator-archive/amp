#!/bin/bash

for id in $(amp -k stack ls -q)
do
  amp -k team resource rm $id | pcregrep -q $id
done
