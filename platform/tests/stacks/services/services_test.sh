#!/bin/bash

amp -s localhost stack services global | grep -q "global"
