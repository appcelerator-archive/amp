#!/bin/bash

amp -k user signup --name user1 --password password --email email@user1.amp

find $HOME/.config/amp -name 'localhost*'

amp -k user rm user1
