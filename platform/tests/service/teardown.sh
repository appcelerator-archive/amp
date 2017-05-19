#!/bin/bash

amp="amp -s localhost"
$amp stack rm global
$amp stack rm replicated
$amp stack rm pinger
$amp stack rm counter
