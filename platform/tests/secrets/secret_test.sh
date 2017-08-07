#!/bin/bash

# on success, amp secret create returns an alphanumeric ID of length 25
amp -k secret create TEST platform/tests/secrets/create/test | pcregrep -q '[[:alnum:]]{25}'

amp -k secret ls | pcregrep -q 'TEST'

amp -k secret rm TEST | pcregrep -q 'TEST'


echo hello | amp secret create HELLO - | pcregrep -q '[[:alnum:]]{25}'

amp -k secret ls | pcregrep -q 'HELLO'

amp -k secret rm HELLO | pcregrep -q 'HELLO'
