#!/usr/bin/env bash

# Starts the API in the background and waits for it to come up. Should only
# be used for CI tests.

go run . --config-file=./configuration/local/config.yaml  > api.log 2>&1 &

while ! grep -q 'Starting HTTP server' api.log
do
  sleep .1
done

exit 0