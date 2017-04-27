#!/bin/bash

ResId=$(amp -k stack ls -q)
amp -k team resource ls | grep -q "$ResId"
