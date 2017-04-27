#!/bin/bash

amp -k login --name su --password password
amp -k logs -m | grep -q "timestamp:"
