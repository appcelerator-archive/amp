#!/bin/bash

amp -k team remove team | pcregrep -q "team"
