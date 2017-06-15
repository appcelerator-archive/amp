#!/bin/bash

amp="amp -s localhost"
$amp login --name su --password password
$amp logs -m | grep -q "timestamp:"
