#!/bin/bash

amp team resource add --org=org --team=team $(amp stack ls -q) |  grep -q "Resource(s) have been added to team."
