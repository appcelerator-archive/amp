#!/bin/bash

amp -k org get org | grep -q "Organization: org"
amp -k org get org | grep -q "Email: sample@user.amp"
