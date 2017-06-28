#!/bin/bash

amp user signup --name user --password password --email email@user.amp --autologin=false
amp user signup --name user1 --password password --email email@user1.amp --autologin=false
amp user signup --name user2 --password password --email email@user2.amp --autologin=false
amp org member add user user1 user2
