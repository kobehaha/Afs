#!/usr/bin/env bash

export LISTEN_ADDRESS=127.0.0.1:7070
export RABBITMQ_SERVER=amqp://test:test123@127.0.0.1:15672 
export STORAGE_ROOT=/var/
export LOG_DIR=/var/logs/afs/
export LOG_LEVEL=DEBUG

if [ ! -d $LOG_DIR ];then
    mkdir -p $LOG_DIR
    touch $LOG_DIR/afs.log 
fi

if [ -d $STORAGE_ROOT ];then
    go run ../main/dataServer.go
else
    mkdir -p $STORAGE_ROOT
    go run ../main/dataServer.go
fi



