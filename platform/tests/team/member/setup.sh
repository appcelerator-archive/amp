#!/bin/bash

amp user signup --name user --password password --email email@user.amp --autologin=false
amp org member add user
