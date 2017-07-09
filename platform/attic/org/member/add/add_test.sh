#!/bin/bash

amp -k org member add user | grep -q "Member(s) have been added to organization."
amp -k org member add user1 user2 | grep -q "Member(s) have been added to organization."
