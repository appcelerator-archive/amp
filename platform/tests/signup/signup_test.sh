#!/usr/bin/env bash

amp -k user signup --name user --password password --email email@user.amp
amp -k user ls -q | pcregrep -q "user"
amp -k user rm user
