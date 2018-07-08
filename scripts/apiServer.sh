#!/usr/bin/env bash

export LISTEN_ADDRESS=127.0.0.1:8010
export RABBITMQ_SERVER=amqp://test:test123@127.0.0.1:15672
export LOG_DIR=/var/logs/afs/
export LOG_LEVEL=DEBUG

if [ ! -d $LOG_DIR ];then
    mkdir -p $LOG_DIR
    touch $LOG_DIR/afs.log 
fi

go run ../main/apiServer.go



