#!/bin/bash

amp -k secret create TEST platform/tests/secrets/create/test | grep -o -w -E -q '[[:alnum:]]{25}'

amp -k secret ls | grep -q 'TEST'

amp -k secret rm TEST | grep -q 'TEST'
