#!/bin/bash

amp -k team member ls | grep -q "user"
amp -k team member ls | grep -q "user1"
amp -k team member ls | grep -q "user2"
