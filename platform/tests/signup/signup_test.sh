#!/usr/bin/env bash

amp user signup --name user --password password --email email@user.amp
amp user ls -q | grep -q "user"
amp user rm user
