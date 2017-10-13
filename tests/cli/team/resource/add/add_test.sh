#!/bin/bash

amp -k team resource add --team=team $(amp -k stack ls -q) |  grep -q "Resource(s) have been added to team."
