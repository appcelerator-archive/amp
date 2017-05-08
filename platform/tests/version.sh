#!/bin/bash

result=$(amp -s localhost version)

version="v[0..9]{1,2}\.[0..9]{1,3}\.[0..9]{1,3}"
echo $result | grep -E "Version:[[:space:]]+$version"
