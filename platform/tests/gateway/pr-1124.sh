#!/bin/bash

curl -k -XPOST https://gw.local.atomiq.io/v1/stack | grep '{"error":"EOF","code":3}'
