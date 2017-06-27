#!/bin/bash

amp team create --team=team --org=org | grep -q "Team has been created in the organization."
