#!/bin/bash

amp org create --org=org --email=sample@user.amp | grep -q "Organization has been created."
