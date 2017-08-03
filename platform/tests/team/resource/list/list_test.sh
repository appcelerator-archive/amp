#!/bin/bash

for id in $(amp -k stack ls -q)
do
  amp -k team resource ls | pcregrep -q $id
done
