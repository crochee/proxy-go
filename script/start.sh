#!bin/bash
set -ex
nohup ./proxy >>./log/console.txt 2>&1 &
