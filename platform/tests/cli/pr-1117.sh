#!/bin/bash

# verify the command 'service tasks' runs without any error

set -e

id=$(docker service ls -q | head -n 1)
amp service tasks $id
