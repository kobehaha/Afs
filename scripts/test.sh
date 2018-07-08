#!/usr/bin/env bash


export LISTEN_ADDRESS=127.0.0.1:8010
echo "Current will test listen server address $LISTEN_ADDRESS"


rm -r -f /var/objects/test
curl -v $LISTEN_ADDRESS/objects/test -XPUT -d"test file objectasdfs"

