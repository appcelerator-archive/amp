#!/bin/bash

amp -k org member add user
amp -k org member role --member=user --role=owner | grep -q "Member role has been changed."
amp -k org member role --member=user1 --role=OWNER | grep -q "Member role has been changed."
