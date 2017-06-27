#!/bin/bash

ResId=$(amp stack ls -q)
amp team resource ls | grep -q "$ResId"
