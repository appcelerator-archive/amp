#!/bin/bash

ResId=$(amp stack ls -q)
amp team resource rm $ResId | grep -q "$ResId"
