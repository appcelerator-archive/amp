#!/bin/bash

amp -k stack services global | pcregrep -q "global"
