BEAT_NAME=ampbeat
BEAT_PATH=github.com/appcelerator/ampbeat
BEAT_GOPATH=$(firstword $(subst :, ,${GOPATH}))
BEAT_URL=https://${BEAT_PATH}
SYSTEM_TESTS=false
TEST_ENVIRONMENT=false
ES_BEATS?=${GOPATH}/src/github.com/appcelerator/amp/vendor/github.com/elastic/beats
GOPACKAGES=$(shell glide novendor)
PREFIX?=.
NOTICE_FILE=NOTICE

# Path to the libbeat Makefile
-include $(ES_BEATS)/libbeat/scripts/Makefile

# Collects all dependencies and then calls update
.PHONY: collect
collect:
