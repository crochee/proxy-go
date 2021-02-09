#!bin/bash

processid=$(ps -ef | grep "proxy" | grep -v grep | awk '{print $2}')

if [ -z "$processid" ]; then
  echo "processid is null."
else
  echo processid:$processid is not null.
  kill -2 $processid
fi