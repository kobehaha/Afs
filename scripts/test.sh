#!/usr/bin/env bash


export LISTEN_ADDRESS=127.0.0.1:8010
echo "Current will test listen server address $LISTEN_ADDRESS"

object="ddddddd"

rm -r -f /var/objects/test
value=$(echo -n $object | openssl dgst -sha256 | base64)
echo "sha256 value --> " + $value

curl -v $LISTEN_ADDRESS/objects/test -XPUT -d$object -H "Digest: SHA-256=$value"

