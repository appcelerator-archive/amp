#!/bin/bash

amp -k secret create TEST platform/tests/secrets/create/test | grep -o -w -E -q '[[:alnum:]]{25}'
