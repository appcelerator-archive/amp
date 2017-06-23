#!/bin/bash

amp login --name su --password password
amp logs -m | grep -q "timestamp:"
