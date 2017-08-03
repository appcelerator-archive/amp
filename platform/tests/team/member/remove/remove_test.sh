#!/bin/bash

amp -k team member rm user | pcregrep -q "user"
