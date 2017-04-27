#!/bin/bash

ResId=$(amp -k stack ls -q)
amp -k team resource rm $ResId | grep -q "$ResId"
