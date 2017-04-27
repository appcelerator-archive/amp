#!/bin/bash

amp -k org member add user
amp -k org member role --member=user --role=owner | grep -q "Member role has been changed."
