#!/usr/bin/env bash


export LISTEN_ADDRESS=127.0.0.1:7070
echo "Current will test listen server address $LISTEN_ADDRESS"


curl -v $LISTEN_ADDRESS/objects/test -XPUT -d"test file objects"

