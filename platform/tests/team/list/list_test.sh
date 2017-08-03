#!/bin/bash

amp -k team ls | pcregrep -q "team"
