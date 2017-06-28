#!/bin/bash

ResId=$(amp stack ls -q)
amp team resource add --org=org --team=team $ResId |  grep -q "Resource has been added to team."
