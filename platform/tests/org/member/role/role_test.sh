#!/bin/bash

amp org member add user
amp org member role --member=user --role=owner | grep -q "Member role has been changed."
