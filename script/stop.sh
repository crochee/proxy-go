#!bin/bash

#id=$(ps -ef |grep "proxy"|grep -v grep |awk '{print $2}')
#if [ -z "$id" ]; then
#	echo "proccessId is null"
#else
#	echo proccessId:$id is not null.
#	kill -2 $id
#fi

ps -ef | grep "proxy" | grep -v grep | awk '{print $2}' | xargs kill -2
