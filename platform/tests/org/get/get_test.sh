#!/bin/bash

amp org get org | grep -q "Organization: org"
amp org get org | grep -q "Email: sample@user.amp"
