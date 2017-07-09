#!/bin/bash

amp -k user signup --name owner --password password --email owner@email.amp || amp -k login --name owner --password password
