#!/bin/bash

ResId=$(amp stack ls -q)
amp team resource perm $ResId write | grep -q "Permission level has been changed."
amp team resource ls | grep -q "TEAM_WRITE"
