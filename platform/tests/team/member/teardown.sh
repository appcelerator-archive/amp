#!/bin/bash

amp -k login --name su --password password
amp -k user rm user user1 user2
amp -k login --name owner --password password
