#!/bin/bash

#amp -k team create team --org=org | grep -q "Team has been created in the organization."
amp -k team create team | grep -q "Team has been created."
